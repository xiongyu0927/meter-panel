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
	al := store.AllStore.AppStore[cluster].List()
	tmp := Style.OrganzieApplicationList(cluster, al)

	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}
	w.Write(data2)
}
