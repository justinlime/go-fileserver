package main

import (
    "fmt"
)

func PrettyBytes(bytes int64) string {
    var pretty string
    b := float64(bytes)
    switch {
    case bytes < 1000:
        pretty = fmt.Sprintf("%f B", b)
    case bytes < 1_000_000:
        pretty = fmt.Sprintf("%0.2f KB", b/1_000)
    case bytes < 1_000_000_000:
        pretty = fmt.Sprintf("%0.2f MB", b/1_000_000)
    case bytes < 1_000_000_000_000:
        pretty = fmt.Sprintf("%0.2f GB", b/1_000_000_000)
    default:
        pretty = fmt.Sprintf("%0.2f TB", b/1_000_000_000_000)
    } 
    return pretty
}

func PrettyTime(seconds float64) string {
    var pretty string
    switch {
    case seconds < 60:
       pretty = fmt.Sprintf("%0.2f secs", seconds) 
    case seconds < 3600:
       pretty = fmt.Sprintf("%0.2f mins", seconds/60) 
    case seconds < 216000:
       pretty = fmt.Sprintf("%0.2f hours", seconds/3600) 
    default:
       pretty = fmt.Sprintf("%0.2f days", seconds/216000) 
    }
    return pretty 
}
