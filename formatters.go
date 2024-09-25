package main

import (
	"fmt"
)

const (
    kilobyte float64 = 1024
    megabyte float64 = 1_048_576 
    gigabyte float64 = 1_073_741_824 
    terabyte float64 = 1_099_511_627_776 
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
