package main

import (
    "os"
    "io"
    "fmt"
    "time"
    "strings"
    "net/http"
    "archive/zip"
    fp "path/filepath"

	"github.com/rs/zerolog/log"
)

var currentDownloads int64 = 0

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

type ZipSizeWriter struct {
    io.Writer
    Written int
}
func(w *ZipSizeWriter) Write(p []byte) (int, error) {
    n, err := w.Writer.Write(p)
    w.Written += n
    return n, err
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

    currentDownloads++

    log.Info().
        Str("requester_ip", GetIP(r)).
        Str("file", ffv.RealPath).
        Int64("current_downloads", currentDownloads).
        Str("available_bandwidth", fmt.Sprintf("%s/s", PrettyBytes(speedLimit/currentDownloads))).
        Msg("New Download")
    // Download
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Length", fmt.Sprintf("%d", ffv.Size))
    begin := time.Now()
    reader := &DownloadReader{Reader: file}
    io.Copy(w, reader)
    currentDownloads--
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

func downloadAllHandle(w http.ResponseWriter, r *http.Request) {
    webPath := strings.TrimPrefix(r.URL.Path, "/downloadall")
    ffv, err := GetFileForVisit(webPath)
    if err != nil {
        log.Error().Err(err).Str("path", r.URL.Path).Msg("Failed to serve download")
        http.NotFound(w,r)
        return
    }

    name := "all.zip"
    if webPath != "/" {
        name = ffv.Name + ".zip"
    }

    sizeWriter := &ZipSizeWriter{Writer: io.Discard}
    zipSizeWriter := zip.NewWriter(sizeWriter)

    zipDownloadWriter := zip.NewWriter(w)

    var canceled bool
    sizeWalker := func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        head, err := zip.FileInfoHeader(info)
        if err != nil {
            return err
        }
        head.Name = strings.TrimPrefix(strings.ReplaceAll(path, ffv.RealPath, ""), "/")
        head.Method = zip.Store
        zw, err := zipSizeWriter.CreateHeader(head)
        if err != nil {
            return err
        }
        if _, err := io.Copy(zw, file); err != nil {
            return err
        }
        return nil
    }
    downloadWalker := func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        head, err := zip.FileInfoHeader(info)
        if err != nil {
            return err
        }
        head.Name = strings.TrimPrefix(strings.ReplaceAll(path, ffv.RealPath, ""), "/")
        head.Method = zip.Store
        zw, err := zipDownloadWriter.CreateHeader(head)
        if err != nil {
            return err
        }
        reader := &DownloadReader{Reader: file}
        if _, err := io.Copy(zw, reader); err != nil {
            canceled = true
        }
        return nil
    }
    // Calc the size of the resulting zip for the header
    if err := fp.Walk(ffv.RealPath, sizeWalker); err != nil {
        log.Error().Err(err).
            Str("dir", ffv.RealPath).
            Msg("Failed to serve directory")
    }
    zipSizeWriter.Close()
    size := int64(sizeWriter.Written)
    // Download that mf
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))
    currentDownloads++
    log.Info().
        Str("requester_ip", GetIP(r)).
        Str("file", ffv.RealPath).
        Int64("current_downloads", currentDownloads).
        Str("available_bandwidth", fmt.Sprintf("%s/s", PrettyBytes(speedLimit/currentDownloads))).
        Msg("New Download")
    begin := time.Now()
    if err := fp.Walk(ffv.RealPath, downloadWalker); err != nil {
        log.Error().Err(err).
            Str("dir", ffv.RealPath).
            Msg("Failed to serve directory")
    }
    zipDownloadWriter.Close()
    l := log.Info().
             Str("requester_ip", GetIP(r)).
             Str("time_elapsed", PrettyTime(time.Since(begin).Seconds())).
             Str("download", ffv.RealPath).
             Int64("current_downloads", currentDownloads)
    if canceled {
        l.Msg("Interrupted Download")
    } else {
        l.Msg("Completed Download")
    } 
    currentDownloads--
}
