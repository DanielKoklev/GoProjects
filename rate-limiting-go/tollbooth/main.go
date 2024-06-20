package main

import (
	"encoding/json"
	"log"
	"net/http"

	tollbooth "github.com/didip/tollbooth/v7"
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

func main(){
	message := Message {
		Status: "Request failed!",
		Body: "The API is at full capacity!",
	}

	jsonMessage, _ := json.Marshal(message)
	tlbhLimiter := tollbooth.NewLimiter(1, nil)
	tlbhLimiter.SetMessageContentType("application/json")
	tlbhLimiter.SetMessage(string(jsonMessage))
	http.Handle("/ping", tollbooth.LimitFuncHandler(tlbhLimiter, endpointHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error while listening on port 8080!!", err)
	}

}