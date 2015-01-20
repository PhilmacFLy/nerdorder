package main

import (
	"path/filepath"
	"strings"
)

func StripExt(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
