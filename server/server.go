package server

import (
	"database/sql"
	_ "embed"
	"strings"
)

//go:embed view/index.html
var htmlTemplate string

//go:embed view/script.js
var scriptJS string

//go:embed view/styles.css
var stylesCSS string

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
