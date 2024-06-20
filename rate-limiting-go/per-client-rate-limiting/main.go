package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Message struct {
	Status string `json:"status"`
	Body string `json:"body"`
}


func perClientLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
	
	type Client struct {
		limiter *rate.Limiter
		lastSeen time.Time
	}

	var (
		mtx sync.Mutex
		Clients = make(map[string]*Client)
	)

	go func(){
		for {
			time.Sleep(time.Minute)
			mtx.Lock()
			for ip, Client := range Clients {
				if time.Since(Client.lastSeen) > 3*time.Minute{
					delete(Clients, ip)
				}
			}
			mtx.Unlock()
		}
	}()
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		mtx.Lock()
		if _, found := Clients[ip]; !found{
			Clients[ip] = &Client{limiter: rate.NewLimiter(2, 4)}
		}
		Clients[ip].lastSeen = time.Now()
		if !Clients[ip].limiter.Allow() {
			mtx.Unlock()
			message := Message {
				Status: "Request failed!",
				Body: "The API is at full capacity!",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&message)
			return
		}
		mtx.Unlock()
		next(w, r)
	})
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
	http.Handle("/ping", perClientLimiter(endpointHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error while trying to listen on port 8080!", err)
	}

}