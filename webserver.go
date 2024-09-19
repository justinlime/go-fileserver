package main

import (
    "os"
    "io"
    fp "path/filepath"
    "net/http"
    tmpl "html/template"
    "github.com/rs/zerolog/log"
)

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
            Str("path", r.URL.Path).
            Msg("New Request")
        handle.ServeHTTP(w, r)
    })
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
    wd, err := os.Getwd()
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to get the CWD")
    }
    // handle unknown routes
    if r.URL.Path != "/" {
        servePath := fp.Join(wd, r.URL.Path)
        if _, err := os.Stat(servePath); err != nil {
            log.Debug().Err(err).Msg("File not found, or can't be accessed")
            http.NotFound(w, r)
            return
        }
        w.Header().Set("Content-Type", "application/octet-stream")
        http.ServeFile(w, r, servePath)
        log.Info().Str("file", servePath).Msg("File served")
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
