package utils

import (
	"crypto/rand"
	"encoding/base64"
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

func GenerateVerificationToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
