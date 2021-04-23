package main

import (
	"github.com/gorilla/mux"
	"github.com/kbase/kbase-sdk-module-jsonschema/paths"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}
	address := ":" + port
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/rpc", paths.HandleRPC).Methods("POST")
	router.HandleFunc("/about", paths.HandleGetAbout).Methods("GET")
	router.HandleFunc("/schema/{path:.*}/{schema}.{version}.json", paths.HandleGetSchema).Methods("GET")
	router.HandleFunc("/schema/{path:.*}/{schema}.json", paths.HandleGetSchema).Methods("GET")
	log.Fatal(http.ListenAndServe(address, router))
}
