package controllers

import (
	"encoding/json"
	"log"
	"meter-panel/store"
	"net/http"
)

// PodViews is api of return node number and status
func PodViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	tmp := store.GetSingleClusterPodsList(cluster)
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
	}
	w.Write(data2)
}
