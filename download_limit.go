package main

import (
    "io"
    "time"
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
