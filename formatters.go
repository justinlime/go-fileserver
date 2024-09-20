package main

import (
    "fmt"
)

func PrettyBytes(bytes int64) string {
    var pretty string
    b := float64(bytes)
    switch {
    case bytes < 1000:
        pretty = fmt.Sprintf("%0.0fB", b)
    case bytes < 1_000_000:
        pretty = fmt.Sprintf("%0.2fKB", b/1_000)
    case bytes < 1_000_000_000:
        pretty = fmt.Sprintf("%0.2fMB", b/1_000_000)
    case bytes < 1_000_000_000_000:
        pretty = fmt.Sprintf("%0.2fGB", b/1_000_000_000)
    default:
        pretty = fmt.Sprintf("%0.2fTB", b/1_000_000_000_000)
    } 
    return pretty
}

func PrettyTime(seconds float64) string {
    var pretty string
    switch {
    case seconds < 60:
       pretty = fmt.Sprintf("%0.2fsecs", seconds) 
    case seconds < 3600:
       pretty = fmt.Sprintf("%0.2fmins", seconds/60) 
    case seconds < 216000:
       pretty = fmt.Sprintf("%0.2fhours", seconds/3600) 
    default:
       pretty = fmt.Sprintf("%0.2fdays", seconds/216000) 
    }
    return pretty 
}
