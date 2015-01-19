package main

import (
	"io/ioutil"
	"path/filepath"
)

type User struct {
	Username string
	Password string
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
			err := l.Load("lists/" + u.Username + "/" + f.Name())
			if err != nil {
				return nil, err
			}
			lists = append(lists, l)
		}
	}
	return lists, nil
}
