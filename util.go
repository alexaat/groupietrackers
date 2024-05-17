package main

import (
	"encoding/json"
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
	display = make(map[string][]int)
	for _, artist := range bands {
		// artist/band
		artistNameDisplay := artist.Name + " - artist/band"
		addToMap(display, artistNameDisplay, artist.ID)

		// members
		for _, member := range artist.Members {
			memberNameDisplay := member + " - member"
			addToMap(display, memberNameDisplay, artist.ID)

		}
		// creation date
		creationDate := strconv.Itoa(artist.CreationDate)
		creationDateDisplay := creationDate + " - creation date"
		addToMap(display, creationDateDisplay, artist.ID)

		// first album date
		firstAlbumDate := artist.FirstAlbum
		firstAlbumDateDisplay := firstAlbumDate + " - first album date"
		addToMap(display, firstAlbumDateDisplay, artist.ID)

		// locations
		for key := range artist.Concerts {
			locationDisplay := key + " - location"
			addToMap(display, locationDisplay, artist.ID)
		}
	}
	return display
}

func fetchData(c chan Error) {
	data, err := getData(api)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(data), &groupies)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	// Get bands
	artists, err := getData(groupies.Artists)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(artists), &bands)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	// Get Relations
	isDataAvailable := true
	relat, err := getData(groupies.Relation)
	if err != nil {
		isDataAvailable = false
	}
	if isDataAvailable {
		err = json.Unmarshal([]byte(relat), &relations)
		if err != nil {
			isDataAvailable = false
		}
	}
	if len(relations.Index) == 0 {
		isDataAvailable = false
	}

	if isDataAvailable {
		// Save relations map to artists
		for _, item := range relations.Index {
			bands[item.ID-1].Concerts = formatLocations(item.DatesLocations)
		}
		c <- Error{Code: http.StatusOK, Message: ""}
		return
	}
	// Get Locations if not avalable
	locat, err := getData(groupies.Locations)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(locat), &locations)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	// Get Dates
	d, err := getData(groupies.Dates)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	err = json.Unmarshal([]byte(d), &dates)
	if err != nil {
		c <- make500Error(err.Error())
		return
	}
	// Construct map using locations and dates
	for index := range bands {
		m := constructConcerts(index + 1)
		bands[index].Concerts = m
	}

	c <- Error{Code: http.StatusOK, Message: ""}
}

func make500Error(message string) Error {
	return Error{Code: http.StatusInternalServerError, Message: "500 INTERNAL SERVER ERROR: " + message}
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
	for _, location := range locations.Index {
		if location.ID == id {
			locationsArr = location.Locations
			break
		}
	}

	datesArr := []string{}
	for _, dates := range dates.Index {
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
