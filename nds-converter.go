package main

import (
	"encoding/hex"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	DefaultPort = "1111"
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
	MaxFileSize, err = strconv.ParseInt(fileSizeString, 10, 64)
	if err != nil {
		MaxFileSize = DefaultMaxFileSize
	}
	MaxFileSize *= MB

	AdditionalBytes, _ = hex.DecodeString(AdditionalBytesHex)

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
		http.Error(w, "Method not supported", http.StatusBadRequest)
	}
}
