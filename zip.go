package main

import (
    "strings"
    "archive/zip"
    fp "path/filepath"
    "io"
    "os"
)

type TestWriter struct {
    io.Writer
    Written int
}

func(w *TestWriter) Write(p []byte) (int, error) {
    n, err := w.Writer.Write(p)
    w.Written += n
    return n, err
}
// We're cooking with this
func ActualZipSize(ffv FileForVisit) int64 {
    writer := &TestWriter{Writer: io.Discard}
    zw := zip.NewWriter(writer)
    walk := func(path string, info os.FileInfo, err error) error {
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
        w, err := zw.CreateHeader(head)
        if err != nil {
            return err
        }
        if _, err := io.Copy(w, file); err != nil {
        }
        return nil
    }
    if err := fp.Walk(ffv.RealPath, walk); err != nil {
        panic(err)
    }
    // This fucks shit up if deferred, call explicitly
    zw.Close()
    return int64(writer.Written)
}
