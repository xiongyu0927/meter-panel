package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

// PipelineViews is api of return node number and status
func CqbViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	// cqbl := store.AllStore.CqbStore.List()
	prj, err := getProjectFromCluster(cluster)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	tmp := Style.OrganizeCqbList(prj)
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(data2)
}
