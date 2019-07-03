package controllers

import (
	"KODE-Blog/api/models"
	"KODE-Blog/api/utils"
	"encoding/json"
	"net/http"
	"strings"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&user)
	response := user.Create()
	if response["status"] == false {
		utils.Respond(w, response)
		return
	}
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&user)
	response := user.Login()
	if response["status"] == false {
		//TODO: CHANGE RESPONSE HTTP STATUS CODE TO 403
		utils.Respond(w, response)
		return
	}
	utils.Respond(w, response)
	return
}

func Logout(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	user.Token = strings.Split(r.Header.Get("Authorization"), " ")[1]
	response := user.Logout()
	if response["status"] == false {
		//TODO: CHANGE RESPONSE HTTP STATUS CODE TO 403
		utils.Respond(w, response)
		return
	}
	utils.Respond(w, response)
	return
}
