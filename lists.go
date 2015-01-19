package main

import (
	"encoding/json"
	"io/ioutil"
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
