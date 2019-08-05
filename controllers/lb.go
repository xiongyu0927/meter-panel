package controllers

import (
	"encoding/json"
	"log"
	"meter-panel/store"
	"net/http"
)

func LbViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	tmp := store.GetSingleClusterLBList(cluster)
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
	}
	w.Write(data2)
}
