package controllers

import (
	"log"
	"meter-panel/store"
	"net/http"

	"k8s.io/apimachinery/pkg/labels"
)

const (
	key   string = "statefulset.kubernetes.io/pod-name"
	value string = "prometheus-kube-prometheus-0"
	// key2   string = "app"
	// value2 string = "jenkins"
)

func CpuViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	var tmp []byte
	if store.ProUseEnv {
		tmp = store.ListSingleClusterCpu(store.ProCfg[cluster])
	} else {
		labelset := labels.Set(map[string]string{key: value}).AsSelector()
		lister := store.AllLister.PodLister[cluster]
		if lister == nil {
			store.AddNewClusterResource(cluster)
			lister = store.AllLister.PodLister[cluster]
		}
		pl, _ := lister.List(labelset)
		if pl == nil {
			log.Println(EPrometheus)
			http.Error(w, EPrometheus, 400)
			return
		}
		log.Println(pl[0].Name)
		tmp = store.ListSingleClusterCpu(pl[0].Status.PodIP)
	}
	w.Write(tmp)
	// fakedata := "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{},\"value\":[1564972370.447,\"0.5997867933742264\"]}]}}"
	// w.Write([]byte(fakedata))
}

func MemViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	var tmp []byte
	if store.ProUseEnv {
		tmp = store.ListSingleClusterMem(store.ProCfg[cluster])
	} else {
		labelset := labels.Set(map[string]string{key: value}).AsSelector()
		lister := store.AllLister.PodLister[cluster]
		if lister == nil {
			store.AddNewClusterResource(cluster)
			lister = store.AllLister.PodLister[cluster]
		}
		pl, _ := lister.List(labelset)
		if pl == nil {
			log.Println(EPrometheus)
			http.Error(w, EPrometheus, 400)
			return
		}
		tmp = store.ListSingleClusterMem(pl[0].Status.PodIP)
	}
	w.Write(tmp)
	// fakedata := "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{},\"value\":[1563271125.722,\"33.71522723799342\"]}]}}"
	// w.Write([]byte(fakedata))
}

func AlertViews(w http.ResponseWriter, r *http.Request) {
	cluster := r.FormValue("cluster")
	var tmp []byte
	if store.ProUseEnv {
		tmp = store.ListSingleClusterAlerts(store.ProCfg[cluster])
	} else {
		labelset := labels.Set(map[string]string{key: value}).AsSelector()
		lister := store.AllLister.PodLister[cluster]
		if lister == nil {
			store.AddNewClusterResource(cluster)
			lister = store.AllLister.PodLister[cluster]
		}
		pl, _ := lister.List(labelset)
		if pl == nil {
			log.Println(EPrometheus)
			http.Error(w, EPrometheus, 400)
			return
		}
		log.Println("Get prometheus pod address" + pl[0].Status.PodIP)
		tmp = store.ListSingleClusterAlerts(pl[0].Status.PodIP)
	}
	w.Write(tmp)
	// fakedata := "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{},\"value\":[1563271125.722,\"33.71522723799342\"]}]}}"
	// w.Write([]byte(fakedata))
}
