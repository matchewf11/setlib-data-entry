package main

import (
	"fmt"
	"log"
	"net/http"
	"setlib-data-entry/server"
	db "setlib-data-entry/storage"
)

// starting point of the program
func main() {

	db, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	svr := server.New(db)

	http.HandleFunc("/", svr.HandleGet)
	http.HandleFunc("/preview", svr.HandlePreview)
	http.HandleFunc("/submit", svr.HandleSubmit)

	fmt.Println("http://localhost:8080/")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
