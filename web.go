package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const portNumber = ":8080"

func artistHandler(w http.ResponseWriter, r *http.Request) {
	c := make(chan Error)
	go fetchData(c)
	e := <-c
	if e.Code != http.StatusOK {
		errorHandler(w, e.Code, e.Message)
		return
	}
	path := r.URL.Path
	if path != "/" {

		t, e := template.ParseFiles("templates/artist.html")
		if e != nil {
			errorHandler(w, http.StatusNotFound, "404 NOT FOUND artist.html")
			return
		}

		ids := []int{}
		idsString := strings.TrimSpace(path[1:])
		idsSlice := strings.Fields(idsString)
		for _, idStr := range idsSlice {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				errorHandler(w, http.StatusBadRequest, "400 BAD REQUEST")
				return
			}
			ids = append(ids, id)
		}

		bandsToDysplay := []Artist{}
		for _, id := range ids {
			for _, artist := range bands {
				if artist.ID == id {
					bandsToDysplay = append(bandsToDysplay, artist)
				}
			}
		}

		if len(bandsToDysplay) == 0 {
			errorHandler(w, http.StatusBadRequest, "400 BAD REQUEST")
			return
		}
		e = t.Execute(w, bandsToDysplay)

		if e != nil {
			errorHandler(w, http.StatusInternalServerError, "500 INTERNAL SERVER ERROR")
			return
		}
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		errorHandler(w, http.StatusNotFound, "404 NOT FOUND")
	}

	display = createSearchMaster()
	searchObject = SearchObject{}
	searchObject.Artists = bands
	searchObject.Display = display
	err = t.Execute(w, searchObject)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, "500 SERVER ERROR")
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errorHandler(w, http.StatusBadRequest, "400 BAD REQUEST")
		return
	}

	searchRequest := strings.ToLower(strings.TrimSpace(r.FormValue("artist")))

	IDs := []int{}
	for key := range display {

		keyFormat := strings.ToLower(key)
		keyFormatNoSuffix := removeSuffixes(keyFormat)
		searchRequestNoSuffix := removeSuffixes(searchRequest)

		if strings.Contains(keyFormat, searchRequest) && strings.Contains(keyFormatNoSuffix, searchRequestNoSuffix) {
			if ids, ok := display[key]; ok {
				for _, id := range ids {
					if !contains(IDs, id) {
						IDs = append(IDs, id)
					}
				}
			}
		}
	}
	if len(IDs) > 0 {
		url := intsToUrl(IDs)
		http.Redirect(w, r, "/"+url, http.StatusFound)
	} else {
		t, e := template.ParseFiles("templates/notfound.html")
		if e != nil {
			errorHandler(w, http.StatusNotFound, "404 NOT FOUND notfound.html")
			return
		}
		err := t.Execute(w, searchRequest)
		if err != nil {
			errorHandler(w, http.StatusInternalServerError, "500 SERVER ERROR")
		}
	}
}

func errorHandler(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	template, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, message, code)
		return
	}
	err = template.Execute(w, Error{Message: message, Code: code})
	if err != nil {
		http.Error(w, message, code)
	}
}
