package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type List struct {
	Name  string
	Owner string
	Items []ListItem
}

type Order struct {
	Name  string
	Items []OrderItem
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
	ioutil.WriteFile(l.Name+".json", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (l *List) Load(filename string) error {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &l)
	if err != nil {
		return err
	}
	return nil
}

func (o *Order) Save() error {
	b, err := json.MarshalIndent(&o, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile(o.Name+".json", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (o *Order) Load(filename string) error {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &o)
	if err != nil {
		return err
	}
	return nil
}

func LoadOrders() ([]Order, error) {
	var orders []Order
	files, err := ioutil.ReadDir("orders")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			o := Order{}
			err := o.Load("orders/" + f.Name())
			if err != nil {
				return nil, err
			}
			orders = append(orders, o)
		}
	}
	return orders, nil
}
