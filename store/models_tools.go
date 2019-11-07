package store

import (
	"fmt"
	"log"
	"meter-panel/pkg/api/v1/k8s/crd/application"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func isNeedDoSomeThingA(m meta) bool {
	_, ok := AllStore.ClientSet[m.clustername]

	if m.appname != "" && m.or == nil && ok {
		return true
	}
	return false
}

func isNeedDoSomeThingU(info1, info2 meta) bool {
	_, ok := AllStore.ClientSet[info1.clustername]
	_, ok2 := AllStore.ClientSet[info2.clustername]

	if info1.appname != info2.appname && ok && ok2 {
		return true
	}

	return false
}

func (m *Models) createApp(info meta) {
	if m.avoidConflict(info) {
		app := newApplication(info)
		_, err := AllStore.ClientSet[info.clustername].App(info.namespace).Create(app)
		if err != nil {
			log.Println(err)
		}
		m.Lock()
		delete(m.IsCreating, info.appname)
		m.Unlock()
	}
}

func (m *Models) avoidConflict(info meta) bool {
	_, err := AllStore.ClientSet[info.clustername].App(info.namespace).Get(info.appname, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	str := fmt.Sprintf("applications.app.k8s.io \"%v\" not found", info.appname)
	str2 := fmt.Sprintf("%v", err)
	if str2 == str {
	} else {
		log.Println("will not go to the waiting spaces")
		return false
	}

	for {
		if m.IsCreating[info.name] {
			_, err := AllStore.ClientSet[info.clustername].App(info.namespace).Get(info.appname, metav1.GetOptions{})
			if err != nil {
				log.Println(err)
			}
			str := fmt.Sprintf("applications.app.k8s.io \"%v\" not found", info.appname)
			str2 := fmt.Sprintf("%v", err)
			if str2 == str {
				log.Println(info.appname + " is creating, waiting for it")
				time.Sleep(time.Duration(1) * time.Second)
				continue
			}
			return true
		}

		m.Lock()
		m.IsCreating[info.appname] = true
		m.Unlock()
		return true
	}
}

func isEmptyApplication(info meta) bool {
	_, ok := AllStore.ClientSet[info.clustername]
	if !ok {
		return false
	}

	labelset := labels.Set(info.labels).AsSelector()
	svc, err := AllLister.SvcLister[info.clustername].List(labelset)
	if err != nil {
		log.Println(err)
	}
	sf, err := AllLister.StatefulSetLister[info.clustername].List(labelset)
	if err != nil {
		log.Println(err)
	}
	dep, err := AllLister.DeploymentLister[info.clustername].List(labelset)
	if err != nil {
		log.Println(err)
	}
	if len(svc)+len(sf)+len(dep) == 0 {
		log.Println(info.appname + " is a empty app")
		return true
	}
	log.Println(info.appname + " is not a empty app")
	return false
}

func newApplication(m meta) *application.Application {
	ac := strings.SplitN(m.labels[key], ".", -1)
	tmp := &application.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       appKind,
			APIVersion: appApiVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        ac[0],
			Namespace:   m.namespace,
			Labels:      map[string]string{appLabelKey: ac[0]},
			Annotations: map[string]string{appAnotation1: "", appAnotation2: ""},
		},
		Spec: application.ApplicationSpec{
			ComponentGroupKinds: []metav1.GroupKind{
				metav1.GroupKind{Group: depGroup, Kind: depKind},
				metav1.GroupKind{Group: sfGroup, Kind: sfKind},
				metav1.GroupKind{Group: svcGroup, Kind: svcKind},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: m.labels,
			},
		},
	}
	return tmp
}
