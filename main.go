package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
	http.HandleFunc("/", inputConsul)
	log.Fatal(http.ListenAndServe(":7000", nil))
}
