package handler

import (
	"encoding/json"
	api "groupietrackers"
	models "groupietrackers/models"
	services "groupietrackers/services"
	"net/http"
	"strconv"
	"strings"
)

func formatLocation(location string) string {
	location = strings.Replace(location, "-", ", ", -1)
	location = strings.Replace(location, "_", " ", -1)
	location = strings.Title(location)
	location = strings.Replace(location, ", Usa", ", USA", -1)
	location = strings.Replace(location, ", Uk", ", UK", -1)
	return location
}

func formatLocations(data map[string][]string) map[string][]string {
	result := make(map[string][]string)
	for key, value := range data {
		result[formatLocation(key)] = value
	}
	return result
}

func createSearchMaster() map[string][]int {
	api.DisplayInstance = make(map[string][]int)
	for _, artist := range api.BandsInstance {
		// artist/band
		artistNameDisplay := artist.Name + " - artist/band"
		addToMap(api.DisplayInstance, artistNameDisplay, artist.ID)

		// members
		for _, member := range artist.Members {
			memberNameDisplay := member + " - member"
			addToMap(api.DisplayInstance, memberNameDisplay, artist.ID)

		}
		// creation date
		creationDate := strconv.Itoa(artist.CreationDate)
		creationDateDisplay := creationDate + " - creation date"
		addToMap(api.DisplayInstance, creationDateDisplay, artist.ID)

		// first album date
		firstAlbumDate := artist.FirstAlbum
		firstAlbumDateDisplay := firstAlbumDate + " - first album date"
		addToMap(api.DisplayInstance, firstAlbumDateDisplay, artist.ID)

		// locations
		for key := range artist.Concerts {
			locationDisplay := key + " - location"
			addToMap(api.DisplayInstance, locationDisplay, artist.ID)
		}
	}
	return api.DisplayInstance
}

func fetchData(c chan models.Error) {
	data, err := services.GetData(services.API)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(data), &api.GroupiesInstance)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	// Get bands
	artists, err := services.GetData(api.GroupiesInstance.Artists)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(artists), &api.BandsInstance)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	// Get Relations
	isDataAvailable := true
	relat, err := services.GetData(api.GroupiesInstance.Relation)
	if err != nil {
		isDataAvailable = false
	}
	if isDataAvailable {
		err = json.Unmarshal([]byte(relat), &api.RelationsInstance)
		if err != nil {
			isDataAvailable = false
		}
	}
	if len(api.RelationsInstance.Index) == 0 {
		isDataAvailable = false
	}

	if isDataAvailable {
		// Save relations map to artists
		for _, item := range api.RelationsInstance.Index {
			api.BandsInstance[item.ID-1].Concerts = formatLocations(item.DatesLocations)
		}
		c <- models.Error{Code: http.StatusOK, Message: ""}
		return
	}
	// Get Locations if not avalable
	locat, err := services.GetData(api.GroupiesInstance.Locations)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(locat), &api.LocationsInstance)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	// Get Dates
	d, err := services.GetData(api.GroupiesInstance.Dates)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(d), &api.DatesInstance)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	// Construct map using locations and dates
	for index := range api.BandsInstance {
		m := constructConcerts(index + 1)
		api.BandsInstance[index].Concerts = m
	}

	c <- models.Error{Code: http.StatusOK, Message: ""}
}

func make500Error(message string) models.Error {
	return models.Error{Code: http.StatusInternalServerError, Message: "500 INTERNAL SERVER ERROR: " + message}
}

func intsToUrl(data []int) string {
	url := ""
	for _, item := range data {
		idStr := strconv.Itoa(item)
		url += idStr + " "
	}
	return url
}

func addToMap(myMap map[string][]int, key string, value int) {
	if ids, ok := myMap[key]; ok {
		if !contains(ids, value) {
			ids = append(ids, value)
			myMap[key] = ids
		}
	} else {
		myMap[key] = []int{value}
	}
}

func contains(slice []int, value int) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func constructConcerts(id int) map[string][]string {
	result := make(map[string][]string)

	locationsArr := []string{}
	for _, location := range api.LocationsInstance.Index {
		if location.ID == id {
			locationsArr = location.Locations
			break
		}
	}

	datesArr := []string{}
	for _, dates := range api.DatesInstance.Index {
		if dates.ID == id {
			datesArr = dates.Dates
		}
	}
	datesString := strings.Join(datesArr, " ")
	datesByStar := strings.Split(datesString, "*")[1:]

	if len(datesByStar) != len(locationsArr) {
		return result
	}
	for i := 0; i < len(datesByStar); i++ {
		dates := strings.Split(datesByStar[i], " ")
		datesAdj := []string{}
		for j := 0; j < len(dates); j++ {
			date := strings.TrimSpace(dates[j])
			if date != "" {
				datesAdj = append(datesAdj, date)
			}
		}

		result[formatLocation(locationsArr[i])] = datesAdj
	}
	return result
}

func removeSuffixes(data string) string {
	data = strings.TrimSuffix(data, "- artist/band")
	data = strings.TrimSuffix(data, "- member")
	data = strings.TrimSuffix(data, "- creation date")
	data = strings.TrimSuffix(data, "- first album date")
	data = strings.TrimSuffix(data, "- location")
	return strings.TrimSpace(data)
}
