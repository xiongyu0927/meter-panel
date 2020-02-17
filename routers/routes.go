package routers

import (
	"meter-panel/controllers"
	_ "meter-panel/exporter"
	"net/http"
	"net/http/pprof"

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
	// doesn't need parmeter
	Mux.HandleFunc("/api/v1/ceb/capacity", controllers.CapacityViews)
	Mux.HandleFunc("/health", controllers.HealthViews)
	// Debug api
	Mux.Handle("/metrics", promhttp.Handler())
	Mux.HandleFunc("/debug/pprof/", pprof.Index)
	Mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	Mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	Mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	Mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
