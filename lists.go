package main

type List struct {
	ID      int
	ShopID  int
	OwnerID int
	Items   []Item
}

type Item struct {
	Name  string
	Artnr string
	Descr string
	Preis float64
}
