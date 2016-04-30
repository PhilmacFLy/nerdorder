package main

import (
	"crypto/sha512"
	"html"
	"path/filepath"
	"strings"
)

func StripExt(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func BuildMessage(template string, message string) string {
	message = html.EscapeString(message)
	return strings.Replace(template, "$MESSAGE$", message, -1)
}

func hashPassword(password string) []byte {
	h512 := sha512.New()
	return h512.Sum([]byte(password))
}
