package main

import (
    "os"
    "fmt"
    "flag"
    "embed"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
    p int
    port string

    debug bool
    DirToServe string

    ms float64
    speedLimit int64

    //go:embed templates
    tmplFS embed.FS
    //go:embed htmx.min.js
    htmx embed.FS
)

func init() {
    flag.StringVar(&DirToServe, "dir", ".", "Directory to serve")
    flag.IntVar(&p, "port", 6900, "Port to listen on.")
    flag.Float64Var(&ms, "speed", 1, "Maximum speed the server can serve ")
    flag.BoolVar(&debug, "debug", false, "Enable debug logs")
    flag.Parse()
    if DirToServe == "." {
        dir, err := os.Getwd(); 
        if err != nil {
            log.Fatal().Err(err).Msg("Failed to stat CWD")
        }
        DirToServe = dir
    }
    port = fmt.Sprintf(":%d", p)
    if debug {
        InitLogger(zerolog.DebugLevel)
    } else {
        InitLogger(zerolog.InfoLevel)
    }
    speedLimit = int64(1000 * 1000 * ms)
    log.Info().
        Str("max-speed", fmt.Sprintf("%0.3f MB/s", ms)).
        Str("port", port).
        Str("serving-directory", DirToServe).
        Msg("Using the following config")
}

func main() {
    StartServer() 
}
