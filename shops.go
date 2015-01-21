package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type Shop struct {
	Name      string
	Threshold float64
	Website   string
}

func (s *Shop) Save() error {
	b, err := json.MarshalIndent(&s, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile("shops/"+s.Name+".json", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *Shop) Load() error {
	body, err := ioutil.ReadFile("shops/" + s.Name + ".json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return err
	}
	return nil
}

func LoadShops() ([]Shop, error) {
	var shops []Shop
	files, err := ioutil.ReadDir("shops")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			s := Shop{}
			s.Name = StripExt(f.Name())
			err := s.Load()
			if err != nil {
				return nil, err
			}
			shops = append(shops, s)
		}
	}
	return shops, nil
}
