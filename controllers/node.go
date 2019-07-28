package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"meter-panel/store"
	"net/http"
)

// NodeViews is api of return node number and status
func NodeViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	tmp := store.StoreAllClusterNodeList[cluster]
	fmt.Println(tmp)
	data2, err := json.Marshal(tmp)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(data2)
	w.Write(data2)
}
