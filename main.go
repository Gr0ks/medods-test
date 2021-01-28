package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
}

func main() {
	// init gorilla/mux router
	r := mux.NewRouter()

	r.HandleFunc("api/auth", getNewToken).Methods("POST")
	r.HandleFunc("api/auth/refreshToken", refreshToken).Methods("POST")
	log.Fatal(http.ListenAndServe(":8089", r))
}





