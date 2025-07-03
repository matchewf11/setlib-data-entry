package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// need to validate latex
// need to validate before saving as json
// need to validate before doing the sql

const (
	storageDir    = "storage"
	imgDir        = "images"
	latexImg      = "preview"
	jsonFileName  = "data.json"
	sqlFileName   = "data.sql"
	plainFileName = "data.txt"
	port          = ":8080"
	jsonFormat    = `
  {
    "section": "%s",
    "difficulty": "%s",
    "problem": "%s"
  },`
	latexFormat = `\documentclass[preview]{standalone}
\usepackage{amsmath}
\begin{document}
%s
\end{document}`
	sqlFormat = `
INSERT INTO problems (course, difficulty, problem_text)
VALUES (%s, %s, %s);`
	plainFormat = "%s\n%s\n%s\n"
)

//go:embed index.html
var htmlTemplate string

var jsonFile, sqlFile, plainFile *os.File

// Init the file vars
func init() {

	dirToMake := []string{storageDir, imgDir}

	for _, dir := range dirToMake {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	filesToMake := map[string]**os.File{
		jsonFileName:  &jsonFile,
		sqlFileName:   &sqlFile,
		plainFileName: &plainFile,
	}

	for name, file := range filesToMake {
		filePath := filepath.Join(storageDir, name)
		temp, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		*file = temp
	}

}

// starting point of the program
func main() {

	defer func() {
		sqlFile.Close()
		jsonFile.Close()
		plainFile.Close()
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, htmlTemplate) })
	http.HandleFunc("/preview", previewHandler)
	http.HandleFunc("/submit", submitHandler)

	fmt.Printf("http://localhost%s/\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}

}

// compile latex and display
func previewHandler(w http.ResponseWriter, r *http.Request) {

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
		jsonFormat:  jsonFile,
		sqlFormat:   sqlFile,
		plainFormat: plainFile,
	}

	for fileFormat, file := range filesToWrite {
		_, err = fmt.Fprintf(file, fileFormat, currForm.Section, currForm.Difficulty, currForm.Problem)
		if err != nil {
			http.Error(w, "Unable to write to files", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Form submitted successfully"})
}
