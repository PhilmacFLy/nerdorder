package main

import (
	"path/filepath"
	"strings"
)

func StripExt(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func BuildMessage(template string, message string) string {
	return strings.Replace(template, "$MESSAGE$", message, -1)
}
