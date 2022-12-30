package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

type Image struct {
	Content string `json:"image"`
}

var count int = 0
var m sync.Mutex

// Given a base64 string of a JPEG, encodes it into an JPEG image test.jpg
func base64toJpg(data string) {

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	fmt.Println("base64toJpg", bounds, formatString)

	//Encode from image format to writer
	pngFilename := "test.jpg"
	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = jpeg.Encode(f, m, &jpeg.Options{Quality: 75})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Jpg file", pngFilename, "created")

}

// ----- HANDLER FUNCTIONS ------

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func add(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	for i := 0; i < 10; i++ {
		count += 1
	}
	fmt.Fprintf(w, "Level: "+strconv.Itoa(count))
	m.Unlock()
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

	var url string
	if result.Level == "4" {
		url = "http://" + clients.Clients[1].Ip + ":" + clients.Clients[1].Port
	} else if result.Level == "3" {
		url = "http://" + clients.Clients[0].Ip + ":" + clients.Clients[0].Port
	}

	fmt.Println("\nURL:>", url)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

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

func getImage(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var image Image
	err := json.Unmarshal(reqBody, &image)

	if err != nil {
		fmt.Printf("Could not unmarshal Json %s\n", err)
	}

	//image_url := "data:image/jpeg;base64," + image.Content

	fmt.Printf("Image bitmap: " + image.Content)

	// decode the image bitmap
	base64toJpg(image.Content)

	// save the image in the src folder

}

func main() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/wronglevel", wronglevel)
	myRouter.HandleFunc("/add", add)
	myRouter.HandleFunc("/image", getImage)

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	// check

	fmt.Printf("Starting server at port 10000\n")
	log.Fatal(http.ListenAndServe(":10000", myRouter))

}
