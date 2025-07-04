package main

import (
	"fmt"
	"log"
	"net/http"
	"setlib-data-entry/server"
	"setlib-data-entry/storage"
)

// starting point of the program
func main() {

	storage, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	svr := server.New(storage)

	http.HandleFunc("/", svr.HandleGet)
	http.HandleFunc("/preview", svr.HandlePreview)
	http.HandleFunc("/submit", svr.HandleSubmit)

	fmt.Println("http://localhost:8080/")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
