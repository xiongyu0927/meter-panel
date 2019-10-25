package controllers

import (
	"encoding/json"
	"log"
	"meter-panel/store"
	"net/http"
)

// PipelineViews is api of return node number and status
func ProjectViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	prjl := store.AllStore.ProjectStore.List()
	tmp := Style.OrganizeProjectList(cluster, prjl)
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(data2)
}
