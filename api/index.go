package handler

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed templates/*
var templates embed.FS

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Main Page</h1>")
}
