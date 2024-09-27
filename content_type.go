package main

import (
	fp "path/filepath"
)

var (
    imageTypes = []string{".png", ".jpg", ".jpeg", ".webp", ".ico"}
)

func ContentType(filePath string) string {
    for _, it := range imageTypes {
        if fp.Ext(filePath) == it {
            return "image"
        }
    }
    return "generic"
}
