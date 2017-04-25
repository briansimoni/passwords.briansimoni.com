package main

import (
	"encoding/base64"
	"github.com/gorilla/securecookie"
	"io/ioutil"
)

func passwordGenerator() string {
	x := base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	return x
}

// need to check that the password is 32 bytes
func getMasterKey() (string, error){
	dat, err := ioutil.ReadFile("master-key")
	if err != nil {
		return "", err
	}
	return string(dat), nil
}
