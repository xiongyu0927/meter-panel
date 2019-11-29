package store

import (
	"fmt"
	"log"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Models) AddServiceFlow(svc *v1.Service) {
	info := getServiceMeta(svc)
	if isNeedDoSomeThingA(info) {
		_, err := AllStore.ClientSet[info.clustername].App(info.namespace).Get(info.appname, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
		}
		str := fmt.Sprintf("applications.app.k8s.io \"%v\" not found", info.appname)
		str2 := fmt.Sprintf("%v", err)
		if str2 == str {
			log.Println("create application" + info.name)
			m.createApp(info)
			m.updateServiceOR(info)
		} else {
			updateApplication(info)
			log.Println("update deployment" + info.name)
			m.updateServiceOR(info)
		}
	}
}

func (m *Models) UpdateServiceFlow(svc1, svc2 *v1.Service) {
	info1 := getServiceMeta(svc1)
	info2 := getServiceMeta(svc2)
	if isNeedDoSomeThingU(info1, info2) {
		// 处理新标签
		_, err := AllStore.ClientSet[info2.clustername].App(info2.namespace).Get(info2.appname, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
		}
		str := fmt.Sprintf("applications.app.k8s.io \"%v\" not found", info2.appname)
		str2 := fmt.Sprintf("%v", err)
		if str2 == str {
			log.Println("create application" + info2.appname)
			m.createApp(info2)
			m.updateServiceOR(info2)
		} else {
			updateApplication(info2)
			m.updateServiceOR(info2)
		}

		// 处理旧标签
		_, err = AllStore.ClientSet[info1.clustername].App(info1.namespace).Get(info1.appname, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
			// return
		}
		str = fmt.Sprintf("applications.app.k8s.io \"%v\" not found", info1.appname)
		str2 = fmt.Sprintf("%v", err)
		// 旧label是否还有Application，没有则结束，
		// 有则判断该application是否还有意义，有则结束，没有则删除
		if str2 == str {
			log.Println("old application doesn't have needn't do any thing")
		} else {
			if isEmptyApplication(info1) {
				err := AllStore.ClientSet[info1.clustername].App(info1.namespace).Delete(info1.appname, &metav1.DeleteOptions{})
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func (m *Models) DeleteServiceFlow(svc *v1.Service) {
	info := getServiceMeta(svc)
	if isEmptyApplication(info) {
		err := AllStore.ClientSet[info.clustername].App(info.namespace).Delete(info.appname, &metav1.DeleteOptions{})
		if err != nil {
			log.Println(err)
		}
	}
}

func getServiceMeta(svc *v1.Service) meta {
	var app, cnm string

	label := svc.GetLabels()
	v := label[key2]
	neededLabel := make(map[string]string)
	neededLabel[key2] = v

	ac := strings.SplitN(v, ".", -1)
	if ac != nil && len(ac) == 2 {
		if _, ok := AllStore.ClientSet[ac[1]]; ok {
			cnm = ac[1]
			app = ac[0]
		}
	}

	m := meta{
		name:        svc.GetName(),
		namespace:   svc.GetNamespace(),
		labels:      neededLabel,
		or:          svc.GetOwnerReferences(),
		appname:     app,
		clustername: cnm,
		kind:        svcKind,
	}
	return m
}

func (m *Models) updateServiceOR(info meta) {
	for i := 0; i < 3; i++ {
		str := changeService(info)
		if !strings.Contains(str, modifiederr) {
			log.Println("update deployment ownerreference succssed")
			break
		}
	}
}

func changeService(info meta) string {
	v, err := AllStore.ClientSet[info.clustername].App(info.namespace).Get(info.appname, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	t := true
	f := false
	svc, err := AllLister.ClientSet[info.clustername].CoreV1().Services(info.namespace).Get(info.name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	this := metav1.OwnerReference{
		APIVersion:         appApiVersion,
		Kind:               appKind,
		Name:               info.appname,
		UID:                v.UID,
		Controller:         &f,
		BlockOwnerDeletion: &t,
	}
	// 如果有Application的or，则直接覆盖
	if len(svc.OwnerReferences) != 0 {
		for i, v := range svc.OwnerReferences {
			if v.Kind == appKind {
				svc.OwnerReferences[i] = this
			} else {
				svc.OwnerReferences = append(svc.OwnerReferences, this)
			}
		}
	} else {
		svc.OwnerReferences = append(svc.OwnerReferences, this)
	}
	svc.Labels[key] = info.appname + "." + info.namespace
	_, err = AllLister.ClientSet[info.clustername].CoreV1().Services(info.namespace).Update(svc)
	if err != nil {
		log.Println(err)
		str := fmt.Sprintf("%v", err)
		return str
	}

	return ""
}
