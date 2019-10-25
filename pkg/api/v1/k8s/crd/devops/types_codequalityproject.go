package devops

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CodeQualityProject save CodeQualityTool Project info
type CodeQualityProject struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Specification of the desired behavior of the CodeQualityProject.
	// +optional
	Spec CodeQualityProjectSpec `json:"spec"`
	// Most recently observed status of the CodeQualityProject.
	// Populated by the system.
	// Read-only.
	// +optional
	Status CodeQualityProjectStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CodeQualityProjectList is a list of CodeQualityProject objects.
type CodeQualityProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CodeQualityProject `json:"items"`
}

// CodeQualityProject presents CodeQualityProject's spec
type CodeQualityProjectSpec struct {
	// CodeQualityTool defines the CodeQualityTool in spec
	CodeQualityTool LocalObjectReference `json:"codeQualityTool"`
	// CodeQualityBinding defines the CodeQualityBinding in spec
	CodeQualityBinding LocalObjectReference `json:"codeQualityBinding"`
	// CodeRepository defines the CodeRepository in spec
	CodeRepository LocalObjectReference `json:"codeRepository"`
	// Project defines CodeQualityProject info
	Project CodeQualityProjectInfo `json:"project"`
}

// CodeQualityProjectInfo presents CodeQualityProject info
type CodeQualityProjectInfo struct {
	// ProjectKey defines key in CodeQualityProjectInfo
	ProjectKey string `json:"projectKey"`
	// ProjectName defines display name in CodeQualityProjectInfo
	ProjectName string `json:"projectName"`
	// CodeAddress defines code address in CodeQualityProjectInfo
	CodeAddress string `json:"codeAddress"`
	// LastAnalysis defines the last analysis date of this project
	LastAnalysis *metav1.Time `json:"lastAnalysisDate"`
}

// CodeQualityProjectInfo presents CodeQualityProject status
type CodeQualityProjectStatus struct {
	// Current condition of the service.
	// One of: "Creating" or "Ready" or "Error" or "WaitingToDelete".
	// +optional
	Phase ServiceStatusPhase `json:"phase"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
	// LastUpdate is the latest time when updated the service.
	// +optional
	LastUpdate *metav1.Time `json:"lastUpdated"`
	// HTTPStatus is http status of the service.
	// +optional
	HTTPStatus *HostPortStatus `json:"http,omitempty"`
	// Conditions defines analyze info
	CodeQualityConditions []CodeQualityCondition `json:"conditions"`
}

// CodeQualityCondition presents CodeQualityProject analyze info
type CodeQualityCondition struct {
	BindingCondition
	// Branch defines analyze code branch, default is master
	Branch string `json:"branch"`
	// IsMain defines whether the branch is the main branch
	IsMain bool `json:"isMain"`
	// QualityGate defines project use which quality gate
	QualityGate string `json:"qualityGate"`
	// Public defines project visible
	Visibility string `json:"visibility"`
	// Metrics define a series of metrics of this project
	Metrics map[string]CodeQualityAnalyzeMetric `json:"metrics"`
}

// CodeQualityAnalyzeResult present CodeQualityProject analyze result
type CodeQualityAnalyzeMetric struct {
	// Name defines the name of this metric
	Name string `json:"name"`
	// Value defines the value of this metric
	Value string `json:"value"`
	// Level defines the level of the value
	// +optional
	Level string `json:"level,omitempty"`
}
