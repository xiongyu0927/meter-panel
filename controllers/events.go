package controllers

import (
	"encoding/json"
	"log"
	"meter-panel/store"
	"net/http"
)

type CebEvent struct {
	Total   int64
	Warning int64
}

func EventsViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	tmp, err := GetEventList(cluster)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(data2)
}

func GetEventList(cluster string) (CebEvent, error) {
	a := store.EsClient
	var total int64
	var warning int64
	var tmp CebEvent
	err := store.EsClient.SearchCount(cluster, a.Index[0])
	if err != nil {
		return tmp, err
	}
	for _, v := range a.Data[cluster] {
		warning = warning + v[1]
		total = warning + v[0]
	}
	tmp = CebEvent{
		Total:   total,
		Warning: warning,
	}
	return tmp, nil
}
