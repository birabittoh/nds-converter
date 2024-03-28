package main

import (
	"bytes"
	"embed"
	"net/http"
	"text/template"
)

const templatesDirectory = "templates/"

var (
	//go:embed templates/index.html
	templates     embed.FS
	indexTemplate = template.Must(template.ParseFS(templates, templatesDirectory+"index.html"))
)

func getHandler(w http.ResponseWriter) {
	buf := &bytes.Buffer{}
	err := indexTemplate.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}
