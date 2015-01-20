package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var fronterr string
var ordererr string

const errormessage = `<div class="alert alert-danger" role="alert">
  <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span>
  <span class="sr-only">Error:</span>
  $MESSAGE$
</div>`

type Frontpage struct {
	Lists    []List
	Message  template.HTML
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
	if fronterr != "" {
		fp.Message = template.HTML(fronterr)
		fronterr = ""
	}

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

	e := l.Load()

	if e != nil {
		log.Println(l)
		return
	}

	a := r.FormValue("action")

	var err error

	switch a {
	case "add":
		i := ListItem{}
		i.Name = r.FormValue("name")
		i.Artnr = r.FormValue("artnr")
		i.Count, e = strconv.Atoi(r.FormValue("count"))
		if e != nil {
			log.Println(e)
			fronterr = BuildMessage(errormessage, e.Error())
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		i.Preis, e = strconv.ParseFloat(r.FormValue("preis"), 64)
		if e != nil {
			log.Println(e)
			fronterr = BuildMessage(errormessage, e.Error())
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		err = l.AddItem(i)
	case "delete":
	}

	if err != nil {
		log.Println(err)
		fronterr = BuildMessage(errormessage, err.Error())
	}
	http.Redirect(w, r, "/", http.StatusFound)

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