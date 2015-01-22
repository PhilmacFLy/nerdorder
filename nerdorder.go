package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var fronterr string
var ordererr string
var ordercount int

const errormessage = `<div class="alert alert-danger" role="alert">
  <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span>
  <span class="sr-only">Error:</span>
  $MESSAGE$
</div>`

type Frontpage struct {
	Lists      []List
	Message    template.HTML
	Username   string
	Ordercount int
	Shops      []Shop
}

type Orderpage struct {
	Orders     []Order
	Username   string
	Ordercount int
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
		fronterr = BuildMessage(errormessage, err.Error())
		log.Println(err)
	}

	fp.Shops, err = LoadShops()

	if err != nil {
		fronterr = BuildMessage(errormessage, err.Error())
		log.Println(err)
	}

	t, err := template.ParseFiles("templates/main.html")

	if err != nil {
		log.Fatal(err)
	}

	fp.Username = u.Username
	fp.Ordercount = ordercount

	if fronterr != "" {
		fp.Message = template.HTML(fronterr)
		fronterr = ""
	}

	err = t.Execute(w, &fp)
	if err != nil {
		log.Println(err)
	}
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
	op.Ordercount = ordercount

	err = t.Execute(w, &op)
	if err != nil {
		log.Println(err)
	}

}

func loginhandler(w http.ResponseWriter, r *http.Request) {
}

func listchangehandler(w http.ResponseWriter, r *http.Request) {

	u := User{}
	u.Username = "PhilmacFLy"

	l := List{}
	l.Name = r.FormValue("list")
	l.Owner = u.Username

	a := r.FormValue("action")

	var err error

	if a != "new" {
		e := l.Load()

		if e != nil {
			log.Println(e)
			fronterr = BuildMessage(errormessage, err.Error())
		}
	}

	switch a {
	case "add":
		i := ListItem{}
		i.Name = r.FormValue("name")
		if i.Name == "" {
			e := errors.New("Kein Name angegeben")
			log.Println(e)
			fronterr = BuildMessage(errormessage, e.Error())
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		i.Artnr = r.FormValue("artnr")
		if i.Artnr == "" {
			e := errors.New("Kein Name angegeben")
			log.Println(e)
			fronterr = BuildMessage(errormessage, e.Error())
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		var e error
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
		artnr := r.FormValue("artnr")
		err = l.RemoveItem(artnr)
	case "new":
		err = l.Create()
	case "remove":
		err = l.Delete()
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
