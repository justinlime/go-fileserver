package main

import (
    "strings"
    "net/http"
)

// Get the remote requesters' IP
func GetIP(r *http.Request) string {
    IP := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
    if IP == "" {
        IP = r.Header.Get("X-Real-IP")
    }
    if IP == "" {
        IP = strings.Split(r.RemoteAddr, ":")[0]
    }
    return IP
}
