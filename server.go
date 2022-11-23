package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/wronglevel", wronglevel)
	myRouter.HandleFunc("/add", add)
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument

	fmt.Printf("Starting server at port 10000\n")
	log.Fatal(http.ListenAndServe("192.168.0.112:10000", myRouter))

}

var count int = 0

func add(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		count += 1
	}
	fmt.Fprintf(w, "Level: "+strconv.Itoa(count))

}
func main() {
	handleRequests()
}

type Result struct {
	Level   string `json:"level"`
	ShelfNo string `json:"shelfno"`
}

func wronglevel(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var result Result
	err := json.Unmarshal(reqBody, &result)

	if err != nil {
		fmt.Printf("Could not unmarshal Json %s\n", err)
	}
	fmt.Printf("Level: " + result.Level + ", ShelfNo: " + result.ShelfNo)
	fmt.Fprintf(w, "Level: "+result.Level+", ShelfNo: "+result.ShelfNo)

	// Match the pattern using XML
	// using placeholder
	url := "http://192.168.0.205:8080/?level=" + result.Level + "&shelfno=" + result.ShelfNo
	fmt.Println("URL:>", url)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
