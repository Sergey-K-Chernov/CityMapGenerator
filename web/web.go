package main

import (
	"fmt"
	"net/http"
	"html/template"
	"strconv"
	gen "chirrwick.com/projects/city/generator"
	"chirrwick.com/projects/city/city_map"
)

type BordersResponse struct{
	Error string
	Map string
	MinR float64
	MaxR float64
	NCorners int
	Variation float64
}

func mainHandler(w http.ResponseWriter, r *http.Request){
	fmt.Printf("Serving: %s for %s\n", r.Host, r.URL.Path)
	index_template := template.Must(template.ParseGlob("./html/index.gohtml"))
	index_template.ExecuteTemplate(w, "index.gohtml", nil)
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

	fmt.Printf("")

	channel := make(chan city_map.Map)

	go gen.GenerateBorders(channel, initials)

	city_map := <- channel

	fmt.Println(city_map)

	return true, resp, initials
}

func bordersHandler(w http.ResponseWriter, r *http.Request){
	fmt.Printf("Serving %s for %s", r.Host, r.URL.Path)
	index_template := template.Must(template.ParseGlob("./html/borders.gohtml"))

	success, response, initial_values := readBorderParams(r)
	fmt.Println()
	fmt.Println("Borders read successfully? -", success)
	fmt.Println(initial_values)
	index_template.ExecuteTemplate(w, "borders.gohtml", response)
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/borders", bordersHandler)
	err := http.ListenAndServe("193.168.173.245:80", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}
