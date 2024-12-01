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
	"strings"
	"strconv"
	"encoding/json"
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


const COOKIE_MAX_SIZE = 4000

func jsonToCookieStrings(json_value []byte) (strs []string) {
	str := string(json_value)
	str = strings.ReplaceAll(str, "\"", "%22")
	
	size := len(str)
	n := size / COOKIE_MAX_SIZE + 1
	for i := 0; i < n; i++ {
		i1 := i*COOKIE_MAX_SIZE
		i2 := ((i+1)*COOKIE_MAX_SIZE)
		i2 = min(i2, len(str))
		strs = append(strs, str[i1:i2])
	}
	return strs
}


func cookieStringsToJson(strs []string) []byte {
	str := ""
	for _, s := range strs {
		str += s
	}
	
	str = strings.ReplaceAll(str, "%22", "\"")
	return []byte(str)
}


func setMapCookies(m city_map.Map, w http.ResponseWriter) {
	map_json, err := json.Marshal(m)

	if err != nil {
		fmt.Println("Map to json error: ")
		fmt.Println(err)
		return
	}

	cs := jsonToCookieStrings(map_json)
	cookie := &http.Cookie{
		Name: "MapCookiesNum",
		Value: strconv.Itoa(len(cs)),
		MaxAge: 3600,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	
	for i, s := range cs {
		fmt.Printf("Set Cookie %d\n", i)
		cookie := &http.Cookie{
			Name: "Map" + strconv.Itoa(i),
			Value: s,
			MaxAge: 3600,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)
	}
}

func getMapFromCookies(r *http.Request) (m city_map.Map) {
	cookie, err := r.Cookie("MapCookiesNum")
	if err != nil {
		fmt.Println("No cookies found: ")
		fmt.Println(err)
		return
	}

	n, err := strconv.Atoi(cookie.Value)
	if err != nil {
		fmt.Println("Error reading cookies")
		fmt.Println(err)
		return
	}

	strs := make([]string, 0)
	for i := range n {
		cookie, err :=  r.Cookie("Map" + strconv.Itoa(i))
		if err != nil {
			fmt.Println("Error reading cookies")
			fmt.Println(err)
			return
		}
		strs = append(strs, cookie.Value)
	}

	map_json := cookieStringsToJson(strs)

	err = json.Unmarshal(map_json, &m)
	if err != nil {
		fmt.Println("Error unmarshalling map:")
		fmt.Println(err)
		return
	}

	return
}