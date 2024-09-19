package main

import (
    "os"
    "strings"
    fp "path/filepath"
    "github.com/rs/zerolog/log"
)

type FileForVisit struct {
    Name string    
    RealPath string
    WebPath string
    IsDir bool
}

func GetFiles(realPath string) []FileForVisit {
    var available []FileForVisit
    files, err := os.ReadDir(realPath)
    if err != nil {
        log.Error().Err(err).Msg("Failed to read dir")
    }
    for _, file := range files {
        if file.IsDir() {
            rp := fp.Join(realPath, file.Name())
            a := FileForVisit {
                Name: file.Name(),
                RealPath: rp,
                WebPath: strings.ReplaceAll(rp, DirToServe, ""),
                IsDir: true,
            }
            available = append(available, a)
        }
    }
    for _, file := range files {
        if !file.IsDir() {
            rp := fp.Join(realPath, file.Name())
            a := FileForVisit {
                Name: file.Name(),
                RealPath: rp,
                WebPath: strings.ReplaceAll(rp, DirToServe, ""),
                IsDir: false,
            }
            available = append(available, a)
        }
    }
    return available
}
