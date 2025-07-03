package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// compile latex and display
func (s *server) HandlePreview(w http.ResponseWriter, r *http.Request) {

	const latexImg = "preview"

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

	latex := fmt.Sprintf(`\documentclass[preview]{standalone}
\usepackage{amsmath}
\begin{document}
%s
\end{document}`, currForm.Problem)

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
func (s *server) HandleSubmit(w http.ResponseWriter, r *http.Request) {

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

	err = s.storage.InsertProblem(currForm.Section, currForm.Difficulty, currForm.Problem)
	if err != nil {
		http.Error(w, "unable to write to db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Form submitted successfully"})
}

func (s *server) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, s.html)
}
