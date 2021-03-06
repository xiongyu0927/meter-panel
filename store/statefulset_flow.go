package store

import (
	"fmt"
	"log"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Models) AddStatefulSetFlow(sf *appsv1.StatefulSet) {
	info := getStatefulSetMeta(sf)
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
			m.updateStatefulSetOR(info)
		} else {
			updateApplication(info)
			log.Println("update deployment" + info.name)
			m.updateStatefulSetOR(info)
		}
	}
}

func (m *Models) UpdateStatefulSetFlow(sf1, sf2 *appsv1.StatefulSet) {
	info1 := getStatefulSetMeta(sf1)
	info2 := getStatefulSetMeta(sf2)
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
			m.updateStatefulSetOR(info2)
		} else {
			updateApplication(info2)
			m.updateStatefulSetOR(info2)
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

func (m *Models) DeleteStatefulSetFlow(sf *appsv1.StatefulSet) {
	info := getStatefulSetMeta(sf)
	if isEmptyApplication(info) {
		err := AllStore.ClientSet[info.clustername].App(info.namespace).Delete(info.appname, &metav1.DeleteOptions{})
		if err != nil {
			log.Println(err)
		}
	}
}

func (m *Models) updateStatefulSetOR(info meta) {
	for i := 0; i < 3; i++ {
		str := changeStatefulSet(info)
		if !strings.Contains(str, modifiederr) {
			log.Println("update deployment ownerreference succssed")
			break
		}
	}
}

func changeStatefulSet(info meta) string {
	v, err := AllStore.ClientSet[info.clustername].App(info.namespace).Get(info.appname, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	t := true
	f := false
	sf, err := AllLister.ClientSet[info.clustername].AppsV1().StatefulSets(info.namespace).Get(info.name, metav1.GetOptions{})
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
	if len(sf.OwnerReferences) != 0 {
		for i, v := range sf.OwnerReferences {
			if v.Kind == appKind {
				sf.OwnerReferences[i] = this
			} else {
				sf.OwnerReferences = append(sf.OwnerReferences, this)
			}
		}
	} else {
		sf.OwnerReferences = append(sf.OwnerReferences, this)
	}
	sf.Labels[key] = info.appname + "." + info.namespace
	sf.Spec.Template.Labels[key] = info.appname + "." + info.namespace
	AllLister.ClientSet[info.clustername].AppsV1().StatefulSets(info.namespace).Update(sf)
	if err != nil {
		log.Println(err)
		str := fmt.Sprintf("%v", err)
		return str
	}

	return ""
}

func getStatefulSetMeta(sf *appsv1.StatefulSet) meta {
	var app, cnm string

	label := sf.GetLabels()
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
		name:        sf.GetName(),
		namespace:   sf.GetNamespace(),
		labels:      neededLabel,
		or:          sf.GetOwnerReferences(),
		appname:     app,
		clustername: cnm,
		kind:        sfKind,
	}
	return m
}
