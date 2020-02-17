package controllers

import (
	"encoding/json"
	"meter-panel/pkg/api/v1/k8s/constome"
	"net/http"
)

var (
	Style constome.OrganizeData

	healthInfo, _ = json.Marshal("I'm fine thank you And you?")
)

const (
	EPrometheus string = "Can't find the correct prometheus address"
)

func init() {
	Style = constome.NewCebStyle()
}

func HealthViews(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write(healthInfo)
}
