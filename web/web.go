package main

import (
	"fmt"
	"net/http"
	"html/template"
	"chirrwick.com/projects/city/city_map"
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"encoding/base64"
	md "chirrwick.com/projects/city/draw"
)

func mainHandler(w http.ResponseWriter, r *http.Request){
	fmt.Printf("Serving: %s for %s\n", r.Host, r.URL.Path)
	index_template := template.Must(template.ParseGlob("./html/index.gohtml"))
	index_template.ExecuteTemplate(w, "index.gohtml", nil)
}

func makeImageString(city_map city_map.Map) (string, bool) {
	img := image.NewRGBA(image.Rect(0, 0, 512, 512))
	green := color.RGBA{0, 255, 0, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{green}, image.ZP, draw.Src)

	map_img := md.Draw(city_map)
	img = map_img

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return "", false
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return str, true
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/borders", bordersHandler)
	http.HandleFunc("/roads", roadsHandler)
	http.HandleFunc("/areas", areasHandler)
	http.HandleFunc("/blocks", blocksHandler)
	err := http.ListenAndServe("193.168.173.245:80", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

}
