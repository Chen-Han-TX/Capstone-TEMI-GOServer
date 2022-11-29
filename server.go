package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Clients struct {
	XMLName xml.Name `xml:"clients"`
	Clients []Client `xml:"client"`
}

type Client struct {
	XMLName xml.Name `xml:"client"`
	Name    string   `xml:"name"`
	Ip      string   `xml:"ip"`
	Port    string   `xml:"port"`
}

type Result struct {
	Level   string `json:"level"`
	ShelfNo string `json:"shelfno"`
}

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
	// check

	fmt.Printf("Starting server at port 10000\n")
	log.Fatal(http.ListenAndServe(":10000", myRouter))

}

var count int = 0
var m sync.Mutex

func add(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	for i := 0; i < 10; i++ {
		count += 1
	}
	fmt.Fprintf(w, "Level: "+strconv.Itoa(count))

	m.Unlock()

}
func main() {
	handleRequests()
}

func wronglevel(w http.ResponseWriter, r *http.Request) {
	m.Lock()
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
	xmlFile, err := os.Open("addresses.xml")
	if err != nil {
		fmt.Println(err)
	}

	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var clients Clients
	xml.Unmarshal(byteValue, &clients)
	fmt.Fprintf(w, clients.Clients[0].Name)
	fmt.Fprintf(w, clients.Clients[0].Ip)
	fmt.Fprintf(w, clients.Clients[0].Port)

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
	m.Unlock()

}
