package controllers

import "meter-panel/pkg/api/v1/k8s/constome"

var Style constome.OrganizeData

const (
	EPrometheus string = "Can't find the correct prometheus address"
)

func init() {
	Style = constome.NewCebStyle()
}
