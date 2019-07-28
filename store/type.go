package store

import (
	k8s "meter-panel/pkg/api/v1/k8s"
)

type NodeCache struct {
	StoreAllClusterNodeList map[string]k8s.HumanSingleClusterNodeList
}
