package controllers

import (
	"encoding/json"
	"log"
	"meter-panel/store"
	"net/http"

	"k8s.io/apimachinery/pkg/labels"
)

// PodViews is api of return node number and status
func PodViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	lister := store.AllLister.PodLister[cluster]
	if lister == nil {
		store.AddNewClusterResource(cluster)
		lister = store.AllLister.PodLister[cluster]
	}
	pl, err := lister.List(labels.Everything())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	tmp := Style.OrganizePodsList(cluster, pl)
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(data2)
}
