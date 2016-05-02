package main

import (
	"crypto/sha512"
	"html"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
)

var sc = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

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

func SetCookie(w http.ResponseWriter, u string) error {

	value := map[string]string{
		"name": u,
	}

	encoded, err := sc.Encode("nerdorder", value)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:  "nerdorder",
		Value: encoded,
		Path:  "/",
	}
	cookie.Expires = time.Now().Add(10 * 365 * 24 * time.Hour)

	http.SetCookie(w, cookie)
	return nil
}

func GetCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("nerdorder")

	if err != nil {
		return "", err
	}
	value := make(map[string]string)
	err = sc.Decode("nerdorder", cookie.Value, &value)

	if err != nil {
		return "", err
	}
	return value["name"], nil

}

func RemoveCookie(w http.ResponseWriter, r *http.Request) {
	expire := time.Now().AddDate(0, 0, 1)

	cookieMonster := &http.Cookie{
		Name:    "nerdorder",
		Expires: expire,
		Value:   "",
	}

	// http://golang.org/pkg/net/http/#SetCookie
	// add Set-Cookie header
	http.SetCookie(w, cookieMonster)
}
