package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	startRouter()
}

/*
	A simple HTTP listener, listening on localhost:8080/svc GET
	Parameters:
	Returns:
*/
func startRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/svc", GetService).Methods(http.MethodGet)
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

/*
	A simple functional service to reply a request to the endpoint localhost:8080/svc GET
	It will then call the function that verifies if the requester is allowed to call this
*/
func GetService(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if isAllowed() {
			w.Write([]byte("Allowed"))
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Not Allowed"))
		}
		return
	}
}

/*
	Makes a decision whether UserID is allowed to perform the current request.
	Parameters:
	Returns bool,
		returns true if requester is within the limits set(requests and window time)
*/
func isAllowed() bool {
	var MaxRequests uint64
	var WindowDuration string

	// Using Environment variables to allow customisation of MAX_REQUESTS
	if req, exist := os.LookupEnv("MAX_REQUESTS"); exist != true {
		MaxRequests = 3
	} else {
		n, err := strconv.ParseInt(req, 10, 64)
		if err != nil {
			log.Fatalf("%d of type %T", n, n)
		}
		MaxRequests = uint64(n)
	}

	// Using Environment variables to allow customisation of WINDOW_DURATION
	if winDuration, exist := os.LookupEnv("WINDOW_DURATION"); exist != true {
		WindowDuration = "1m"
	} else {
		WindowDuration = winDuration
	}

	Duration, _ := time.ParseDuration(WindowDuration)

	sw := GetSWInstance()
	sw.MaxRequests = MaxRequests
	sw.WindowDuration = Duration

	var limiterInterface RateLimiterInterface = sw

	// Here the userid is being mocked, but could be retrieved by any means
	// such as IP, userID from http body/header, etc
	res, err := limiterInterface.IsRequestAllowed("user123")

	if err != nil {
		log.Fatalf("Error trying to allow user: %v", err)
	}

	return res
}
