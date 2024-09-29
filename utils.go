package main

import (
    "os"
    "mime"
    "strings"
    "net/http"
    fp "path/filepath"
    // "archive/zip"
)

// Get the remote requesters' IP
func GetIP(r *http.Request) string {
    IP := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
    if IP == "" {
        IP = r.Header.Get("X-Real-IP")
    }
    if IP == "" {
        IP = strings.Split(r.RemoteAddr, ":")[0]
    }
    return IP
}

func ZipDir(dir FileForVisit) {
    return
}

// Infer a mimetype based on a given file, default being application/octet-stream
// if it can't be properly inferred
func InferMimeType(realPath string) string {
    file, err := os.Open(realPath)
    if err != nil {
        return "application/octet-stream"
    }
    defer file.Close()
    // http.DetectContentType checks the first 512 bytes to infer
    buf := make([]byte, 512)
    file.Read(buf)
    mtype := http.DetectContentType(buf)
    // Fallback to inferring based on extension if http.DetectContentType fails
    if mtype == "application/octet-stream" {
        newMType := mime.TypeByExtension(fp.Ext(realPath)) 
        if newMType != "" {
           mtype = newMType 
        }
    }
    return mtype
}
