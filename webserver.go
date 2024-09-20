package main

import (
    "os"
    "io"
    "fmt"
    "time"
    "strings"
    fp "path/filepath"
    "net/http"
    tmpl "html/template"
    "github.com/rs/zerolog/log"
)

var currentDownloads int64 = 0

type DownloadWriter struct {
    http.ResponseWriter
    Progress *int64
}
func (dw DownloadWriter) Write(b []byte) (int, error) {
    n, err := dw.ResponseWriter.Write(b) 
    *dw.Progress += int64(n)
    return n, err
}



func StartServer() {
    log.Info().Str("port", port[1:]).Msg("Webserver started")
    http.Handle("/", middle(rootHandle))
    http.Handle("/deps/htmx.min.js", middle(htmxHandle))
    http.ListenAndServe(port, nil)
}

func middle(next func(http.ResponseWriter, *http.Request)) http.Handler {
    handle := http.HandlerFunc(next)
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Debug().
            Str("type", r.Method).
            Str("requester-ip", GetIP(r)).
            Str("requested-url", fp.Join(r.Host, r.URL.Path)).
            Msg("New Request")
        handle.ServeHTTP(w, r)
    })
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
    // For all other paths not explicitly defined
    if r.URL.Path != "/" {
        downloadHandle(w, r)
        return
    }
    t := tmpl.Must(tmpl.ParseFS(
        tmplFS,
        "templates/base.html",
        "templates/dir.html",
    ))
    w.Header().Set("Content-Type", "text/html")
    if err := t.Execute(w, GetFiles(DirToServe)); err != nil {
        log.Error().Err(err).Msg("Failed to execute template")
    }
}

func htmxHandle(w http.ResponseWriter, r *http.Request) {
    h, err := htmx.Open("htmx.min.js")
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to open htmx")
    }
    w.Header().Set("Content-Type", "application/javascript")
    io.Copy(w, h)
}

func downloadHandle(w http.ResponseWriter, r *http.Request) {
    var prog int64
    // Wraped http.ResponseWriter to tracked progress
    nw := DownloadWriter{
        w,
        &prog,
    }
    servePath := fp.Join(DirToServe, r.URL.Path)
    if _, err := os.Stat(servePath); err != nil {
        log.Debug().Err(err).Msg("File not found, or can't be accessed")
        http.NotFound(w, r)
        return
    }
    w.Header().Set("Content-Type", "application/octet-stream")
    // Limit download rate
    currentDownloads += 1
    as := fmt.Sprintf("%0.3f MB/s", float64(speedLimit)/float64(currentDownloads)/1000000)
    begin := time.Now()
    var canceled bool
    // Open the file and find the file size
    file, err := os.Open(servePath)
    if err != nil {
        log.Error().Err(err).Msg("Failed to open download file")
        http.NotFound(w, r)
        return
    }
    fileInfo , err := os.Stat(servePath)
    if err != nil {
        log.Error().Err(err).Msg("Failed to stat download file")
        http.NotFound(w, r)
        return
    }
    // Set the header to ensure it downloads
    w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
    log.Info().
        Str("requester-ip", GetIP(r)).
        Str("file", servePath).
        Int64("current-downloads", currentDownloads).
        Str("available-bandwidth", as).
        Msg("New Download")
    // Download
    for range time.Tick(1 * time.Second) {
        // Check the speed again every tick
        availableSpeed := speedLimit / currentDownloads
        if _, err = io.CopyN(nw, file, availableSpeed); err != nil {
            if *nw.Progress != fileInfo.Size() {
                canceled = true 
            }
            break
        }
    }
    duration := fmt.Sprintf("%0.2fs", time.Since(begin).Seconds())
    currentDownloads -= 1
    l := log.Info().
           Str("requester-ip", GetIP(r)).
           Int64("current-downloads", currentDownloads).
           Int64("downloaded-size", *nw.Progress).
           Str("elapsed", duration)
    if !canceled {
        l.Msg("Download Complete")
    } else {
        l.Msg("Download Interrupted")
    }
}

// Get the remote requesters' IP
func GetIP(r *http.Request) string {
    IP := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
    if IP == "" {
        IP = r.Header.Get("X-Real-IP")
    }
    if IP == "" {
        IP = strings.Split(r.RemoteAddr, ":")[0]
    }
    return IP
}
