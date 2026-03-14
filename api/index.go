package handler

import (
	"embed"
	"net/http"
	handlers "groupietrackers/handlers"
)

//go:embed templates/*
var templates embed.FS

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		handlers.ArtistHandler(w, r)
	}
}
