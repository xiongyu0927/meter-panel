package controllers

import (
	"encoding/json"
	"log"
	capacity "meter-panel/capacity_report"
	"net/http"
)

// AppViews is api of return node number and status
func ClusterScalCapacityViews(w http.ResponseWriter, r *http.Request) {
	tmp := capacity.GetClusterScalData()
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(data2)
}
