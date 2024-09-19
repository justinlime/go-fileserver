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


func StartServer() {
    log.Info().Str("port", port[1:]).Msg("Webserver started")
    http.Handle("/", middle(rootHandle))
    http.Handle("/deps/htmx.min.js", middle(htmxHandle))
    http.ListenAndServe(port, nil)
}

func middle(next func(http.ResponseWriter, *http.Request)) http.Handler {
    handle := http.HandlerFunc(next)
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        IP := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
        if IP == "" {
            IP = r.Header.Get("X-Real-IP")
        }
        log.Debug().
            Str("type", r.Method).
            Str("IP", IP).
            Str("path", r.URL.Path).
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
        wd, err := os.Getwd()
        if err != nil {
            log.Fatal().Err(err).Msg("Failed to get the CWD")
        }
        servePath := fp.Join(wd, r.URL.Path)
        if _, err := os.Stat(servePath); err != nil {
            log.Debug().Err(err).Msg("File not found, or can't be accessed")
            http.NotFound(w, r)
            return
        }
        w.Header().Set("Content-Type", "application/octet-stream")
        // Limit download rate
        currentDownloads += 1
        as := fmt.Sprintf("%0.3f MB/s", float64(speedLimit)/float64(currentDownloads)/1000000)
        log.Debug().
            Str("available-bandwidth", as).
            Int64("current-downloads", currentDownloads).
            Msg("New Download")
        for range time.Tick(1 * time.Second) {
            // Check the speed again every tick
            availableSpeed := speedLimit / currentDownloads
            file, err := os.Open(servePath)
            if err != nil {
                log.Error().Err(err).Msg("Failed to open download file")
                break
            }
            fileInfo , err := os.Stat(servePath)
            if err != nil {
                log.Error().Err(err).Msg("Failed to stat download file")
                break
            }
            w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
            if _, err = io.CopyN(w, file, availableSpeed); err != nil {
                break
            }
        }
        currentDownloads -= 1
        log.Debug().
            Int64("current-downloads", currentDownloads).
            Msg("Download Finished")
        log.Info().Str("file", servePath).Msg("File served")
}
