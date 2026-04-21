package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	validIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	maxNoteSize  int64 = 1048576
	readOnly     bool  = false
)

func init() {
	if v := os.Getenv("MAX_NOTE_SIZE"); v != "" {
		if size, err := strconv.ParseInt(v, 10, 64); err == nil && size > 0 {
			maxNoteSize = size
		}
	}
	if v := os.Getenv("READ_ONLY"); v == "true" || v == "1" {
		readOnly = true
	}
}

func registerHandlers(dataDir string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", serveIndex)
	mux.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		listNotes(w, r, dataDir)
	})
	mux.HandleFunc("/note/", func(w http.ResponseWriter, r *http.Request) {
		handleNote(w, r, dataDir)
	})

	http.Handle("/", mux)
}

func handleNote(w http.ResponseWriter, r *http.Request, dataDir string) {
	id := strings.TrimPrefix(r.URL.Path, "/note/")
	if id == "" || !validIDRegex.MatchString(id) {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		readNote(w, r, dataDir, id)
	case http.MethodPost:
		writeNote(w, r, dataDir, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(indexHTML))
}

func readNote(w http.ResponseWriter, r *http.Request, dataDir, id string) {
	path := filepath.Join(dataDir, id+".txt")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Printf("Error reading note %s: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(content)
}

func writeNote(w http.ResponseWriter, r *http.Request, dataDir, id string) {
	if readOnly {
		http.Error(w, "Read only mode", http.StatusForbidden)
		return
	}

	content, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxNoteSize+1))
	if err != nil {
		http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	if int64(len(content)) > maxNoteSize {
		http.Error(w, "Note size exceeds limit", http.StatusRequestEntityTooLarge)
		return
	}

	path := filepath.Join(dataDir, id+".txt")
	if err := os.WriteFile(path, content, 0644); err != nil {
		log.Printf("Error writing note %s: %v", id, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func listNotes(w http.ResponseWriter, r *http.Request, dataDir string) {
	files, err := os.ReadDir(dataDir)
	if err != nil {
		log.Printf("Error listing notes: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var notes []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			notes = append(notes, strings.TrimSuffix(file.Name(), ".txt"))
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(notes)
}
