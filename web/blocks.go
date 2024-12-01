package main

import (
	"fmt"
	"net/http"
	"html/template"
	"strconv"
	gen "chirrwick.com/projects/city/generator"
	"chirrwick.com/projects/city/city_map"
	"encoding/json"
)

type BlocksResponse struct{
	Default bool
	Error string
	Map string
	Image string
	MIN_SIZE float64
	MAX_SIZE float64
}

func readBlocksParams(r *http.Request) (bool, BlocksResponse, gen.InitialValuesBlocks, city_map.Map) {
	resp := BlocksResponse{Error: "", Map: "{}", Default: false}

	fmt.Println("Prepare initials for areas")
	var initials gen.InitialValuesBlocks

	fmt.Println("ReadMap")
/*
	map_string := r.FormValue("map")
	resp.Map = map_string
	map_json := []byte(map_string)
	var city_map city_map.Map
	err := json.Unmarshal(map_json, &city_map)
	if err != nil {
		resp.Error = "Cannot get map from you"
		return false, resp, initials, city_map
	}
	resp.Default = true
*/

	city_map := getMapFromCookies(r)

	fmt.Println("Ok. Read min")

        min_size, err := strconv.ParseFloat(r.FormValue("min"), 64)
        if err != nil {
		resp.Error = "Cannot read min block size"
                return false, resp, initials, city_map
	}
	resp.MIN_SIZE = min_size

	fmt.Println("Ok. Read max")

	max_size, err := strconv.ParseFloat(r.FormValue("max"), 64)
	if err != nil {
                resp.Error = "Cannot read max size"
                return false, resp, initials, city_map
        }
	resp.MAX_SIZE = max_size

	fmt.Println("Ok.")

	initials.Size.Min = min_size
	initials.Size.Max = max_size

	resp.Default = false
	return true, resp, initials, city_map
}

func generateBlocks(initial_values gen.InitialValuesBlocks, cm city_map.Map) city_map.Map {
	channel := make(chan city_map.Map) // Bad API. Needs refactoring

	cm.Blocks = gen.GenerateBlocks(cm, channel, initial_values)

	return cm
}



func handleGetBlocks(w http.ResponseWriter, r *http.Request) BlocksResponse {
	response := BlocksResponse{Default: true, Error: "", Map: "{}"}

	city_map := getMapFromCookies(r)

	img, success := makeImageString(city_map)
	if success {
	    response.Image = img
	}

	return response
}


func handlePostBlocks(w http.ResponseWriter, r *http.Request) BlocksResponse {
	
	success, response, initial_values, city_map := readBlocksParams(r)

	if !success {
		img, success := makeImageString(city_map)
		if success {
		    response.Image = img
		}
		return response
	}

	fmt.Println()
	fmt.Println("Blocks read successfully? -", success)
	fmt.Println(initial_values)

	city_map = generateBlocks(initial_values, city_map)

	fmt.Println("Ok")

	map_json, err := json.Marshal(city_map)
	if err != nil {
                response.Error = "Error while generating or serializing areas"
		response.Map = "{}"
        }
	fmt.Println(map_json)
	response.Map = string(map_json)

	setMapCookies(city_map, w)

	img, success := makeImageString(city_map)
	if success {
		response.Image = img
	}

	return response
}


func blocksHandler(w http.ResponseWriter, r *http.Request){
	fmt.Printf("Serving %s for %s", r.Host, r.URL.Path)
	index_template := template.Must(template.ParseGlob("./html/blocks.gohtml"))

	var response BlocksResponse
	switch r.Method {
		case http.MethodGet:
			response = handleGetBlocks(w,r)
		case http.MethodPost:
			response = handlePostBlocks(w,r)
	}

	index_template.ExecuteTemplate(w, "blocks.gohtml", response)
}
