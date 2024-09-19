package main

import (
    "os"
    "fmt"
    "time"
    fp "path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	ansiReset    = "\033[0m"
    ansiPurp     = "\033[35m"
	ansiYellow   = "\033[93m"
	ansiGreen    = "\033[92m"
	ansiRed      = "\033[91m"
	ansiBlue     = "\033[34m"
	ansiCyan     = "\033[36m"
	ansiWhite    = "\033[97m"
	ansiBack     = "\033[90m"

    ansiBPurp     = "\033[1;35m"
	ansiBYellow   = "\033[1;93m"
	ansiBGreen    = "\033[1;92m"
	ansiBRed      = "\033[1;91m"
	ansiBBlue     = "\033[1;34m"
	ansiBCyan     = "\033[1;36m"
	ansiBWhite    = "\033[1;97m"
	ansiBBack     = "\033[1;90m"
)

func InitLogger(level zerolog.Level) {
    output := zerolog.ConsoleWriter{Out: os.Stdout}
    output.FormatLevel = func(i interface{}) string {
        var ansi string
        var level string
        switch i.(string) {
        case "trace":
            level = "[TRC]"
            ansi = ansiBCyan
        case "debug":
            level = "[DBG]"
            ansi = ansiBPurp 
        case "info":
            level = "[INF]"
            ansi = ansiBGreen
        case "warn":
            level = "[WRN]"
            ansi = ansiBYellow
        case "error":
            level = "[ERR]"
            ansi = ansiBRed
        case "fatal":
            level = "[FTL]"
            ansi = ansiBRed
        default:
            level = "[???]"
            ansi = ansiBRed
        }
        return fmt.Sprintf("%s%s%s", ansi, level, ansiReset)
    }
    output.FormatMessage = func(i interface{}) string {
        return fmt.Sprintf("%s%s%s", ansiBWhite, i, ansiReset)
    }
    output.FormatFieldName = func(i interface{}) string {
        return fmt.Sprintf("%s%s=%s", ansiBCyan, i, ansiReset)
    }
    output.FormatFieldValue = func(i interface{}) string {
        return fmt.Sprintf("%s%s%s", ansiPurp, i, ansiReset)
    }
    output.FormatTimestamp = func(i interface{}) string {
        t := time.Now().Format(time.Stamp)
        return fmt.Sprintf("%s%s%s", ansiBBlue, t, ansiReset)
    }
    output.FormatCaller = func(i interface{}) string {
        val, ok := i.(string)
        if !ok {
            return ""
        }
        return fmt.Sprintf("%s%s%s", ansiBlue, fp.Base(val), ansiReset)
    }
    if level == zerolog.DebugLevel || level == zerolog.TraceLevel {
        log.Logger = zerolog.New(output).Level(level).With().Timestamp().Caller().Logger()
    } else {
        log.Logger = zerolog.New(output).Level(level).With().Timestamp().Logger()
    }
}
