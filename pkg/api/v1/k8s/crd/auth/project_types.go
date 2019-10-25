/*
Copyright 2018 Alauda.io.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProjectClusters is the attribute for ProjectSpec
type ProjectClusters struct {
	// Name of the cluster
	Name string `json:"name"`
	// Type of the cluster
	Type string `json:"type"`
	// Quota store the quota info for Project
	Quota corev1.ResourceList `json:"quota"`
}

func (pc *ProjectClusters) UnmarshalJSON(value []byte) error {
	var projectCluster struct {
		// Name of the cluster
		Name string `json:"name"`
		// Type of the cluster
		Type string `json:"type"`
		// Quota store the quota info for Project
		Quota map[string]json.RawMessage `json:"quota"`
	}
	json.Unmarshal(value, &projectCluster)

	var quota = corev1.ResourceList{}

	for resourceName, quantityBytes := range projectCluster.Quota {
		quantity := resource.Quantity{}
		// only record the correct quantity
		if err := json.Unmarshal([]byte(string(quantityBytes)), &quantity); err == nil {
			quota[corev1.ResourceName(resourceName)] = quantity
		}
	}

	*pc = ProjectClusters{
		Name:  projectCluster.Name,
		Type:  projectCluster.Type,
		Quota: quota,
	}

	return nil
}

// ProjectSpec defines the desired state of Project
type ProjectSpec struct {
	// Clusters contains the clusters associated with this Project
	Clusters []*ProjectClusters `json:"clusters,omitempty"`
}

// ProjectStatus defines the observed state of Project
type ProjectStatus struct {
	// Phase record the state of project
	Phase string `json:"phase,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Project is the Schema for the projects API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProjectSpec   `json:"spec,omitempty"`
	Status ProjectStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProjectList contains a list of Project
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}

//
// func init() {
// 	SchemeBuilder.Register(&Project{}, &ProjectList{})
// }
//
// // SetDefaults make sure that Labels and Annotations are initialized
// func (object *Project) SetDefaults() {
// 	if len(object.Labels) == 0 {
// 		object.Labels = map[string]string{}
// 	}
//
// 	if len(object.Annotations) == 0 {
// 		object.Annotations = map[string]string{}
// 	}
// }
