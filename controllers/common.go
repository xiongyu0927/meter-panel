package controllers

import (
	"errors"
	"meter-panel/pkg/api/v1/k8s/constome"
	cs "meter-panel/pkg/api/v1/k8s/crd/cluster"
	"meter-panel/store"

	"github.com/spf13/viper"
)

var Style constome.OrganizeData

const (
	EPrometheus string = "Can't find the correct prometheus address"
)

func init() {
	style := viper.GetString("PROJECT")
	switch style {
	case "CEB":
		Style = constome.NewCebStyle()
	case "YILI":
		Style = constome.NewYiLiStyle()
	}
}

func getProjectFromCluster(cluster string) ([]string, error) {
	cl := store.AllStore.ClusterStore.List()
	for _, v := range cl {
		t, ok := v.(*cs.Cluster)
		if ok && t.Name != "" && t.Name == cluster {
			return t.Finalizers, nil
		}
	}
	return nil, errors.New(cluster + " does not exit")
}
