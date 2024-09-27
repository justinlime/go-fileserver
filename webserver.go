package main

import (
	"fmt"
	tmpl "html/template"
	"io"
	"net/http"
	"os"
	fp "path/filepath"
	"strings"
	"time"
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
    if r.URL.Path == "/" || len(r.URL.Path) >= 5 && r.URL.Path[0:5] == "/open" {
        openHandle(w, r)
        return
    }
    // Genericly handle all the rest
    http.FileServer(http.Dir(DirToServe)).ServeHTTP(w, r)
}

// Serve the embedded deps
func embedHandle(w http.ResponseWriter, r *http.Request) {
    fs := http.FileServer(http.FS(embedFS))
    fs.ServeHTTP(w, r)
}

func openHandle(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    context, err := GetFileForVisit(strings.TrimPrefix(r.URL.Path, "/open"))
    if err != nil {
        log.Error().Err(err).Msg("Failed to get file for visit")
    }
    var t *tmpl.Template
    if context.IsDir {
        t = tmpl.Must(tmpl.ParseFS(
            embedFS,
            "embed/templates/base.html",
            "embed/templates/page.html",
            "embed/templates/catalog.html",
        ))
    } else {
        if context.Ext == ".png" {
            t = tmpl.Must(tmpl.ParseFS(
                embedFS,
                "embed/templates/base.html",
                "embed/templates/page.html",
                "embed/templates/open.html",
                "embed/templates/preview-image.html",
            ))
        } else {
            t = tmpl.Must(tmpl.ParseFS(
                embedFS,
                "embed/templates/base.html",
                "embed/templates/page.html",
                "embed/templates/open.html",
                "embed/templates/preview-generic.html",
            ))
        }
    }
    if err := t.Execute(w, context); err != nil {
        log.Error().Err(err).Msg("Failed to execute template")
    }
    return
}

func directoryHandle(w http.ResponseWriter, r *http.Request) {
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
