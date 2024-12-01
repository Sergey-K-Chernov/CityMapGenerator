package main

import (
	"fmt"
	"net/http"
	"html/template"
	"strconv"
	gen "chirrwick.com/projects/city/generator"
	"chirrwick.com/projects/city/city_map"
)

type RoadsResponse struct{
	Default bool
	Error string
	Image string
	MinR float64
	MaxR float64
	NCenters int
	Branching int
}

func readRoadsParams(r *http.Request) (bool, RoadsResponse, gen.InitialValuesRoads, city_map.Map) {
	resp := RoadsResponse{Error: ""}
	resp.Default = false
	var city_map city_map.Map

	fmt.Println("Prepare initials")
	var initials gen.InitialValuesRoads

	fmt.Println("ReadMap")

	city_map = getMapFromCookies(r)
	fmt.Println("Ok. Read min r")

	min_r, err := strconv.ParseFloat(r.FormValue("min_r"), 32)
        if err != nil {
		resp.Error = "Cannot read min radius"
		return false, resp, initials, city_map
	}
	resp.MinR = min_r

	fmt.Println("Ok. Read max r")

        max_r, err := strconv.ParseFloat(r.FormValue("max_r"), 32)
        if err != nil {
		resp.Error = "Cannot read max radius"
                return false, resp, initials, city_map
	}
	resp.MaxR = max_r

	fmt.Println("Ok. Read n centers")

	n_centers, err := strconv.ParseInt(r.FormValue("n_centers"), 10, 32)
        if err != nil {
                resp.Error = "Cannot read centers"
                return false, resp, initials, city_map
        }
	resp.NCenters = int(n_centers)

	fmt.Println("Ok. Read branching")

	branching, err := strconv.ParseInt(r.FormValue("branching"), 10, 32)
	if err != nil {
                resp.Error = "Cannot read road exits"
                return false, resp, initials, city_map
        }
	resp.Branching = int(branching)

	fmt.Println("Ok.")

	initials.Raduis.Min = min_r
	initials.Raduis.Max = max_r
	initials.NumCenters = int(n_centers)
	initials.Branching = int(branching)

	resp.Default = false
	return true, resp, initials, city_map
}

func generateRoads(initial_values gen.InitialValuesRoads, cm city_map.Map) city_map.Map {
	channel := make(chan city_map.Map) // Bad API. Needs refactoring

	cm.Roads = gen.GenerateRoads(cm, channel, initial_values)

	return cm
}


func handleGetRoads(w http.ResponseWriter, r *http.Request) (RoadsResponse) {
	response := RoadsResponse{Default: true, Error: ""}
	var city_map city_map.Map

	city_map = getMapFromCookies(r)

	img, success := makeImageString(city_map)
	if success {
		response.Image = img
	}

	return response
}


func handlePostRoads(w http.ResponseWriter, r *http.Request) RoadsResponse {
	success, response, initial_values, city_map := readRoadsParams(r)

	if !success {
		img, success := makeImageString(city_map)
		if success {
		    response.Image = img
		}
		return response
	}

	fmt.Println()
	fmt.Println("Roads read successfully? -", success)
	fmt.Println(initial_values)

	city_map = generateRoads(initial_values, city_map)

	fmt.Println("Ok")

	setMapCookies(city_map, w)

	img, success := makeImageString(city_map)
	if success {
		response.Image = img
	}

	return response
}


func roadsHandler(w http.ResponseWriter, r *http.Request){
	fmt.Printf("Serving %s for %s", r.Host, r.URL.Path)
	index_template := template.Must(template.ParseGlob("./html/roads.gohtml"))

	var response RoadsResponse
	switch r.Method {
		case http.MethodGet:
			response = handleGetRoads(w, r)
		case http.MethodPost:
			response = handlePostRoads(w,r)
	}

	index_template.ExecuteTemplate(w, "roads.gohtml", response)
}
