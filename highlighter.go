package main

import (
	"fmt"
	"os"
	"bytes"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/alecthomas/chroma/v2/formatters/html"
)

func HighlightText(filePath string) (string, error) {
    var highlighted bytes.Buffer
    contents, err := os.ReadFile(filePath)
    if err != nil {
        return "", fmt.Errorf("Failed to read file for highlighting: %v", err)
    }
    lexer := lexers.Match(filePath)
    if lexer == nil {
        lexer = lexers.Fallback
    }
    style := styles.Get("catppuccin-mocha")
    formatter := html.New(
        html.WithLineNumbers(true),
        html.WithLinkableLineNumbers(true, "line"),
        html.WrapLongLines(false),
    )
    iterator, err := lexer.Tokenise(nil, string(contents))
    if err != nil {
        return "", fmt.Errorf("Failed to create interator for highlighting: %v", err)
    }
    err = formatter.Format(&highlighted, style, iterator)
    if err != nil {
        return "", fmt.Errorf("Failed to format highlighted contents: %v", err)
    }
    return highlighted.String(), nil
}
