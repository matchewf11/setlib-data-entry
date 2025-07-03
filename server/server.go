package server

import (
	"database/sql"
	_ "embed"
	"log"
	"os"
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
	db   *sql.DB
	html string
}

func New(db *sql.DB) *server {

	finalHTML := strings.Replace(htmlTemplate, "<!--SCRIPT-->", "<script>"+scriptJS+"</script>", 1)
	finalHTML = strings.Replace(finalHTML, "<!--STYLES-->", "<style>"+stylesCSS+"</style>", 1)

	return &server{
		db:   db,
		html: finalHTML,
	}
}
