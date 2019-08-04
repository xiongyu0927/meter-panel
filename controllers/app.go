package controllers

import (
	"encoding/json"
	"log"
	"meter-panel/store"
	"net/http"
)

// AppViews is api of return node number and status
func AppViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	tmp := store.TmpGetAPP(cluster)
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
	}
	w.Write(data2)
}
