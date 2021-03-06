package store

import (
	"fmt"
	"log"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (m *Models) AddDeploymentFlow(dep *appsv1.Deployment) {
	info := getDeploymentMeta(dep)
	if isNeedDoSomeThingA(info) {
		_, err := AllStore.ClientSet[info.clustername].App(info.namespace).Get(info.appname, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
		}
		str := fmt.Sprintf("applications.app.k8s.io \"%v\" not found", info.appname)
		str2 := fmt.Sprintf("%v", err)
		if str2 == str {
			log.Println("create application " + info.name)
			m.createApp(info)
			m.updateDeploymentOR(info)
		} else {
			updateApplication(info)
			log.Println("update deployment " + info.name)
			m.updateDeploymentOR(info)
		}
	}
}

func (m *Models) UpdateDeploymentFlow(dep1, dep2 *appsv1.Deployment) {
	info1 := getDeploymentMeta(dep1)
	info2 := getDeploymentMeta(dep2)
	if isNeedDoSomeThingU(info1, info2) {
		// 处理新标签
		_, err := AllStore.ClientSet[info2.clustername].App(info2.namespace).Get(info2.appname, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
		}
		str := fmt.Sprintf("applications.app.k8s.io \"%v\" not found", info2.appname)
		str2 := fmt.Sprintf("%v", err)
		if str2 == str {
			log.Println("create application " + info2.appname)
			m.createApp(info2)
			m.updateDeploymentOR(info2)
		} else {
			updateApplication(info2)
			m.updateDeploymentOR(info2)
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

func (m *Models) DeleteDeploymentFlow(dep *appsv1.Deployment) {
	info := getDeploymentMeta(dep)
	if isEmptyApplication(info) {
		err := AllStore.ClientSet[info.clustername].App(info.namespace).Delete(info.appname, &metav1.DeleteOptions{})
		if err != nil {
			log.Println(err)
		}
	}
}

func getDeploymentMeta(dep *appsv1.Deployment) meta {
	var app, cnm string

	label := dep.GetLabels()
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
		name:        dep.GetName(),
		namespace:   dep.GetNamespace(),
		labels:      neededLabel,
		or:          dep.GetOwnerReferences(),
		appname:     app,
		clustername: cnm,
		kind:        depKind,
	}
	return m
}

func (m *Models) updateDeploymentOR(info meta) {
	for i := 0; i < 3; i++ {
		err := changeDeployment(info)
		if err == nil {
			break
		}
	}
}

func changeDeployment(info meta) error {
	v, err := AllStore.ClientSet[info.clustername].App(info.namespace).Get(info.appname, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	t := true
	f := false
	dep, err := AllLister.ClientSet[info.clustername].AppsV1().Deployments(info.namespace).Get(info.name, metav1.GetOptions{})
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
	if len(dep.OwnerReferences) != 0 {
		for i, v := range dep.OwnerReferences {
			if v.Kind == appKind {
				dep.OwnerReferences[i] = this
			} else {
				dep.OwnerReferences = append(dep.OwnerReferences, this)
			}
		}
	} else {
		dep.OwnerReferences = append(dep.OwnerReferences, this)
	}
	dep.Labels[key] = info.appname + "." + info.namespace
	dep.Spec.Template.Labels[key] = info.appname + "." + info.namespace
	_, err = AllLister.ClientSet[info.clustername].AppsV1().Deployments(info.namespace).Update(dep)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
