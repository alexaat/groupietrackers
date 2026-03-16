package handler

import (
	"embed"
	handlers "groupietrackers/handlers"
	"net/http"
)

//go:embed templates/*
var templates embed.FS

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/search":
		handlers.SearchHandler(w, r, templates)
	default:
		handlers.ArtistHandler(w, r, templates)
	}
}
