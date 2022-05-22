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

func startRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/svc", GetService).Methods(http.MethodGet)
	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

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

func isAllowed() bool {
	var MaxRequests uint64
	var WindowDuration string

	if req, exist := os.LookupEnv("MAX_REQUESTS"); exist != true {
		MaxRequests = 3
	} else {
		n, err := strconv.ParseInt(req, 10, 64)
		if err != nil {
			log.Fatalf("%d of type %T", n, n)
		}
		MaxRequests = uint64(n)
	}

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
