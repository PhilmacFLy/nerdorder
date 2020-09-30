package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type User struct {
	Username string
	Password string
	Email    string
}

func (u *User) loadLists() ([]List, error) {
	var lists []List
	files, err := ioutil.ReadDir("lists/" + u.Username)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			l := List{}
			l.Owner = u.Username
			l.Name = StripExt(f.Name())
			err := l.Load()
			if err != nil {
				return nil, err
			}
			lists = append(lists, l)
		}
	}
	return lists, nil
}

func (u *User) Load() error {
	body, err := ioutil.ReadFile("users/" + u.Username + ".json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Register() error {
	if _, err := os.Stat("users/" + u.Username + ".json"); err == nil {
		return errors.New("User already exists")
	}
	return u.Save()
}

func (u *User) Save() error {
	b, err := json.MarshalIndent(&u, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile("users/"+u.Username+".json", b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Remove() error {
	return os.Remove("users/" + u.Username + ".json")
}
