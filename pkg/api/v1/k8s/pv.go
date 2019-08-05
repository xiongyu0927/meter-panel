package k8s

import (
	"encoding/json"
	"log"
	"meter-panel/configs"
	"strconv"
	"strings"
)

var (
	pvslist                      SingleClusterPvList
	NilHumanSingleClusterPvsList HumanSingleClusterPvsList
	_HumanAllClusterPvsList      HumanAllClusterPvsList = make(map[string]HumanSingleClusterPvsList)
)

func ListAllClusterPvs(k8sconfigs configs.HumanAllK8SConfigs) (HumanAllClusterPvsList, error) {
	for k, v1 := range k8sconfigs {
		tmp, err := ListSingleClusterPvs(v1)
		if err != nil {
			return nil, err
		}
		_HumanAllClusterPvsList[k] = tmp
	}
	return _HumanAllClusterPvsList, nil
}

func ListSingleClusterPvs(k8sconfig configs.HumanSingleK8sConfigs) (HumanSingleClusterPvsList, error) {
	K8sRequest.Host = k8sconfig.EndPoint
	K8sRequest.BearToken = k8sconfig.Token
	K8sRequest.Path = "/api/v1/persistentvolumes"
	data, err := K8sRequest.Get()
	if err != nil {
		return NilHumanSingleClusterPvsList, err
	}
	err = json.Unmarshal(data, &pvslist)

	if err != nil {
		return NilHumanSingleClusterPvsList, err
	}
	pvstatus := make(map[string]string)
	var initstorage int
	PvsDetail(pvslist.Items, pvstatus, &initstorage)

	if pvslist.Metadata.Continue != "" {
		K8sRequest.Path = "/api/v1/persistentvolumes?limit=500&continue=" + podslist.Metadata.Continue
		data, err = K8sRequest.Get()
		if err != nil {
			return NilHumanSingleClusterPvsList, err
		}
		err = json.Unmarshal(data, &pvslist)
		if err != nil {
			return NilHumanSingleClusterPvsList, err
		}

		PvsDetail(pvslist.Items, pvstatus, &initstorage)
	}
	log.Println(initstorage)
	storage := ToGMK(initstorage)
	log.Println(storage)
	var tmp2 = HumanSingleClusterPvsList{
		PvStatus: pvstatus,
		AllStore: storage,
	}
	return tmp2, nil
}

func PvsDetail(item []pv, pvstatus map[string]string, initstorage *int) {
	for _, v2 := range item {
		pvstatus[v2.Metadata.Name] = v2.Spec.Capacity.Storage
		num := v2.Spec.Capacity.Storage
		x := ToB(num)
		*initstorage = *initstorage + x
	}
}

func ToB(num string) int {
	if strings.Contains(num, "Gi") {
		tmp := strings.SplitN(num, "G", -1)
		x, err := strconv.Atoi(tmp[0])
		if err != nil {
			log.Println(err)
		}
		x = x << 30
		return x
	}

	if strings.Contains(num, "Mi") {
		tmp := strings.SplitN(num, "M", -1)
		x, err := strconv.Atoi(tmp[0])
		if err != nil {
			log.Println(err)
		}
		x = x << 20
		return x
	}

	if strings.Contains(num, "Ki") {
		tmp := strings.SplitN(num, "K", -1)
		x, err := strconv.Atoi(tmp[0])
		if err != nil {
			log.Println(err)
		}
		x = x << 10
		return x
	}

	if strings.Contains(num, "G") {
		tmp := strings.SplitN(num, "G", -1)
		x, err := strconv.Atoi(tmp[0])
		if err != nil {
			log.Println(err)
		}
		x = x * 1000 * 1000 * 1000
		return x
	}

	if strings.Contains(num, "M") {
		tmp := strings.SplitN(num, "M", -1)
		x, err := strconv.Atoi(tmp[0])
		if err != nil {
			log.Println(err)
		}
		x = x * 1000 * 1000
		return x
	}

	if strings.Contains(num, "K") {
		tmp := strings.SplitN(num, "K", -1)
		x, err := strconv.Atoi(tmp[0])
		if err != nil {
			log.Println(err)
		}
		x = x * 1000
		return x
	}

	return 0
}

func ToGMK(s int) string {
	suffix := ""
	b := s
	if s > (1 << 30) {
		suffix = "Gi"
		b = s / (1 << 30)
	} else if s > (1 << 20) {
		suffix = "Mi"
		b = s / (1 << 20)
	} else if s > (1 << 10) {
		suffix = "Ki"
		b = s / (1 << 10)
	}
	tmp := strconv.Itoa(b) + suffix
	return tmp
}
