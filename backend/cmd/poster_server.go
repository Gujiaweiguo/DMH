package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	postersDir := "/opt/data/posters"

	http.HandleFunc("/api/v1/posters/", func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/api/v1/posters/")

		if filename == "" || filename == "." {
			http.Error(w, "Invalid filename", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join(postersDir, filename)

		info, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		if info.IsDir() {
			http.Error(w, "Invalid filename", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeFile(w, r, filePath)
	})

	fmt.Println("Static file server listening on :8889")
	http.ListenAndServe(":8889", nil)
}
