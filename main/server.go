package main

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "hello world")
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	// Use the thumbnailHandler function
	http.HandleFunc("/api/thumbnail", thumbnailHandler)

	// Serve static files from the frontend/dist directory
	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", fs)

	// Start the server.
	fmt.Println("Server listening on port 3000")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}

type thumbnailRequest struct {
	Url string `json:"url"`
}

type screenshotAPIRequest struct {
	token string `json:"token"`
	url string `json:"url"`
	width int `json:"width"`
	height int `json:"height"`
	output string `json:"output"`
	file_type string `json:"file_type"`
	wait_for_event string `json:"wait_for_event"`
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	var decoded thumbnailRequest

	// Try to decode the request into the thumbnailRequest struct
	err := json.NewDecoder(r.Body).Decode(&decoded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a HTTP request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://shot.screenshotapi.net/screenshot?token=GSFSNEE-45VMGZE-G660VB5-EJCPNQD&url=%s&width=1920&height=1080&output=json&file_type=png&wait_for_event=load&thumbnail_width=600", decoded.Url), nil)

	// Execute HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	checkError(err)

	// Close response at end of function
	defer response.Body.Close()

	// Read raw response into Go struct
	type screenshotAPIResponse struct {
		Screenshot string `json"screenshot"`
	}
	var apiResponse screenshotAPIResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	checkError(err)

	// Pass back the screenshot URL to the frontend
	_, err = fmt.Fprintf(w, `{ "screenshot": "%s" }`, apiResponse.Screenshot)
	checkError(err)
}