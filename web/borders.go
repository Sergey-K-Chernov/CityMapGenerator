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

type BordersResponse struct{
	Error string
	Map string
	Image string
	MinR float64
	MaxR float64
	NCorners int
	Variation float64
}

func readBorderParams(r *http.Request) (bool, BordersResponse, gen.InitialValuesMap) {
	resp := BordersResponse{Error: "", Map: "{}"}
	var initials gen.InitialValuesMap

	min_r, err := strconv.ParseFloat(r.FormValue("min_r"), 32)
        if err != nil {
		resp.Error = "Cannot read min radius"
		return false, resp, initials
	}
	resp.MinR = min_r

        max_r, err := strconv.ParseFloat(r.FormValue("max_r"), 32)
        if err != nil {
		resp.Error = "Cannot read max radius"
                return false, resp, initials
	}
	resp.MaxR = max_r

	n_corners, err := strconv.ParseInt(r.FormValue("n_corners"), 10, 32)
        if err != nil {
                resp.Error = "Cannot read corners"
                return false, resp, initials
        }
	resp.NCorners = int(n_corners)

	vartn, err := strconv.ParseFloat(r.FormValue("variation"), 32)
	if err != nil {
                resp.Error = "Cannot read variation"
                return false, resp, initials
        }
	resp.Variation = vartn

	initials.Raduis.Min = min_r
	initials.Raduis.Max = max_r
	initials.NumSides = int(n_corners)
	initials.VertexShift = vartn

	return true, resp, initials
}

func generateBorders(initial_values gen.InitialValuesMap) city_map.Map {
        channel := make(chan city_map.Map)

        go gen.GenerateBorders(channel, initial_values)

        city_map := <- channel

        return city_map
}


func bordersHandler(w http.ResponseWriter, r *http.Request){
	fmt.Printf("Serving %s for %s", r.Host, r.URL.Path)
	index_template := template.Must(template.ParseGlob("./html/borders.gohtml"))

	success, response, initial_values := readBorderParams(r)

	if !success {
		index_template.ExecuteTemplate(w, "borders.gohtml", response)
		return
	}

	//fmt.Println()
	//fmt.Println("Borders read successfully? -", success)
	//fmt.Println(initial_values)

	city_map := generateBorders(initial_values)

	map_json, err := json.Marshal(city_map)
	if err != nil {
                response.Error = "Error while generating map"
		response.Map = "{}"
        }
	//fmt.Println(map_json)
	response.Map = string(map_json)

	img, success := makeImageString(city_map)
	if success {
		response.Image = img
	}

	index_template.ExecuteTemplate(w, "borders.gohtml", response)
}