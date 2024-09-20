package main

import (
    "os"
    "io"
    "fmt"
    "time"
    fp "path/filepath"
    "net/http"
    tmpl "html/template"
    "github.com/rs/zerolog/log"
)

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

func downloadHandle(w http.ResponseWriter, r *http.Request) {
    servePath := fp.Join(DirToServe, r.URL.Path)
    // Limit download rate
    currentDownloads += 1
    // Open the file and find the file size
    file, err := os.Open(servePath)
    if err != nil {
        log.Error().Err(err).Msg("Failed to open download file")
        http.NotFound(w, r)
        return
    }
    fileInfo, _ := os.Stat(servePath)
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
    var canceled bool
    for range time.Tick(1 * time.Second) {
        // Check the speed again every tick
        availableSpeed := speedLimit / currentDownloads
        if n, err := io.CopyN(w, file, availableSpeed); err != nil {
            // FIXME n only returns in chunks, not the acual written amount
            if n != fileInfo.Size() {
                canceled = false
                // canceled = true
            }
            break
        }
    }
    currentDownloads -= 1
    l := log.Info().
             Str("requester_ip", GetIP(r)).
             Str("time_elapsed", PrettyTime(time.Since(begin).Seconds())).
             Int64("current-downloads", currentDownloads)
    if canceled {
        l.Msg("Download Interrupted")
        return
    } 
    l.Str("downloaded_size", PrettyBytes(fileInfo.Size())).
      Msg("Download Complete")
}
