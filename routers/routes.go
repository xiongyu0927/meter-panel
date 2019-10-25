package routers

import (
	"meter-panel/controllers"
	_ "meter-panel/exporter"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Mux is router
var Mux *http.ServeMux = http.NewServeMux()

func init() {
	// parmeter requierd cluster name string
	Mux.HandleFunc("/api/v1/k8s/nodes", controllers.NodeViews)
	Mux.HandleFunc("/api/v1/k8s/pods", controllers.PodViews)
	Mux.HandleFunc("/api/v1/k8s/app", controllers.AppViews)
	Mux.HandleFunc("/api/v1/k8s/pvs", controllers.PvViews)
	Mux.HandleFunc("/api/v1/k8s/lbs", controllers.LbViews)
	Mux.HandleFunc("/api/v1/k8s/cpu", controllers.CpuViews)
	Mux.HandleFunc("/api/v1/k8s/mem", controllers.MemViews)
	Mux.HandleFunc("/api/v1/k8s/alerts", controllers.AlertViews)
	Mux.HandleFunc("/api/v1/k8s/events", controllers.EventsViews)
	Mux.HandleFunc("/api/v1/k8s/pipelines", controllers.PipelineViews)
	Mux.HandleFunc("/api/v1/k8s/codequality", controllers.CqbViews)
	Mux.HandleFunc("/api/v1/k8s/projects", controllers.ProjectViews)
	Mux.Handle("/metrics", promhttp.Handler())
}
