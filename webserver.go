package main

import (
    "os"
    "io"
    "fmt"
    "time"
    "net/http"
    fp "path/filepath"
    tmpl "html/template"
    "github.com/rs/zerolog/log"


)

type DownloadReader struct {
    io.Reader
    Progress int64
    LimitBucket int
}

func (w *DownloadReader) Read(p []byte) (int, error) {
    w.Limit(p)
    n, err := w.Reader.Read(p) 
    w.Progress += int64(n)
    return n, err
}

func (w *DownloadReader) Limit(p []byte) {
    availableSpeed := speedLimit / currentDownloads
    if len(p) + w.LimitBucket >= int(availableSpeed) {
        time.Sleep(time.Second * 1)
        w.LimitBucket = 0
    }
    w.LimitBucket += len(p)
}

var currentDownloads int64 = 0

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
            Str("requester_ip", GetIP(r)).
            Str("requested_url", fp.Join(r.Host, r.URL.Path)).
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
        "templates/sizes.html",
    ))
    w.Header().Set("Content-Type", "text/html")
    if err := t.Execute(w, GetFiles(DirToServe)); err != nil {
        log.Error().Err(err).Msg("Failed to execute template")
    }
}

// Serve the HTMX dep
func htmxHandle(w http.ResponseWriter, r *http.Request) {
    h, err := jsdepsFS.Open("jsdeps/htmx.min.js")
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to open htmx")
    }
    w.Header().Set("Content-Type", "application/javascript")
    io.Copy(w, h)
}

func pageHandle(w http.ResponseWriter, r *http.Request) {
    return
}

func downloadHandle(w http.ResponseWriter, r *http.Request) {
    servePath := fp.Join(DirToServe, r.URL.Path)
    fileInfo, err := os.Stat(servePath)
    if err != nil {
        log.Error().Err(err).Str("file", servePath).Msg("Failed to stat file")
        http.NotFound(w, r)
        return
    }
    // Open the file and find the file size
    file, err := os.Open(servePath)
    if err != nil {
        log.Error().Err(err).Msg("Failed to open download file")
        http.NotFound(w, r)
        return
    }
    defer file.Close()

    currentDownloads += 1

    log.Info().
        Str("requester_ip", GetIP(r)).
        Str("file", servePath).
        Int64("current_downloads", currentDownloads).
        Str("available_bandwidth", fmt.Sprintf("%s/s", PrettyBytes(speedLimit/currentDownloads))).
        Msg("New Download")
    // Download
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
    begin := time.Now()
    reader := &DownloadReader{Reader: file}
    io.Copy(w, reader)
    currentDownloads -= 1
    l := log.Info().
             Str("requester_ip", GetIP(r)).
             Str("time_elapsed", PrettyTime(time.Since(begin).Seconds())).
             Str("downloaded_size", PrettyBytes(reader.Progress)).
             Str("download", servePath).
             Int64("current_downloads", currentDownloads)
    if fileInfo.Size() != reader.Progress {
        l.Msg("Interrupted Download")
    } else {
        l.Msg("Completed Download")
    } 
}
