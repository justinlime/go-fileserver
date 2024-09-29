package main

import (
	"fmt"
)

const (
    kilobyte float64 = 1024
    megabyte float64 = 1_048_576 
    gigabyte float64 = 1_073_741_824 
    terabyte float64 = 1_099_511_627_776 

    minute   float64 = 60
    hour     float64 = 3600
    day      float64 = 216_000
)

func PrettyBytes(bytes int64) string {
    var pretty string
    b := float64(bytes)
    switch {
    case b < kilobyte:
        pretty = fmt.Sprintf("%0.0f B", b)
    case b < megabyte:
        pretty = fmt.Sprintf("%0.2f KB", b/kilobyte)
    case b < gigabyte:
        pretty = fmt.Sprintf("%0.2f MB", b/megabyte)
    case b < terabyte:
        pretty = fmt.Sprintf("%0.2f GB", b/gigabyte)
    default:
        pretty = fmt.Sprintf("%0.2f TB", b/terabyte)
    } 
    return pretty
}

func PrettyTime(seconds float64) string {
    var pretty string
    switch {
    case seconds < minute:
       pretty = fmt.Sprintf("%0.2fsecs", seconds) 
    case seconds < hour:
       pretty = fmt.Sprintf("%0.2fmins", seconds/minute) 
    case seconds < day:
       pretty = fmt.Sprintf("%0.2fhours", seconds/hour) 
    default:
       pretty = fmt.Sprintf("%0.2fdays", seconds/day) 
    }
    return pretty 
}
