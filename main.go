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

    //go:embed embed
    embedFS embed.FS
)

func init() {
    flag.StringVar(&DirToServe, "dir", ".", "Directory to serve")
    flag.IntVar(&p, "port", 6900, "Port to listen on.")
    flag.Float64Var(&ms, "speed", 20, "Maximum speed the server can serve")
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
    speedLimit = int64(1024 * 1024 * ms)
    log.Info().
        Str("max_speed", fmt.Sprintf("%s/s", PrettyBytes(speedLimit))).
        Str("port", port).
        Str("serving_directory", DirToServe).
        Msg("Using the following config")
}

func main() {
    StartServer() 
}
