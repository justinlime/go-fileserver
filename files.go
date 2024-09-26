package main

import (
	"fmt"
	"os"
	fp "path/filepath"
)

type FileForVisit struct {
    Name string    
    Size int64
    PrettySize string
    RealPath string
    WebPath string
    IsDir bool
    Files []FileForVisit
}

func GetFileForVisit(webPath string) (FileForVisit, error) {
    realPath := fp.Join(DirToServe, webPath)
    fileInfo, err := os.Stat(realPath)
    if err != nil {
        return FileForVisit{}, fmt.Errorf("Failed to stat the file for visit: %v", err)
    }
    ffv := FileForVisit {
        Name: fileInfo.Name(),
        RealPath: realPath,
        WebPath: webPath,
        IsDir: fileInfo.IsDir(),
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
