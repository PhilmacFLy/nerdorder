package main

import (
	"log"
	"net/http"
	"text/template"
)

type Frontpage struct {
	Lists    []List
	Username string
}

func statichandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func listshandler(w http.ResponseWriter, r *http.Request) {

	u := User{}

	u.Username = "PhilmacFLy"

	var err error

	fp := Frontpage{}

	fp.Lists, err = u.loadLists()

	if err != nil {
		log.Println(err)
	}

	t, err := template.ParseFiles("templates/main.html")

	if err != nil {
		log.Fatal(err)
	}

	fp.Username = u.Username

	t.Execute(w, &fp)
}

func orderhandler(w http.ResponseWriter, r *http.Request) {
}

func loginhandler(w http.ResponseWriter, r *http.Request) {
}

func main() {
	serveSingle("/favicon.ico", "static/favicon.ico")
	http.HandleFunc("/", listshandler)
	http.HandleFunc("/order", orderhandler)
	http.HandleFunc("/login", loginhandler)
	http.HandleFunc("/static/", statichandler)
	http.ListenAndServe(":8080", nil)

}
