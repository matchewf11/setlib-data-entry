package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// save as both json and sql insert queries
func (s *server) HandleSubmit(w http.ResponseWriter, r *http.Request) {

	var currForm struct {
		Problem    string `json:"problem"`
		Section    string `json:"section"`
		Difficulty string `json:"difficulty"`
		Username   string `json:"username"`
		Type       string `json:"type"`
		Subject    string `json:"subject"`
	}

	err := json.NewDecoder(r.Body).Decode(&currForm)
	if err != nil {
		http.Error(w, "unable to write to db", http.StatusInternalServerError)
		return
	}

	err = s.storage.InsertProblem(currForm.Section, currForm.Difficulty, currForm.Problem, currForm.Username, currForm.Subject, currForm.Type)
	if err != nil {
		http.Error(w, "unable to write to db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Form submitted successfully"})
}

// handle the get request
func (s *server) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, s.html)
}

// compile latex and display
func (s *server) HandlePreview(w http.ResponseWriter, r *http.Request) {

	var currForm struct {
		Problem string `json:"problem"`
	}

	err := json.NewDecoder(r.Body).Decode(&currForm)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	imgData, err := getPng(currForm.Problem)
	if err != nil {
		http.Error(w, "Could not make png", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(imgData)
}

func getPng(problem string) ([]byte, error) {
	const latexImg = "preview"

	defer func() {
		fileTypesToDelete := []string{".aux", ".log", ".pdf", ".tex"}
		for _, fileType := range fileTypesToDelete {
			os.Remove(latexImg + fileType)
		}
		os.Remove("texput.log")
	}()

	latex := fmt.Sprintf(`\documentclass[preview]{standalone}
\usepackage{amsmath}
\usepackage{graphicx}
\begin{document}
%s
\end{document}`, problem)

	err := os.WriteFile(latexImg+".tex", []byte(latex), 0644)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", latexImg+".tex")
	cmd.Dir, cmd.Stdout, cmd.Stderr = ".", os.Stdout, os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	pngPath := filepath.Join(imgDir, latexImg+".png")

	cmd = exec.Command("magick", "-density", "300", latexImg+".pdf", "-quality", "90", pngPath)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	imgData, err := os.ReadFile(pngPath)
	if err != nil {
		return nil, err
	}

	return imgData, nil
}
