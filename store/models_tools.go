package store

import (
	"fmt"
	"log"
	"meter-panel/pkg/api/v1/k8s/crd/application"
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
				time.Sleep(1 * time.Second)
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
	tmp := &application.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       appKind,
			APIVersion: appApiVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        m.appname,
			Namespace:   m.namespace,
			Labels:      map[string]string{key: m.appname},
			Annotations: map[string]string{appAnotation1: "", appAnotation2: ""},
		},
		Spec: application.ApplicationSpec{
			ComponentGroupKinds: []metav1.GroupKind{
				metav1.GroupKind{Group: depGroup, Kind: depKind},
				metav1.GroupKind{Group: sfGroup, Kind: sfKind},
				metav1.GroupKind{Group: svcGroup, Kind: svcKind},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{key: m.appname + "." + m.namespace},
			},
		},
	}
	return tmp
}

// 当我创建好应用后，appliction应该是三种资源都在group里，但是平台更新一次后，group里只会显示
// 他找得到的那种资源，其他的会被删掉。所以当新增同一个Application下的其他类型资源后，我需要考虑
// 这个Application是否已经被平台更新过了
func updateApplication(m meta) {
	app, err := AllStore.ClientSet[m.clustername].App(m.namespace).Get(m.appname, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}

	for _, v := range app.Spec.ComponentGroupKinds {
		if v.Kind == m.kind {
			return
		}
	}

	app.Spec.ComponentGroupKinds = []metav1.GroupKind{
		metav1.GroupKind{Group: depGroup, Kind: depKind},
		metav1.GroupKind{Group: sfGroup, Kind: sfKind},
		metav1.GroupKind{Group: svcGroup, Kind: svcKind},
	}

	_, err = AllStore.ClientSet[m.clustername].App(m.namespace).Update(app)
	if err != nil {
		log.Println(err)
	}
}
