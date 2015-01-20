package main

import (
	"encoding/json"
	"io/ioutil"
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
