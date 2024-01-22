package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type GeoJson struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Crs  struct {
		Type       string `json:"type"`
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
	} `json:"crs"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			Name string `json:"name"`
			Iso  string `json:"iso"`
		} `json:"properties"`
		Geometry struct {
			Type        string          `json:"type"`
			Coordinates [][][][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

type GeoJsonData struct {
	DB []GeoJson
}

type GeoJsonDB struct {
	DB []GeoJsonData
}

func (d GeoJsonDB) read() {
	files, err := os.ReadDir("files/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("reading ...")
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}

}

func forward2Request(LD string) {
	posturl := "http://127.0.0.1:7900"
	body := []byte(LD)

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("forward LD", LD, "to", posturl, res)
}

func inputConsul(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)

	fmt.Println("got data: ", bodyString)
	forward2Request(bodyString)
}

func main() {
	fmt.Println("--- Geozone Lookup Manager ---")

	var myDB GeoJsonDB
	myDB.read()

	http.HandleFunc("/", inputConsul)
	log.Fatal(http.ListenAndServe(":7000", nil))
}
