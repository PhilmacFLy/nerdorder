package main

import (
	"bytes"
	"errors"
	"fmt"
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
	Message    template.HTML
	Username   string
	Ordercount int
}

type Shoppage struct {
	Shops      []Shop
	Message    template.HTML
	Username   string
	Ordercount int
}

type Accountpage struct {
	Message    template.HTML
	Username   string
	Ordercount int
}

type Loginpage struct {
	Message template.HTML
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
	var err error

	u.Username, err = GetCookie(r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
	}

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
	var err error

	u.Username, err = GetCookie(r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
	}

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

	if fronterr != "" {
		op.Message = template.HTML(fronterr)
		fronterr = ""
	}

	err = t.Execute(w, &op)
	if err != nil {
		log.Println(err)
	}

}

func loginprocess(w http.ResponseWriter, r *http.Request) error {
	u := User{}
	u.Username = r.FormValue("username")
	err := u.Load()
	if err != nil {
		fmt.Println(err)
		return errors.New("Username or Password wrong")
	}
	h := hashPassword(r.FormValue("password"))
	if bytes.Compare(u.Password, h) != 0 {
		return errors.New("Username or Password wrong")
	}
	err = SetCookie(w, u.Username)
	if err != nil {
		return errors.New("WAWAWAWAWA Cookies not allowed")
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func loginhandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var loginerror string
	a := r.FormValue("action")
	if a == "do" {
		err = loginprocess(w, r)
		if err == nil {
			return
		} else {
			loginerror = BuildMessage(errormessage, err.Error())
		}
	}
	t, err := template.ParseFiles("templates/login.html")

	if err != nil {
		log.Fatal(err)
	}

	lp := Loginpage{template.HTML(loginerror)}

	err = t.Execute(w, &lp)
	if err != nil {
		log.Println(err)
	}
}

func listchangehandler(w http.ResponseWriter, r *http.Request) {

	u := User{}
	var err error

	u.Username, err = GetCookie(r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	l := List{}
	l.Name = r.FormValue("list")
	l.Owner = u.Username

	a := r.FormValue("action")

	if a != "new" {
		e := l.Load()

		if e != nil {
			log.Println(e)
			fronterr = BuildMessage(errormessage, err.Error())
			http.Redirect(w, r, "/", http.StatusFound)
			return
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
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)

}

func shopshandler(w http.ResponseWriter, r *http.Request) {
	u := User{}
	var err error

	u.Username, err = GetCookie(r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	sp := Shoppage{}

	sp.Shops, err = LoadShops()

	if err != nil {
		log.Println(err)
	}

	t, err := template.ParseFiles("templates/shops.html")

	if err != nil {
		log.Fatal(err)
	}

	sp.Username = u.Username
	sp.Ordercount = ordercount

	if fronterr != "" {
		sp.Message = template.HTML(fronterr)
		fronterr = ""
	}

	err = t.Execute(w, &sp)
	if err != nil {
		log.Println(err)
	}

}

func shopschangehandler(w http.ResponseWriter, r *http.Request) {
	var err error

	_, err = GetCookie(r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	a := r.FormValue("action")
	fmt.Println(a)
	var thresh int
	switch a {
	case "remove":
		s := Shop{}
		s.Name = r.FormValue("name")
		err = s.Remove()
	case "add":
		s := Shop{}
		s.Name = r.FormValue("name")
		thresh, err = strconv.Atoi(r.FormValue("thresh"))
		if err != nil {
			break
		}
		s.Threshold = float64(thresh)
		s.Website = r.FormValue("website")
		err = s.Save()
	}
	if err != nil {
		fmt.Println(err)
		fronterr = BuildMessage(errormessage, err.Error())
	}

	http.Redirect(w, r, "/shops", http.StatusFound)
}

func registerhandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var registererror string
	a := r.FormValue("action")
	if a == "do" {
		err = registerprocess(w, r)
		if err == nil {
			return
		} else {
			registererror = BuildMessage(errormessage, err.Error())
		}
	}
	t, err := template.ParseFiles("templates/register.html")
	if err != nil {
		log.Fatal(err)
	}

	lp := Loginpage{template.HTML(registererror)}

	err = t.Execute(w, &lp)
	if err != nil {
		log.Println(err)
	}
}

func registerprocess(w http.ResponseWriter, r *http.Request) error {
	u := User{}
	u.Username = r.FormValue("username")
	u.Password = hashPassword(r.FormValue("password"))
	u.Email = r.FormValue("email")
	err := u.Register()
	if err != nil {
		return err
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}
}

func accountprocess(w http.ResponseWriter, r *http.Request, u User) {
	a := r.FormValue("action")

	switch a {
	case "change":
		u.Password = hashPassword(r.FormValue("password"))
		u.Email = r.FormValue("email")
		http.Redirect(w, r, "/", http.StatusFound)
	case "delete":
		u.Remove()
		RemoveCookie(w, r)
	}
}

func accounthandler(w http.ResponseWriter, r *http.Request) {
	u := User{}
	var err error
	var accounterror string

	u.Username, err = GetCookie(r)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	a := r.FormValue("action")
	if a == "do" {
		accountprocess(w, r, u)
		return
	}

	op := Orderpage{}

	op.Orders, err = LoadOrders()

	if err != nil {
		log.Println(err)
	}

	t, err := template.ParseFiles("templates/account.html")
	if err != nil {
		log.Fatal(err)
	}

	ap := Accountpage{}
	ap.Username = u.Username
	ap.Ordercount = op.Ordercount
	ap.Message = template.HTML(accounterror)

	err = t.Execute(w, &ap)
	if err != nil {
		log.Println(err)
	}
}

func logoffhandler(w http.ResponseWriter, r *http.Request) {
	RemoveCookie(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func main() {
	serveSingle("/favicon.ico", "static/favicon.ico")
	http.HandleFunc("/", listshandler)
	http.HandleFunc("/orders", ordershandler)
	http.HandleFunc("/list", listchangehandler)
	http.HandleFunc("/login", loginhandler)
	http.HandleFunc("/logoff", logoffhandler)
	http.HandleFunc("/shops", shopshandler)
	http.HandleFunc("/shop", shopschangehandler)
	http.HandleFunc("/register", registerhandler)
	http.HandleFunc("/account", accounthandler)
	http.HandleFunc("/static/", statichandler)
	http.ListenAndServe(":8080", nil)

}
