package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const imgDir = "images"

//go:embed view/index.html
var htmlTemplate string

//go:embed view/script.js
var scriptJS string

var storage struct {
	jsonFile, sqlFile, plainFile *os.File
	sqlDb                        *sql.DB
}

// Init the file vars
func init() {

	const storageDir = "storage"

	dirToMake := []string{storageDir, imgDir}

	for _, dir := range dirToMake {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	filesToMake := map[string]**os.File{
		"data.json": &storage.jsonFile,
		"data.sql":  &storage.sqlFile,
		"data.txt":  &storage.plainFile,
	}

	for name, file := range filesToMake {
		filePath := filepath.Join(storageDir, name)
		temp, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		*file = temp
	}

	db, err := sql.Open("sqlite3", filepath.Join(storageDir, "data.db"))
	if err != nil {
		log.Fatal(err)
	}

	const createTable = `
	CREATE TABLE IF NOT EXISTS problems (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		section TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		problem TEXT NOT NULL
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	storage.sqlDb = db

}

// starting point of the program
func main() {

	defer func() {
		storage.sqlFile.Close()
		storage.jsonFile.Close()
		storage.plainFile.Close()
		storage.sqlDb.Close()
	}()

	http.HandleFunc("/preview", previewHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		finalHTML := strings.Replace(htmlTemplate, "{{SCRIPT}}", scriptJS, 1)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, finalHTML)
	})

	fmt.Println("http://localhost:8080/")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

// compile latex and display
func previewHandler(w http.ResponseWriter, r *http.Request) {

	const latexImg = "preview"

	const latexFormat = `\documentclass[preview]{standalone}
\usepackage{amsmath}
\begin{document}
%s
\end{document}`

	defer func() {
		fileTypesToDelete := []string{".aux", ".log", ".pdf", ".tex"}
		for _, fileType := range fileTypesToDelete {
			os.Remove(latexImg + fileType)
		}
	}()

	var currForm struct {
		Problem string `json:"problem"`
	}

	err := json.NewDecoder(r.Body).Decode(&currForm)
	if err != nil || currForm.Problem == "" {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	latex := fmt.Sprintf(latexFormat, currForm.Problem)
	err = os.WriteFile(latexImg+".tex", []byte(latex), 0644)
	if err != nil {
		http.Error(w, "unable to write the latex to file", http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", latexImg+".tex")
	cmd.Dir, cmd.Stdout, cmd.Stderr = ".", os.Stdout, os.Stderr

	err = cmd.Run()
	if err != nil {
		http.Error(w, "Error compiling LaTeX", http.StatusInternalServerError)
		return
	}

	pngPath := filepath.Join(imgDir, latexImg+".png")

	cmd = exec.Command("magick", "-density", "300", latexImg+".pdf", "-quality", "90", pngPath)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	err = cmd.Run()
	if err != nil {
		http.Error(w, "ImageMagick conversion failed", http.StatusInternalServerError)
		return
	}

	imgData, err := os.ReadFile(pngPath)
	if err != nil {
		http.Error(w, "Error reading image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(imgData)
}

// save as both json and sql insert queries
func submitHandler(w http.ResponseWriter, r *http.Request) {

	const insertSqliteFormat = `INSERT INTO problems (section, difficulty, problem) VALUES (?, ?, ?);`

	const jsonFormat = `
  {
    "section": "%s",
    "difficulty": "%s",
    "problem": "%s"
  },`

	const sqlFormat = `
INSERT INTO problems (course, difficulty, problem_text)
VALUES ('%s', '%s', '%s');`

	const plainFormat = "%s\n%s\n%s\n"

	var currForm struct {
		Problem    string `json:"problem"`
		Section    string `json:"section"`
		Difficulty string `json:"difficulty"`
	}

	err := json.NewDecoder(r.Body).Decode(&currForm)
	if err != nil || currForm.Problem == "" || currForm.Section == "" || currForm.Difficulty == "" {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	filesToWrite := map[string]*os.File{
		jsonFormat:  storage.jsonFile,
		sqlFormat:   storage.sqlFile,
		plainFormat: storage.plainFile,
	}

	for fileFormat, file := range filesToWrite {
		_, err = fmt.Fprintf(file, fileFormat, currForm.Section, currForm.Difficulty, currForm.Problem)
		if err != nil {
			http.Error(w, "Unable to write to files", http.StatusInternalServerError)
			return
		}
	}

	_, err = storage.sqlDb.Exec(insertSqliteFormat, currForm.Section, currForm.Difficulty, currForm.Problem)
	if err != nil {
		http.Error(w, "unable to write to db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Form submitted successfully"})
}
