package controllers

import (
	"encoding/json"
	"log"
	"meter-panel/store"
	"net/http"
)

func CpuViews(w http.ResponseWriter, r *http.Request) {
	// cluster := r.FormValue("cluster")
	// tmp, err := store.GetSingleClusterCpu(cluster)
	// if err != nil {
	// 	log.Println(err)
	// }
	// data2, err := json.Marshal(tmp)
	// if err != nil {
	// 	log.Println(err)
	// }
	fakedata := "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{},\"value\":[1564972370.447,\"0.5997867933742264\"]}]}}"
	w.Write([]byte(fakedata))
}

func MemViews(w http.ResponseWriter, r *http.Request) {
	// cluster := r.FormValue("cluster")
	// tmp, err := store.GetSingleClusterMem(cluster)
	// if err != nil {
	// 	log.Println(err)
	// }
	// data2, err := json.Marshal(tmp)
	// if err != nil {
	// 	log.Println(err)
	// }
	fakedata := "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{},\"value\":[1563271125.722,\"33.71522723799342\"]}]}}"
	w.Write([]byte(fakedata))
}

func AlertViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	tmp, err := store.GetSingleClusterAlert(cluster)
	if err != nil {
		log.Println(err)
	}
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
	}
	w.Write(data2)
}
