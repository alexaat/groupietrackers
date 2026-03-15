package handler

import (
	"embed"
	"fmt"
	handlers "groupietrackers/handlers"
	models "groupietrackers/models"
	"net/http"
)

var (
	GroupiesInstance     models.Groupies
	BandsInstance        []models.Artist
	LocationsInstance    models.Locations
	DatesInstance        models.Dates
	RelationsInstance    models.Relations
	SearchObjectInstance models.SearchObject
	DisplayInstance      map[string][]int
)

//go:embed templates/*
var templates embed.FS

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		//handlers.ArtistHandler(w, r, templates)
		fmt.Fprintf(w, "<h1>Main Page</h1>")

	case "/search":
		handlers.SearchHandler(w, r)
	default:
		handlers.ErrorHandler(w, http.StatusNotFound, "404 Not Found", templates)
	}
}
