package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Requestdata []string

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

const (
	active   int = 0
	inactive     = 1
	blocked      = 2
)

type GeoJsonData struct {
	name       string
	from       time.Time
	to         time.Time
	curr_state int
	DB         GeoJson
}

func (d GeoJsonData) init(in GeoJson, name string) GeoJsonData {
	fmt.Println("init", name)
	d.name = name
	d.from = time.Now()
	d.to = d.from.Add(25 * time.Hour)
	d.curr_state = active
	d.DB = in
	return d
}

type GeoJsonDB struct {
	DB []GeoJsonData
}

func (d GeoJsonDB) read() []GeoJsonData {
	file, err := os.Open("files/")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	list, err := file.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range list {
		fmt.Println(f.Name())

		//b1 := make([]byte, 100)
		contents, err := os.ReadFile("files/" + f.Name())
		if err != nil {
			fmt.Println("File reading error", err)
		}
		//fmt.Println("Contents of file:", string(contents))

		//jsonFile := "files/" + file.Name()
		//byteValue, _ := ioutil.Read(file)
		var data GeoJson
		json.Unmarshal(contents, &data)

		var df GeoJsonData

		d.DB = append(d.DB, df.init(data, f.Name()))
		//fmt.Println("db->", d.DB)

	}
	return d.DB
}

func (d GeoJsonDB) init() {
	files, err := os.ReadDir("files/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("reading ...")
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
}

type validationResult struct {
	Ld       string `json:"LD"`
	ID       int    `json:"Id"`
	Ipadress string `json:"Ipadress"`
}

func (d GeoJsonDB) validate(a string) bool {
	files, err := os.ReadDir("files/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("reading ...")
	for _, file := range files {
		name := file.Name()
		name = name[:len(name)-8]
		if name == a {
			return true
		}
		fmt.Println(name, a)
	}
	return false
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
func isValid(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "OK")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	//bodyString := string(bodyBytes)
	var Requestdata []string

	err = json.Unmarshal(bodyBytes, &Requestdata)

	if err != nil {
		fmt.Println(err)
	}

	var myDB GeoJsonDB

	for i := range Requestdata {
		fmt.Println(Requestdata[i])
		if myDB.validate(Requestdata[i]) == true {
			var ret validationResult
			ret.Ld = Requestdata[i]
			ret.ID = 2 //later read from config
			ret.Ipadress = "3.16.183.50"

			j, _ := json.Marshal(ret)
			fmt.Fprintf(w, string(j))
		}
	}
	//fmt.Println(myDB.DB)
	w.WriteHeader(http.StatusNotFound)

}

func inputConsul(w http.ResponseWriter, r *http.Request) {

	var myDB GeoJsonDB
	myDB.read()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	bodyString := string(bodyBytes)

	fmt.Println("got data: ", bodyString)
	bodyString = strings.ReplaceAll(bodyString, ".service.consul", "")
	if myDB.validate(bodyString) == true {
		var ret validationResult
		ret.Ld = bodyString
		ret.ID = 2 //later read from config
		ret.Ipadress = "3.16.183.50"

		j, _ := json.Marshal(ret)
		forward2Request(string(j))
	}
	fmt.Fprintf(w, "NOK")
}

func whoami(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GeozoneLookup Manager")
}

func main() {
	fmt.Println("--- Geozone Lookup Manager (:7000)---")

	http.HandleFunc("/", inputConsul)
	http.HandleFunc("/isvalid", isValid)
	http.HandleFunc("/whoami", whoami)
	log.Fatal(http.ListenAndServe(":7000", nil))
}
