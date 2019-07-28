package routers

import (
	"meter-panel/controllers"
	"net/http"
)

// Mux is router
var Mux *http.ServeMux = http.NewServeMux()

func init() {
	// parmeter requierd cluster name string
	Mux.HandleFunc("/api/v1/k8s/node", controllers.NodeViews)
}