package controllers

import (
	"encoding/json"
	"log"
	"meter-panel/pkg/api/v1/k8s/constome"
	"meter-panel/store"
	"net/http"

	"k8s.io/apimachinery/pkg/labels"
)

func LbViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")

	lister := store.AllLister.SvcLister[cluster]
	if lister == nil {
		store.AddNewClusterResource(cluster)
		lister = store.AllLister.SvcLister[cluster]
	}
	sl, err := lister.List(labels.Everything())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	tmp1 := Style.OrganizeSvcList(cluster, sl)
	t1, ok := tmp1.([]constome.CebClusterSvcList)
	if !ok {
		log.Println("Internal error: Wrong Type")
		http.Error(w, "Internal error: Wrong Type", 400)
		return
	}

	lister2 := store.AllLister.IngressLister[cluster]
	if lister2 == nil {
		store.AddNewClusterResource(cluster)
		lister2 = store.AllLister.IngressLister[cluster]
	}
	il, err := lister2.List(labels.Everything())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	tmp2 := Style.OrganizeIngressList(cluster, il)
	t2, ok := tmp2.([]constome.CebClusterIngressList)
	if !ok {
		log.Println("Internal error: Wrong Type")
		http.Error(w, "Internal error: Wrong Type", 400)
		return
	}

	var tmp3 = constome.CebLB{
		L4: t1,
		L7: t2,
	}

	data2, err := json.Marshal(tmp3)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(data2)
}
