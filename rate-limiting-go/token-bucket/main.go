package main

import (
	"log"
	"net/http"
	"encoding/json"
)

type Message struct {
	Status string `json:"status"`
	Body string `json:"body"`
}

func endpointHandler(writer http.ResponseWriter, request *http.Request){
	writer.Header().Set("Content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	message := Message {
		Status: "Successfull",
		Body: "You reached the API!!!",
	}
	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		return 
	}
}

func main() {
	http.Handle("/ping", rateLimiter(endpointHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error while listening on port 8080 !", err)
	}
}