package main

import (
	"fmt"
	tmpl "html/template"
	"io"
	"net/http"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func StartServer() {
    log.Info().Str("port", port[1:]).Msg("Webserver started")
    http.Handle("/", middle(rootHandle))
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
    if len(r.URL.Path) >= 6 && r.URL.Path[0:6] == "/embed" {
        embedHandle(w, r)
        return
    }
    if len(r.URL.Path) >= 12 && r.URL.Path[0:12] == "/downloadall" {
        downloadAllHandle(w, r)
        return
    }
    if len(r.URL.Path) >= 9 && r.URL.Path[0:9] == "/download" {
        downloadHandle(w, r)
        return
    }
    if len(r.URL.Path) >= 6 && r.URL.Path[0:6] == "/image" || r.URL.Path == "/favicon.ico" {
        imageHandle(w, r)
        return
    }
    if r.URL.Path == "/" || len(r.URL.Path) >= 5 && r.URL.Path[0:5] == "/open" {
        openHandle(w, r)
        return
    }
}

func imageHandle(w http.ResponseWriter, r *http.Request) {
    context, err := GetFileForVisit(strings.TrimPrefix(r.URL.Path, "/image"))
    if err != nil {
        http.NotFound(w, r)
        return
    }
    if strings.Contains(context.MimeType, "image") {
        w.Header().Set("Content-Type", context.MimeType)
        file, err := os.Open(context.RealPath)
        if err != nil {
            http.NotFound(w, r)
            return
        }
        io.Copy(w, file)
    } else {
        http.NotFound(w, r)
    }
}

func embedHandle(w http.ResponseWriter, r *http.Request) {
    http.FileServer(http.FS(embedFS)).ServeHTTP(w, r)
}

func openHandle(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    context, err := GetFileForVisit(strings.TrimPrefix(r.URL.Path, "/open"))
    if err != nil {
        log.Error().Err(err).Msg("Failed to get file for visit")
        http.NotFound(w, r)
        return
    }

    // log.Info().Str("path", context.RealPath).Str("mime", context.MimeType).Msg("MIME")
    var t *tmpl.Template
    if context.IsDir {
        t = tmpl.Must(tmpl.ParseFS(
            embedFS,
            "embed/templates/base.html",
            "embed/templates/page.html",
            "embed/templates/content_catalog.html",
        ))
    } else {
        var preview string
        mtype := context.MimeType
        switch {
        case strings.Contains(mtype, "image"):
            preview = "image" 
        case strings.Contains(mtype, "text"):
            // limit display size to half a MB
            if context.Size > 524_288 {
                preview = "text" 
                context.Extra = map[string]interface{}{
                    "HighlightedSuccess": false,
                }
                break
            }
            preview = "text" 
            hl, err := HighlightText(context.RealPath)
            if err != nil {
                log.Error().Err(err).
                    Str("file", context.RealPath).
                    Msg("Failed to highlight text for file")
            }
            context.Extra = map[string]interface{}{
                "HighlightedText": tmpl.HTML(hl),    
                "HighlightedSuccess": true,
            }
        default:
            preview = "generic"
        }
        t = tmpl.Must(tmpl.ParseFS(
            embedFS,
            "embed/templates/base.html",
            "embed/templates/page.html",
            "embed/templates/content_preview.html",
            fmt.Sprintf("embed/templates/preview_%s.html", preview),
        ))
    }
    if err := t.Execute(w, context); err != nil {
        log.Error().Err(err).Msg("Failed to execute template")
    }
    return
}
