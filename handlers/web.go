package handler

import (
	"embed"
	"fmt"
	models "groupietrackers/models"

	//util "groupietrackers/utils"
	"html/template"
	"net/http"
	//"strconv"
	//"strings"
)

//const portNumber = ":8080"

func ArtistHandler(w http.ResponseWriter, r *http.Request, templates embed.FS) {

	tmpl, err := template.ParseFS(templates, "templates/index.html")
	if err != nil {
		http.Error(w, "Template Parse Fail", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template Execute Fail", http.StatusInternalServerError)
	}

	/*
		c := make(chan models.Error)
		go util.FetchData(c)
		e := <-c
		if e.Code != http.StatusOK {
			ErrorHandler(w, e.Code, e.Message, templates)
			return
		}
		path := r.URL.Path
		if path != "/" {

			tmpl, err := template.ParseFS(templates, "templates/index.html")
			if err != nil {
				ErrorHandler(w, http.StatusNotFound, "404 NOT FOUND artist.html", templates)
				return
			}

			// t, e := template.ParseFiles("templates/artist.html")
			// if e != nil {
			// 	ErrorHandler(w, http.StatusNotFound, "404 NOT FOUND artist.html", templates)
			// 	return
			// }

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

			// e = t.Execute(w, bandsToDysplay)

			// if e != nil {
			// 	ErrorHandler(w, http.StatusInternalServerError, "500 INTERNAL SERVER ERROR", templates)
			// 	return
			// }
			return
		}

		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			ErrorHandler(w, http.StatusNotFound, "404 NOT FOUND", templates)
		}

		models.DisplayInstance = util.CreateSearchMaster()
		models.SearchObjectInstance = models.SearchObject{}
		models.SearchObjectInstance.Artists = models.BandsInstance
		models.SearchObjectInstance.Display = models.DisplayInstance
		err = t.Execute(w, models.SearchObjectInstance)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError, "500 SERVER ERROR", templates)
		}
	*/

}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h3>Search Handler</h3>")
	/*
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
	*/
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
