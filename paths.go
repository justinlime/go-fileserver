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
        return available
    }
    for _, file := range files {
        var size int
        rp := fp.Join(realPath, file.Name())
        fileInfo, err := os.Stat(rp)
        if err != nil {
            log.Error().Err(err).Msg("Failed to stat file")
            continue
        }
        a := FileForVisit {
            Name: file.Name(),
            RealPath: rp,
            WebPath: strings.ReplaceAll(rp, DirToServe, ""),
        }
        if file.IsDir() {
            subPaths, err := GetPathsRecursively(a.RealPath)
            if err != nil {
                log.Error().Err(err).Msg("Failed to calculate size of dir")
                continue
            }
            for _, f := range subPaths {
                info, err := os.Stat(f) 
                if err != nil {
                    log.Error().Err(err).Msg("Failed to calculate size of dir")
                    continue
                }
                size += int(info.Size())
            }
            a.Size = PrettyBytes(int64(size))
            a.IsDir = true
        } else {
            a.Size = PrettyBytes(fileInfo.Size())
            a.IsDir = false
        }
        available = append(available, a)
    }
    return available
}

// Get the path of every single file in a directory (not including the directories themselves)
func GetPathsRecursively(dir string) ([]string, error){
    var paths []string
    var readDir func(string) error
    readDir = func (newDir string)  error {
        files, err := os.ReadDir(newDir)        
        if err != nil {
            return err
        }
        for _, file := range files {
            if file.IsDir() {
                if err := readDir(fp.Join(newDir, file.Name())); err != nil {
                    return err
                }
            } else {
                paths = append(paths, fp.Join(newDir, file.Name()))
            }
        }
        return nil
    }
    if err := readDir(dir); err != nil {
        return []string{}, err
    }
    return paths, nil
}
