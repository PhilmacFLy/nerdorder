package main

import (
	"net/http"
	"text/template"
)

type Frontpage struct {
	lists []List
}

func statichandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/main.html")
	t.Execute(w, nil)

}

func main() {

	http.HandleFunc("/", handler)
	http.HandleFunc("/static/", statichandler)
	http.ListenAndServe(":8080", nil)

}
