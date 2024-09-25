package main

import (
	"os"
	fp "path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type FileForVisit struct {
    Name string    
    Size string
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
            f, err := os.Stat(rp)
            if err != nil {
                log.Error().
                    Str("file", rp).
                    Err(err).
                    Msg("Failed to stat file")
            }
            a := FileForVisit {
                Name: file.Name(),
                RealPath: rp,
                WebPath: strings.ReplaceAll(rp, DirToServe, ""),
                Size: PrettyBytes(f.Size()),
                IsDir: true,
            }
            available = append(available, a)
        }
    }
    for _, file := range files {
        if !file.IsDir() {
            rp := fp.Join(realPath, file.Name())
            f, err := os.Stat(rp)
            if err != nil {
                log.Error().
                    Str("file", rp).
                    Err(err).
                    Msg("Failed to stat file")
            }
            a := FileForVisit {
                Name: file.Name(),
                RealPath: rp,
                WebPath: strings.ReplaceAll(rp, DirToServe, ""),
                Size: PrettyBytes(f.Size()),
                IsDir: false,
            }
            available = append(available, a)
        }
    }
    return available
}
