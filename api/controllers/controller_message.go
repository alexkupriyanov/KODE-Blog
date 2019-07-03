package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/alexkupriyanov/KODE-Blog/api/models"
	"github.com/alexkupriyanov/KODE-Blog/api/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	//TODO: Change method in api documentation
	var message models.Message
	_ = r.ParseMultipartForm(10 * 1024 * 1024)
	message.Text = r.PostForm.Get("text")
	message.Link.Link = r.PostForm.Get("link")
	message.Author.Token = strings.Split(r.Header.Get("Authorization"), " ")[1]
	if message.Author.Token == "" {
		http.Error(w, "You are not authorized", http.StatusUnauthorized)
		return
	}
	resp := message.Create(r)
	if resp["status"] == false {
		utils.Respond(w, resp)
		return
	}
	e := json.NewEncoder(w).Encode(message.ToListModel())
	if e != nil {
		fmt.Println("System error: Can't create json object for:", message)
	}
}

func GetMessageList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := 0
	if len(vars["page"]) != 0 {
		page, _ = strconv.Atoi(vars["page"])
	}
	messages := models.GetMessageList(uint(page))
	e := json.NewEncoder(w).Encode(messages)
	if e != nil {
		fmt.Println("System error: Can't create json object for:", messages)
	}
}

func GetMessageDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id int
	if len(vars["id"]) != 0 {
		id, _ = strconv.Atoi(vars["id"])
	}
	if id == 0 {
		utils.Respond(w, utils.Message(false, "Message not found"))
		return
	}
	message := models.Message{Id: uint(id)}
	resp := message.Details()
	if resp["status"] == false {
		utils.Respond(w, resp)
		return
	}
	mes := message.ToDetailsModel()
	e := json.NewEncoder(w).Encode(mes)
	if e != nil {
		fmt.Println("System error: Can't create json object for:", message)
	}
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	vars := mux.Vars(r)
	var id int
	if len(vars["id"]) != 0 {
		id, _ = strconv.Atoi(vars["id"])
	}
	if id == 0 {
		utils.Respond(w, utils.Message(false, "Message not found"))
		return
	}
	message.Id = uint(id)
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	resp := message.Delete(token)
	if resp["status"] == false {
		utils.Respond(w, resp)
		return
	}
}

func Like(w http.ResponseWriter, r *http.Request) {
	var message models.Message
	vars := mux.Vars(r)
	var id int
	if len(vars["id"]) != 0 {
		id, _ = strconv.Atoi(vars["id"])
	}
	if id == 0 {
		utils.Respond(w, utils.Message(false, "Message not found"))
		return
	}
	message.Id = uint(id)
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	if token == "" {
		http.Error(w, "You are not authorized", http.StatusUnauthorized)
		return
	}
	resp := message.Like(token)
	if resp["status"] == false {
		utils.Respond(w, resp)
		return
	}
	e := json.NewEncoder(w).Encode(message.ToDetailsModel())
	if e != nil {
		fmt.Println("System error: Can't create json object for:", message)
	}
}
