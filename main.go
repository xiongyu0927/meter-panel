package main

import (
	"log"
	"net/http"
	"nsxt-go/routers"
)

func main() {
	err := http.ListenAndServe(":12300", routers.Mux)
	if err != nil {
		log.Fatal(err)
	}
}
