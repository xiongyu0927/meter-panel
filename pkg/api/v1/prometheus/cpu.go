package prometheus

// import (
// 	"log"
// 	"meter-panel/tools"
// )
//
// var ProRequest = tools.Request{
// 	Methoud: "GET",
// 	Host:    "",
// 	Path:    "",
// 	//IsHttps shuold be https or http
// 	IsHTTPS:   "http",
// 	BearToken: "",
// }
//
// func ListSingleClusterCpu(address string) []byte {
// 	ProRequest.Host = address
// 	ProRequest.Path = "/api/v1/query?query=sum(100%20-%20(avg%20by%20(instance)%20(rate(node_cpu%7Bjob%3D%22node-exporter%22%2Cmode%3D%22idle%22%7D%5B5m%5D))%20*%20100))%20%2F%20count(node_cpu%7Bjob%3D%22node-exporter%22%2Cmode%3D%22idle%22%7D)"
// 	data, err := ProRequest.Get()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return data
// }
