package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type List struct {
	Name  string
	Owner string
	Items []ListItem
}

type Order struct {
	Name  string
	Items []OrderItem
	Total float64
}

type ListItem struct {
	Name  string
	Artnr string
	Count int
	Preis float64
}

type OrderItem struct {
	Name        string
	Artnr       string
	Count       int
	Preis       float64
	Gesamtpreis float64
	Owner       string
}

func (l *List) Save() error {
	b, err := json.MarshalIndent(&l, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile("lists/"+l.Owner+"/"+l.Name+".json", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (l *List) Load() error {
	body, err := ioutil.ReadFile("lists/" + l.Owner + "/" + l.Name + ".json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &l)
	if err != nil {
		return err
	}
	return nil
}

func (l *List) AddItem(i ListItem) error {
	for _, it := range l.Items {
		if strings.EqualFold(i.Artnr, it.Artnr) {
			e := errors.New("Item bereits vorhanden")
			return e
		}
	}
	l.Items = append(l.Items, i)
	err := l.Save()
	if err != nil {
		return err
	}

	oi := OrderItem{}
	o := Order{}
	o.Name = l.Name
	oi.Artnr = i.Artnr
	oi.Count = i.Count
	oi.Name = i.Name
	oi.Owner = l.Owner
	oi.Preis = i.Preis
	oi.Gesamtpreis = oi.Preis * float64(oi.Count)

	err = o.AddItem(oi)

	if err != nil {
		l.RemoveItem(i.Artnr)
		return err
	}

	return nil
}

func (l *List) RemoveItem(a string) error {
	for i, it := range l.Items {
		if strings.EqualFold(a, it.Artnr) {
			l.Items = append(l.Items[:i], l.Items[i+1:]...)
		}
	}

	err := l.Save()
	if err != nil {
		return err
	}

	o := Order{}
	o.Name = l.Name

	o.RemoveItem(a)

	return nil
}

func (o *Order) Save() error {
	o.CalcTotal()
	b, err := json.MarshalIndent(&o, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile("orders/"+o.Name+".json", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (o *Order) Load() error {
	body, err := ioutil.ReadFile("orders/" + o.Name + ".json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &o)
	if err != nil {
		return err
	}
	return nil
}

func (o *Order) AddItem(oi OrderItem) error {
	if _, err := os.Stat("orders/" + o.Name + ".json"); os.IsNotExist(err) {
		o.Items = append(o.Items, oi)
		return o.Save()
	}
	if _, err := os.Stat("orders/" + o.Name + ".json"); err == nil {
		e := o.Load()
		if e != nil {
			log.Println(e)
			return e
		}
		o.Items = append(o.Items, oi)
		e = o.Save()
		if e != nil {
			log.Println(e)
			return e
		}

	}
	LoadOrders()
	return nil

}

func (o *Order) RemoveItem(a string) error {
	for i, oi := range o.Items {
		if strings.EqualFold(a, oi.Artnr) {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
		}
	}

	err := o.Save()
	if err != nil {
		return err
	}
	return nil

}

func (o *Order) CalcTotal() {
	for _, oi := range o.Items {
		o.Total += oi.Gesamtpreis
	}
}

func LoadOrders() ([]Order, error) {
	var orders []Order
	files, err := ioutil.ReadDir("orders")
	if err != nil {
		return nil, err
	}
	ordercount = 0
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			o := Order{}
			o.Name = StripExt(f.Name())
			err := o.Load()
			if err != nil {
				return nil, err
			}
			s := Shop{}
			s.Name = o.Name
			s.Load()
			if o.Total >= s.Threshold {
				orders = append(orders, o)
				ordercount++
			}
		}
	}
	return orders, nil
}
