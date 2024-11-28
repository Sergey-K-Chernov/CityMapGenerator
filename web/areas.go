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

type AreasResponse struct{
	Default bool
	Error string
	Map string
	Image string
	NIndustrial int
	PercentageIndustrial float64
	NParks int
	PercentageParks float64
}

func readAreasParams(r *http.Request) (bool, AreasResponse, gen.InitialValuesAreas, city_map.Map) {
	resp := AreasResponse{Error: "", Map: "{}", Default: false}

	fmt.Println("Prepare initials for areas")
	var initials gen.InitialValuesAreas

	fmt.Println("ReadMap")

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

	fmt.Println("Ok. Read indus")

	n_indstrl, err := strconv.ParseInt(r.FormValue("n_industrial"), 10, 32)
	fmt.Println(n_indstrl)
	if err != nil {
		fmt.Println("Error")
		resp.Error = "Cannot read industrial areas number"
		return false, resp, initials, city_map
	}
	resp.NIndustrial = int(n_indstrl)

	fmt.Println("Ok. Read % indus")

        prcntge_indstrl, err := strconv.ParseFloat(r.FormValue("percentage_industrial"), 64)
        if err != nil {
		resp.Error = "Cannot read industrial areas percentage"
                return false, resp, initials, city_map
	}
	resp.PercentageIndustrial = prcntge_indstrl

	fmt.Println("Ok. Read parks")

	n_parks, err := strconv.ParseInt(r.FormValue("n_parks"), 10, 32)
        if err != nil {
                resp.Error = "Cannot read parks"
                return false, resp, initials, city_map
        }
	resp.NParks = int(n_parks)

	fmt.Println("Ok. Read % parks")

	prcntge_parks, err := strconv.ParseFloat(r.FormValue("percentage_parks"), 64)
	if err != nil {
                resp.Error = "Cannot read % parks"
                return false, resp, initials, city_map
        }
	resp.PercentageParks = prcntge_parks

	fmt.Println("Ok.")

	initials.NumIndustrial = int(n_indstrl)
	initials.AreaIndustrial = prcntge_indstrl
	initials.NumParks = int(n_parks)
	initials.AreaParks = prcntge_parks

	resp.Default = false
	return true, resp, initials, city_map
}

func generateAreas(initial_values gen.InitialValuesAreas, cm city_map.Map) city_map.Map {
	channel := make(chan city_map.Map) // Bad API. Needs refactoring

	cm.Areas = gen.GenerateAreas(cm, channel, initial_values)

	return cm
}


func areasHandler(w http.ResponseWriter, r *http.Request){
	fmt.Printf("Serving %s for %s", r.Host, r.URL.Path)
	index_template := template.Must(template.ParseGlob("./html/areas.gohtml"))

	success, response, initial_values, city_map := readAreasParams(r)

	if !success {
		img, success := makeImageString(city_map)
		if success {
		    response.Image = img
		}
		index_template.ExecuteTemplate(w, "areas.gohtml", response)
		
		return
	}

	fmt.Println()
	fmt.Println("Areas read successfully? -", success)
	fmt.Println(initial_values)

	city_map = generateAreas(initial_values, city_map)

	fmt.Println("Ok")

	map_json, err := json.Marshal(city_map)
	if err != nil {
                response.Error = "Error while generating or serializing areas"
		response.Map = "{}"
        }
	fmt.Println(map_json)
	response.Map = string(map_json)

	img, success := makeImageString(city_map)
	if success {
		response.Image = img
	}

	index_template.ExecuteTemplate(w, "areas.gohtml", response)
}
