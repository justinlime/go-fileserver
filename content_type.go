package main

import(
    "os"
    "mime"
    "net/http"
    fp "path/filepath"
)

var overrideToText []string = []string{
    // .json
    "application/json",
    // .sh
    "applicatoin/x-sh",
    // .js
    "application/javascript",
}

// Infer based on file extension. Failing that, infer based on contents of file.
// Returns a default of application/octet-stream
func InferMimeType(realPath string) string {
    ext := fp.Ext(realPath)
    // Infer mime by extension
    mtype := mime.TypeByExtension(ext)
    // Additional mappings, in case mtype TypeByExtension returns none
    if mtype != "" {
        // Override to text for certain types
        for _, t := range overrideToText {
            if mtype == t {
                return "text/plain"
            }
        }
        return mtype
    }
    // Additional mappings if mtype doesn't return anything
    for _, lang := range progExts {
        if ext == lang {
            return "text/plain"
        } 
    }
    // Infer based on contents if extenions aren't resolving
    file, err := os.Open(realPath)
    if err != nil {
        return "application/octet-stream"
    }
    defer file.Close()
    // http.DetectContentType checks the first 512 bytes to infer
    buf := make([]byte, 512)
    file.Read(buf)
    // defaults to application/octet-stream if none are found
    return http.DetectContentType(buf)
}
