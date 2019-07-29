package tools

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
)

// Request define the request info
type Request struct {
	Methoud string
	Host    string
	Path    string
	// IsHttps shuold be https or http
	IsHTTPS   string
	BearToken string
	// Chan is used Transport k8s events when watching the api
	Chan chan map[string]interface{}
}

// tr is InsecureSkipVerify
var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var ChanData = make(map[string][]byte)

// Get function is httpclient do get request
func (its *Request) Get() ([]byte, error) {
	request, err := http.NewRequest("GET", its.IsHTTPS+"://"+its.Host+its.Path, nil)
	if err != nil {
		return nil, err
	}

	// add BearToken auth
	if its.BearToken != "" {
		request.Header.Add("Authorization", "Bearer "+its.BearToken)
	}

	client := http.Client{}
	// add InsecureSkipVerify
	if its.IsHTTPS == "https" {
		client.Transport = tr
	}

	// execute this request
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// get,read and return
	defer resp.Body.Close()
	tmp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return tmp, nil
}
