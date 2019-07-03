package controllers

import (
	"encoding/json"
	"github.com/alexkupriyanov/KODE-Blog/api/models"
	"net/http"
	"strings"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&user)
	err := user.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&user)
	token, err := user.Login()
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var result map[string] interface{}
	result["token"] = token
	_ = json.NewEncoder(w).Encode(result)
	return
}

func Logout(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	user.Token = strings.Split(r.Header.Get("Authorization"), " ")[1]
	err := user.Logout()
	if err!=nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	_ = json.NewEncoder(w).Encode("You are logged out")
	return
}
