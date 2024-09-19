package main

import (
    "fmt"
    "embed"
    "flag"
    "os"
	"github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"
)

var (
    p int
    port string
    debug bool
    DirToServe string
    maxdl int

    //go:embed templates
    tmplFS embed.FS
    //go:embed htmx.min.js
    htmx embed.FS
)

func init() {
    flag.StringVar(&DirToServe, "dir", ".", "Directory to serve")
    flag.IntVar(&p, "port", 6900, "Port to listen on.")
    flag.IntVar(&maxdl, "maxdl", 5, "Maximum amount of downloads to serve at once")
    flag.BoolVar(&debug, "debug", false, "Enable debug logs")
    flag.Parse()
    if DirToServe == "." {
        if dir, err := os.Getwd(); err == nil {
            DirToServe = dir
        }
    }
    port = fmt.Sprintf(":%d", p)
    if debug {
        InitLogger(zerolog.DebugLevel)
    } else {
        InitLogger(zerolog.InfoLevel)
    }
}

func main() {
    StartServer() 
}
