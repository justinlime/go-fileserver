package main

import (
	"os"
	"io"
	"fmt"
	"time"
	"strings"
	"net/http"
	fp "path/filepath"
	tmpl "html/template"

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
                preview = "generic" 
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

func downloadHandle(w http.ResponseWriter, r *http.Request) {
    ffv, err := GetFileForVisit(strings.TrimPrefix(r.URL.Path, "/download"))
    if err != nil {
        log.Error().Err(err).Str("path", r.URL.Path).Msg("Failed to serve download")
        http.NotFound(w,r)
        return
    }
    file, err := os.Open(ffv.RealPath)
    if err != nil {
        log.Error().Err(err).Msg("Failed to open download file")
        http.NotFound(w, r)
        return
    }
    defer file.Close()

    currentDownloads += 1

    log.Info().
        Str("requester_ip", GetIP(r)).
        Str("file", ffv.RealPath).
        Int64("current_downloads", currentDownloads).
        Str("available_bandwidth", fmt.Sprintf("%s/s", PrettyBytes(speedLimit/currentDownloads))).
        Msg("New Download")
    // Download
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Length", fmt.Sprintf("%d", ffv.Size))
    // w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", ffv.Name))
    begin := time.Now()
    reader := &DownloadReader{Reader: file}
    io.Copy(w, reader)
    currentDownloads -= 1
    l := log.Info().
             Str("requester_ip", GetIP(r)).
             Str("time_elapsed", PrettyTime(time.Since(begin).Seconds())).
             Str("downloaded_size", PrettyBytes(reader.Progress)).
             Str("download", ffv.RealPath).
             Int64("current_downloads", currentDownloads)
    if ffv.Size != reader.Progress {
        l.Msg("Interrupted Download")
    } else {
        l.Msg("Completed Download")
    } 
}
