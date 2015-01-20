package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type Frontpage struct {
	Lists    []List
	Username string
}

type Orderpage struct {
	Orders   []Order
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

func ordershandler(w http.ResponseWriter, r *http.Request) {

	u := User{}

	u.Username = "PhilmacFLy"

	var err error

	op := Orderpage{}

	op.Orders, err = LoadOrders()

	if err != nil {
		log.Println(err)
	}

	t, err := template.ParseFiles("templates/orders.html")

	if err != nil {
		log.Fatal(err)
	}

	op.Username = u.Username

	t.Execute(w, &op)
}

func loginhandler(w http.ResponseWriter, r *http.Request) {
}

func listchangehandler(w http.ResponseWriter, r *http.Request) {

	u := User{}
	u.Username = "PhilmacFLy"

	l := List{}
	l.Name = r.FormValue("list")
	l.Owner = u.Username

	err := l.Load()

	if err != nil {
		log.Println(l)
		return
	}

	fmt.Println(l)

}

func main() {
	serveSingle("/favicon.ico", "static/favicon.ico")
	http.HandleFunc("/", listshandler)
	http.HandleFunc("/orders", ordershandler)
	http.HandleFunc("/list", listchangehandler)
	http.HandleFunc("/login", loginhandler)
	http.HandleFunc("/static/", statichandler)
	http.ListenAndServe(":8080", nil)

}
