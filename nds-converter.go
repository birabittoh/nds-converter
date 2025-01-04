package main

import (
	_ "embed"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	MB                 = 1 << 20
	DefaultmaxFileSize = 2
	DefaultPort        = "3000"
)

var (
	//go:embed index.html
	indexPage     []byte
	maxFileSize   int64
	maxFileSizeMB int64
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		println("Could not load .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	fileSizeString := os.Getenv("MAX_SIZE_MB")
	maxFileSize, err = strconv.ParseInt(fileSizeString, 10, 64)
	if err != nil {
		maxFileSize = DefaultmaxFileSize
	}
	println("Max file size:", maxFileSize, "MB.")
	maxFileSizeMB = maxFileSize * MB

	http.HandleFunc("/", mainHandler)
	println("Listening on port " + port)
	http.ListenAndServe(":"+port, nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHandler(w)
		return
	case http.MethodPost:
		postHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter) {
	w.Write(indexPage)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(maxFileSizeMB)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Limit upload file size
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSizeMB)

	// Get the file from the request
	file, h, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file size
	if h.Size > maxFileSizeMB {
		http.Error(w, "File size exceeds "+strconv.FormatInt(maxFileSize, 10)+"MB", http.StatusBadRequest)
		return
	}

	// Convert file
	output, err := Convert(file, h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Write(output)
}
