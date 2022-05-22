package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	startRouter()
}

func startRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/{id}", GetService).Methods("GET")
	http.Handle("/", r)
}

func GetService(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		return
	}

	w.Write([]byte("Allowed"))
}
