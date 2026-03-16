package handler

import (
	"embed"
	models "groupietrackers/models"
	util "groupietrackers/utils"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func ArtistHandler(w http.ResponseWriter, r *http.Request, templates embed.FS) {

	c := make(chan models.Error)
	go util.FetchData(c)
	e := <-c
	if e.Code != http.StatusOK {
		ErrorHandler(w, e.Code, e.Message, templates)
		return
	}
	path := r.URL.Path
	if path != "/" {

		tmpl, err := template.ParseFS(templates, "templates/artist.html")
		if err != nil {
			ErrorHandler(w, http.StatusNotFound, "404 NOT FOUND artist.html", templates)
			return
		}

		ids := []int{}
		idsString := strings.TrimSpace(path[1:])
		idsSlice := strings.Fields(idsString)
		for _, idStr := range idsSlice {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				ErrorHandler(w, http.StatusBadRequest, "400 BAD REQUEST", templates)
				return
			}
			ids = append(ids, id)
		}

		bandsToDysplay := []models.Artist{}
		for _, id := range ids {
			for _, artist := range models.BandsInstance {
				if artist.ID == id {
					bandsToDysplay = append(bandsToDysplay, artist)
				}
			}
		}

		if len(bandsToDysplay) == 0 {
			ErrorHandler(w, http.StatusBadRequest, "400 BAD REQUEST", templates)
			return
		}

		err = tmpl.Execute(w, bandsToDysplay)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "500 INTERNAL SERVER ERROR", templates)
		}

		return
	}

	t, err := template.ParseFS(templates, "templates/index.html")
	if err != nil {
		ErrorHandler(w, http.StatusNotFound, "404 NOT FOUND", templates)
		return
	}

	models.DisplayInstance = util.CreateSearchMaster()
	models.SearchObjectInstance = models.SearchObject{}
	models.SearchObjectInstance.Artists = models.BandsInstance
	models.SearchObjectInstance.Display = models.DisplayInstance
	err = t.Execute(w, models.SearchObjectInstance)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "500 SERVER ERROR", templates)
	}

}

func SearchHandler(w http.ResponseWriter, r *http.Request, templates embed.FS) {
		if r.Method != "POST" {
			ErrorHandler(w, http.StatusBadRequest, "400 BAD REQUEST", templates)
			return
		}

		searchRequest := strings.ToLower(strings.TrimSpace(r.FormValue("artist")))

		IDs := []int{}
		for key := range models.DisplayInstance {

			keyFormat := strings.ToLower(key)
			keyFormatNoSuffix := util.RemoveSuffixes(keyFormat)
			searchRequestNoSuffix := util.RemoveSuffixes(searchRequest)

			if strings.Contains(keyFormat, searchRequest) && strings.Contains(keyFormatNoSuffix, searchRequestNoSuffix) {
				if ids, ok := models.DisplayInstance[key]; ok {
					for _, id := range ids {
						if !util.Contains(IDs, id) {
							IDs = append(IDs, id)
						}
					}
				}
			}
		}
		if len(IDs) > 0 {
			url := util.IntsToUrl(IDs)
			http.Redirect(w, r, "/"+url, http.StatusFound)
		} else {

			tmpl, err := template.ParseFS(templates, "templates/notfound.html")
			if err != nil {
				ErrorHandler(w, http.StatusNotFound, "404 NOT FOUND notfound.html", templates)
				return
			}
			err = tmpl.Execute(w, searchRequest)
			if err != nil {
				ErrorHandler(w, http.StatusInternalServerError, "500 SERVER ERROR", templates)
			}
		}
	
}

func ErrorHandler(w http.ResponseWriter, code int, message string, templates embed.FS) {
	w.WriteHeader(code)
	tmpl, err := template.ParseFS(templates, "templates/error.html")
	if err != nil {
		http.Error(w, message+"Template not found", code)
		return
	}
	err = tmpl.Execute(w, models.Error{Message: message, Code: code})
	if err != nil {
		http.Error(w, message+"Template Execute Fail", code)
	}

	// w.WriteHeader(code)
	// template, err := template.ParseFiles("templates/error.html")
	// if err != nil {
	// 	http.Error(w, message, code)
	// 	return
	// }
	// err = template.Execute(w, Error{Message: message, Code: code})
	// if err != nil {
	// 	http.Error(w, message, code)
	// }

}
