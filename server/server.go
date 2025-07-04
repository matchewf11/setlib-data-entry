package server

import (
	_ "embed"
	"log"
	"os"
	"setlib-data-entry/storage"
	"strings"
)

const imgDir = "images"

//go:embed view/index.html
var htmlTemplate string

//go:embed view/script.js
var scriptJS string

//go:embed view/styles.css
var stylesCSS string

// Init the file vars
func init() {
	err := os.MkdirAll(imgDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

type server struct {
	storage *storage.Storage
	html    string
}

func New(storage *storage.Storage) *server {
	finalHTML := strings.Replace(htmlTemplate, "/*SCRIPT*/", scriptJS, 1)
	finalHTML = strings.Replace(finalHTML, "/*STYLE*/", stylesCSS, 1)
	return &server{
		storage: storage,
		html:    finalHTML,
	}
}
