package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
)

type ConverterFunction func([]byte) []byte

const (
	trimBytes          = 122
	additionalBytesHex = "7C3C2D2D536E69702061626F7665206865726520746F2063726561746520612072617720736176206279206578636C7564696E672074686973204465536D754D4520736176656461746120666F6F7465723A0000010000000100030000000200000000000100000000007C2D4445534D554D4520534156452D7C"
)

var (
	additionalBytes    = MustDecode(hex.DecodeString(additionalBytesHex))
	converterFunctions = map[string]ConverterFunction{
		".dsv": DsvToSav,
		".sav": SavToDsv,
	}
)

func MustDecode(content []byte, err error) []byte {
	if err != nil {
		panic(42)
	}
	return content
}

func DsvToSav(input []byte) []byte {
	return input[:len(input)-trimBytes]
}

func SavToDsv(input []byte) []byte {
	return append(input, additionalBytes...)
}

func Convert(file multipart.File, h *multipart.FileHeader) ([]byte, error) {
	// Check file extension
	fileExt := filepath.Ext(h.Filename)
	converterFunction, exists := converterFunctions[fileExt]
	if !exists {
		return nil, fmt.Errorf("invalid file format: only .dsv and .sav files are allowed")
	}

	// Read file contents
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Return converted file
	return converterFunction(content), nil
}
