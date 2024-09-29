package main

import (
	"fmt"
	"os"
	fp "path/filepath"
	"strings"
)

type FileForVisit struct {
    Name string    
    Size int64
    PrettySize string
    RealPath string
    WebPath string
    ParentWebPath string
    IsDir bool
    MimeType string
    Files []FileForVisit
}

func GetFileForVisit(webPath string) (FileForVisit, error) {
    realPath := fp.Join(DirToServe, webPath)
    fileInfo, err := os.Stat(realPath)
    if err != nil {
        return FileForVisit{}, fmt.Errorf("Failed to stat the file for visit: %v", err)
    }
    if webPath == "" {
        webPath = "/"
    }
    parent := strings.TrimSuffix(strings.TrimSuffix(webPath, fileInfo.Name()), "/")
    if parent == "" {
        parent = "/"
    }
    ffv := FileForVisit {
        Name: fileInfo.Name(),
        RealPath: realPath,
        WebPath: webPath,
        ParentWebPath: parent,
        IsDir: fileInfo.IsDir(),
        MimeType: InferMimeType(realPath),
    }
    if ffv.IsDir {
        var size int64
        files, err := os.ReadDir(realPath)
        if err != nil {
            return FileForVisit{}, fmt.Errorf("Failed to read the dir for visit: %v", err)
        }
        for _, f := range files {
            nestedWebPath := fp.Join(webPath, f.Name()) 
            nestedFFV, err := GetFileForVisit(nestedWebPath)
            if err != nil {
                return FileForVisit{}, fmt.Errorf("Failed to get nested file for visit %s", err)
            }
            ffv.Files = append(ffv.Files, nestedFFV)
            size += nestedFFV.Size
        }
        ffv.Size = size
    } else {
        ffv.Size = fileInfo.Size()
    }
    ffv.PrettySize = PrettyBytes(ffv.Size)
    return ffv, nil
}
