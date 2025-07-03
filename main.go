package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"setlib-data-entry/server"
	db "setlib-data-entry/storage"

	_ "github.com/mattn/go-sqlite3"
)

// Init the file vars
func init() {
	err := os.MkdirAll(server.ImgDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

// starting point of the program
func main() {

	db, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	svr := server.New(db)

	http.HandleFunc("/preview", svr.HandlePreview)
	http.HandleFunc("/submit", svr.HandleSubmit)
	http.HandleFunc("/", svr.HandleGet)

	fmt.Println("http://localhost:8080/")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
