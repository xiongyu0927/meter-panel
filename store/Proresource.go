package store

import (
	pro "meter-panel/pkg/api/v1/prometheus"
)

func GetSingleClusterCpu(cluster string) ([]byte, error) {
	a := StoreAllProConfigs[cluster]
	data, err := pro.ListSingleClusterCpu(a)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetSingleClusterMem(cluster string) ([]byte, error) {
	a := StoreAllProConfigs[cluster]
	data, err := pro.ListSingleClusterMem(a)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetSingleClusterAlert(cluster string) ([]byte, error) {
	a := StoreAllProConfigs[cluster]
	data, err := pro.ListSingleClusterAlert(a)
	if err != nil {
		return nil, err
	}
	return data, nil
}
