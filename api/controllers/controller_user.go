package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/alexkupriyanov/KODE-Blog/api/models"
	"net/http"
	"strings"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&user)
	resp := user.Create()
	if resp["status"] == false {
		http.Error(w, fmt.Sprint(resp["message"]), http.StatusForbidden)
		return
	}
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&user)
	resp := user.Login()
	if resp["status"] == false {
		http.Error(w, fmt.Sprint(resp["message"]), http.StatusForbidden)
		return
	}
	_ = json.NewEncoder(w).Encode(resp["message"])
	return
}

func Logout(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	user.Token = strings.Split(r.Header.Get("Authorization"), " ")[1]
	resp := user.Logout()
	if resp["status"] == false {
		http.Error(w, fmt.Sprint(resp["message"]), http.StatusForbidden)
		return
	}
	_ = json.NewEncoder(w).Encode("You are logged out")
	return
}
