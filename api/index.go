package handler

import (
	"embed"
	"net/http"
	handlers "groupietrackers/handlers"
)

//go:embed templates/*
var templates embed.FS

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/" : handlers.ArtistHandler(w, r)
	case "/search" : handlers.SearchHandler(w, r)
	default: handlers.ErrorHandler(w, http.StatusNotFound, "404 Not Found")
	}
}
