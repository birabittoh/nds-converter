package main

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

type ConverterFunction func([]byte) []byte

const (
	MB                 = 1 << 20
	DefaultMaxFileSize = 2
	TrimBytes          = 122
	AdditionalBytesHex = "7C3C2D2D536E69702061626F7665206865726520746F2063726561746520612072617720736176206279206578636C7564696E672074686973204465536D754D4520736176656461746120666F6F7465723A0000010000000100030000000200000000000100000000007C2D4445534D554D4520534156452D7C"
)

var (
	MaxFileSize        int64
	AdditionalBytes    []byte
	ConverterFunctions = map[string]ConverterFunction{
		".dsv": DsvToSav,
		".sav": SavToDsv,
	}
)

func DsvToSav(input []byte) []byte {
	return input[:len(input)-TrimBytes]
}

func SavToDsv(input []byte) []byte {
	return append(input, AdditionalBytes...)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize)

	// Get the file from the request
	file, h, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file size
	if h.Size > MaxFileSize {
		http.Error(w, "File size exceeds 1MB", http.StatusBadRequest)
		return
	}

	var outputContent []byte

	// Check file extension
	fileExt := filepath.Ext(h.Filename)
	converterFunction, exists := ConverterFunctions[fileExt]
	if !exists {
		http.Error(w, "Invalid file format, only .dsv and .sav files are allowed", http.StatusBadRequest)
		return
	}

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file content", http.StatusInternalServerError)
		return
	}
	outputContent = converterFunction(content)

	w.Header().Set("Content-Length", strconv.Itoa(len(outputContent)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(outputContent)
}
