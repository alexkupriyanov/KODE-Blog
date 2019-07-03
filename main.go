package main

import (
	"github.com/alexkupriyanov/KODE-Blog/api"
	"log"
	"net/http"
)

func main() {
	log.Printf("Server started")

	router := routes.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
