package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)


func HashPassword(p string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashPassword), nil
}


func ParseBody(r *http.Request, x interface{}){
	if b, err := ioutil.ReadAll(r.Body); err == nil{
		if err := json.Unmarshal([]byte(b), x); err != nil {
			return
		}
	}
}



func CheckPassordHash(p, h string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))

	return err == nil
}