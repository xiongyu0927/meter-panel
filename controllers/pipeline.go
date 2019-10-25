package controllers

import (
	"encoding/json"
	"log"
	"meter-panel/store"
	"net/http"
)

// PipelineViews is api of return node number and status
func PipelineViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	pipel := store.AllStore.PipelineStore.List()
	prj, err := getProjectFromCluster(cluster)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	tmp := Style.OrganizePipelineList(prj, pipel)
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(data2)
}
