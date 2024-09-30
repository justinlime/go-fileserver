package main

import(
    "os"
    "mime"
    "net/http"
    fp "path/filepath"
)

var overrideToText []string = []string{
    "application/json",
    "applicatoin/x-sh",
    "application/javascript",
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
    // defaults to application/octet-stream
    mtype := http.DetectContentType(buf)
    // Fallback to inferring based on extension if http.DetectContentType fails
    if mtype == "application/octet-stream" {
        ext := fp.Ext(realPath)
        newMType := mime.TypeByExtension(ext) 
        if newMType != "" {
           mtype = newMType 
        } else {
            for _, lang := range progExts {
                if ext == lang {
                    mtype = "text/plain"
                    break
                } 
            }
        }
    }
    for _, t := range overrideToText {
        if mtype == t {
            mtype = "text/plain"
        }
    }
    return mtype
}
