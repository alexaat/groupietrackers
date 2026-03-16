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
	// case "/":
	// 	handlers.ArtistHandler(w, r, templates)
	case "/search":
		handlers.SearchHandler(w, r, templates)
	// default:
	// 	handlers.ErrorHandler(w, http.StatusNotFound, "404 Not Found", templates)
	// }
	default:
		handlers.ArtistHandler(w, r, templates)
	}
}
