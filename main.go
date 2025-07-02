package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const (
	jsonFileName = "storage/jsonData.json"
	sqlFileName  = "storage/sqlData.sql"
	port         = ":8080"
)

//go:embed index.html
var htmlTemplate string

var jsonFile, sqlFile *os.File

// Init the file vars
func init() {

	err := os.MkdirAll("storage", 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll("images", 0755)
	if err != nil {
		log.Fatal(err)
	}

	json, err := os.OpenFile(jsonFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	sql, err := os.OpenFile(sqlFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	jsonFile, sqlFile = json, sql

}

// starting point of the program
func main() {

	defer func() {
		sqlFile.Close()
		jsonFile.Close()
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, htmlTemplate) })
	http.HandleFunc("/preview", previewHandler)
	http.HandleFunc("/submit", submitHandler)

	fmt.Printf("http://localhost%s/\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}

}

// compile latex and display
func previewHandler(w http.ResponseWriter, r *http.Request) {

	type formInfo struct {
		Problem string `json:"problem"`
	}
	var currForm formInfo

	err := json.NewDecoder(r.Body).Decode(&currForm)
	if err != nil || currForm.Problem == "" {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	latex := `\documentclass[preview]{standalone}
\usepackage{amsmath}
\begin{document}
` + currForm.Problem + `
\end{document}`

	err = os.WriteFile("preview.tex", []byte(latex), 0644)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "unable to write the latex to file", http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "preview.tex")
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error compiling LaTeX", http.StatusInternalServerError)
		return
	}

	cmd = exec.Command("magick", "-density", "300", "preview.pdf", "-quality", "90", "images/preview.png")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		http.Error(w, "ImageMagick conversion failed", http.StatusInternalServerError)
		return
	}

	for _, f := range []string{"preview.aux", "preview.log", "preview.pdf", "preview.tex"} {
		os.Remove(f)
	}

	imgData, err := os.ReadFile("images/preview.png")
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

	type formInfo struct {
		Problem    string `json:"problem"`
		Section    string `json:"section"`
		Difficulty string `json:"difficulty"`
	}

	var currForm formInfo
	err := json.NewDecoder(r.Body).Decode(&currForm)
	if err != nil || currForm.Problem == "" || currForm.Section == "" || currForm.Difficulty == "" {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = fmt.Fprintf(jsonFile, `
{
 "section": "%s",
 "difficulty": "%s",
  "problem": "%s"
},`, currForm.Section, currForm.Difficulty, currForm.Problem)
	if err != nil {
		http.Error(w, "Unable to write to files", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintf(sqlFile, `
INSERT INTO problems (course, difficulty, text)
VALUES (%s, %s, %s);`, currForm.Section, currForm.Difficulty, currForm.Problem)
	if err != nil {
		http.Error(w, "Unable to write to files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Form submitted successfully"})
}
