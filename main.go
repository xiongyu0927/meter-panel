package main

import (
	"log"
	_ "meter-panel/exporter"
	"meter-panel/routers"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":12300", routers.Mux)
	if err != nil {
		log.Fatal(err)
	}
}
