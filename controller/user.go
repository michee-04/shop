package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/michee/e-commerce/model"
	"github.com/michee/e-commerce/utils"
)


func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	utils.ParseBody(r, user)
	u := user.CreateUser()

	res, _ := json.Marshal(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}


func GetUser(w http.ResponseWriter, r *http.Request) {
	u := model.GetUser()
	res, _ := json.Marshal(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}


func GetUserById(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	u, _ := model.GetUserById(userId)
	res, _ := json.Marshal(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}



func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	u := model.DeleteUser(userId)
	res, _ := json.Marshal(u)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

	utils.RespondWithJSON(w, http.StatusOK, "Delete user successful", map[string]string{"token": string(res)})
}



func UpddateUser(w http.ResponseWriter, r *http.Request){
	userUpdate := chi.URLParam(r, "userId")
	utils.ParseBody(r, userUpdate)	
}