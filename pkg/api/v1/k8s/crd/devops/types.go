package devops

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/go-openapi/spec"
	fuzz "github.com/google/gofuzz"
	openapi "k8s.io/kube-openapi/pkg/common"
)

// ListEverything is a list options used to list all objects without any filtering.
var ListEverything = metav1.ListOptions{
	LabelSelector:   labels.Everything().String(),
	FieldSelector:   fields.Everything().String(),
	ResourceVersion: "0",
}

// LocalObjectReference simple local reference for local objects
// contains enough information to let you locate the
// referenced object inside the same namespace
// k8s.io/api/core/v1/types.go
type LocalObjectReference struct {
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// TODO: Add other useful fields. apiVersion, kind, uid?
	// +optional
	Name string `json:"name,omitempty"`
}

// HostPort defines a host/port  construct
type HostPort struct {
	// Host defines the host.
	// +optional
	Host string `json:"host"`

	// AccessURL defines an access URL for the tool
	// useful specially if the API url (host) is different than the
	// Access URL
	// +optional
	AccessURL string `json:"accessUrl"`
}

// GetAccessURL returns access url if set, otherwise returns Host
func (host HostPort) GetAccessURL() (url string) {
	url = host.Host
	if host.AccessURL != "" {
		url = host.AccessURL
	}
	return url
}

// ServiceStatusPhase defines the repo status
type ServiceStatusPhase string

func (phase ServiceStatusPhase) String() string {
	return string(phase)
}

const (
	// ServiceStatusPhaseCreating means the resource is creating
	ServiceStatusPhaseCreating ServiceStatusPhase = StatusCreating
	// ServiceStatusPhaseReady means the connection is ok
	ServiceStatusPhaseReady ServiceStatusPhase = StatusReady
	// ServiceStatusPhaseError means the connection is bad
	ServiceStatusPhaseError ServiceStatusPhase = StatusError
	// ServiceStatusPhaseWaitingToDelete means the resource will be deleted when no resourced reference
	ServiceStatusPhaseWaitingToDelete ServiceStatusPhase = StatusWaitingToDelete
	// ServiceStatusPhaseListTagError means registry list tag detail error
	ServiceStatusPhaseListTagError ServiceStatusPhase = StatusListTagError
	// ServiceStatusNeedsAuthorization needs authorization intervation from user.
	// Generally used on multi-step authorization schemes like oAuth2 etc
	ServiceStatusNeedsAuthorization ServiceStatusPhase = StatusNeedsAuthorization
)

// ServiceStatus defines the status of the service.
type ServiceStatus struct {
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
	// Conditions is a list of BindingCondition objects.
	// +optional
	Conditions []BindingCondition `json:"conditions"`
}

// String print a string for ServiceStatus
func (status ServiceStatus) String() string {
	var buff bytes.Buffer
	if status.Phase != "" {
		buff.WriteString(fmt.Sprintf(`Phase "%s",`, status.Phase))
	}
	if status.Message != "" {
		buff.WriteString(fmt.Sprintf(`Msg "%s",`, status.Message))
	}
	if status.Reason != "" {
		buff.WriteString(fmt.Sprintf(`Reason "%s",`, status.Reason))
	}
	buff.WriteString(fmt.Sprintf(`Conditions len(%d),`, len(status.Conditions)))
	return buff.String()
}

// CleanConditionsLastAttemptByOwner set nil to lastAttempt to conditions of a given owner
func (status ServiceStatus) CleanConditionsLastAttemptByOwner(owner string) ServiceStatus {
	if len(status.Conditions) > 0 {
		for i := 0; i < len(status.Conditions); i++ {
			if status.Conditions[i].Owner == owner {
				status.Conditions[i].LastAttempt = nil
			}
		}
	}
	return status
}

// HostPortStatus defines a status for a HostPort setting
type HostPortStatus struct {
	// StatusCode is the status code of http response
	// +optional
	StatusCode int `json:"statusCode"`
	// Response is the response of the http request.
	// +optional
	Response string `json:"response,omitempty"`
	// Version is the version of the http request.
	// +optional
	Version string `json:"version,omitempty"`
	// Delay means the http request will attempt later
	// +optional
	Delay *time.Duration `json:"delay,omitempty"`
	// Last time we probed the http request.
	// +optional
	LastAttempt *metav1.Time `json:"lastAttempt"`
	// Error Message of http request
	// +optional
	ErrorMessage string `json:"errorMessage"`
}

// Condition generic condition for devops objects
type Condition struct {
	// Type is the type of the condition.
	// +optional
	Type string `json:"type"`
	// Last time we probed the condition.
	// +optional
	LastAttempt *metav1.Time `json:"lastAttempt"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
	// Status is the status of the condition.
	// +optional
	Status string `json:"status"`
}

// BindingCondition defines the resource associated with the binding.
// The binding controller will check the status of the resource periodic and change it's status.
// The resource can be found by "name"+"type"+"binding's namespace"
type BindingCondition struct {
	// Name defines the name.
	// +optional
	Name string `json:"name"`
	// namespace defines the name.
	// +optional
	Namespace string `json:"namespace,omitempty"`
	// Type defines the type.
	// +optional
	Type string `json:"type"`
	// Last time we probed the condition.
	// +optional
	LastAttempt *metav1.Time `json:"lastAttempt"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
	// Status defines the status.
	// +optional
	Status string `json:"status,omitempty"`
	// Owner defins who own current condition
	// +optional
	Owner string `json:"owner,omitempty"`
}

type BindingConditions []BindingCondition

func (conds BindingConditions) Aggregate() (conditions *BindingCondition) {

	if len(conds) == 0 {
		return nil
	}

	for _, cond := range conds {
		// Any error condition will cause status of tbr changing to error
		if cond.Status == StatusError {
			return &cond
		}
	}
	return nil
}

// ReplaceBy remove the conditions that owned by owner and append the `appendConds` conditions
func (conds BindingConditions) ReplaceBy(owner string, appendConds []BindingCondition) []BindingCondition {

	// clean the conditions that controll by current controller
	result := make([]BindingCondition, 0, len(conds))

	for _, cond := range conds {
		if cond.Owner != owner {
			result = append(result, cond)
		}
	}

	//append new conds
	result = append(result, appendConds...)

	return result
}

// Get will return the conditions that owned by owner
func (conds BindingConditions) Get(owner string) []BindingCondition {

	result := make([]BindingCondition, 0, len(conds))

	for _, cond := range conds {
		if cond.Owner == owner {
			result = append(result, cond)
		}
	}

	return result
}

func (conds BindingConditions) Errors() error {
	if len(conds) == 0 {
		return nil
	}

	sb := strings.Builder{}
	for i, cond := range conds {
		if cond.Status == StatusError {
			sb.WriteString(fmt.Sprintf("Error-%d : %s \n", i, cond.Reason))
		}
	}

	msg := sb.String()
	if len(msg) == 0 {
		return nil
	}
	return errors.New(msg)
}

// DisplayName defines a set of readable names
type DisplayName struct {
	// EN is a human readable Chinese name.
	EN string `json:"en"`
	// ZH is a human readable English name.
	ZH string `json:"zh"`
}

// SecretKeySetRef reference of a set of username/api token keys in a Secret
type SecretKeySetRef struct {
	corev1.SecretReference `json:",inline"`
}

// Phaser returns its own phase
type Phaser interface {
	GetPhase() string
}

type KeyValueSet []KeyValue

type KeyValue struct {
	Name  string `json:"name" yaml:"name"`
	Value JsonV  `json:"value" yaml:"value"`
}

type JsonV struct {
	Type         JsonVType        `json:"type"`
	IntVal       int64            `json:"intVal"`
	StringVal    string           `json:"strVal"`
	BoolVal      bool             `json:"boolVal"`
	StringMapVal map[string]JsonV `json:"strMapVal"`
	ArrayVal     []JsonV          `json:"arrayVal"`
}

func (v *JsonV) UnmarshalJSON(value []byte) (err error) {
	v.Type, err = v.typ(value)
	if err != nil {
		return err
	}

	switch v.Type {
	case Null:
		return nil
	case Int:
		return json.Unmarshal(value, &v.IntVal)
	case String:
		return json.Unmarshal(value, &v.StringVal)
	case Bool:
		return json.Unmarshal(value, &v.BoolVal)
	case StringMap:
		return json.Unmarshal(value, &v.StringMapVal)
	case Arrary:
		return json.Unmarshal(value, &v.ArrayVal)
	default:
		return fmt.Errorf("UnKnown type when unmarshal json: %s", string(value))
	}
}
func (v *JsonV) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var nullV = &struct{}{}
	err := unmarshal(nullV)
	if err == nil && nullV == nil {
		v.Type = Null
		return nil
	}

	var intV int64
	err = unmarshal(&intV)
	if err == nil {
		v.Type = Int
		v.IntVal = intV
		return nil
	}

	var boolV bool
	err = unmarshal(&boolV)
	if err == nil {
		v.Type = Bool
		v.BoolVal = boolV
		return nil
	}

	var stringV string
	err = unmarshal(&stringV)
	if err == nil {
		v.Type = String
		v.StringVal = stringV
		return nil
	}

	var stringMapV = map[string]JsonV{}
	err = unmarshal(&stringMapV)
	if err == nil {
		v.Type = StringMap
		v.StringMapVal = stringMapV
		return nil
	}

	var arrayV = []JsonV{}
	err = unmarshal(&arrayV)
	if err == nil {
		v.Type = Arrary
		v.ArrayVal = arrayV
		return nil
	}

	return fmt.Errorf("UnKnown type when unmarshal yaml")
}

func (v JsonV) MarshalJSON() ([]byte, error) {
	switch v.Type {
	case Null:
		return json.Marshal(nil)
	case Int:
		return json.Marshal(v.IntVal)
	case String:
		return json.Marshal(v.StringVal)
	case Bool:
		return json.Marshal(v.BoolVal)
	case StringMap:
		return json.Marshal(v.StringMapVal)
	case Arrary:
		return json.Marshal(v.ArrayVal)
	default:
		return []byte{}, fmt.Errorf("impossible V.Type: %#v", v.Type)
	}
}

func (v JsonV) MarshalYAML() (interface{}, error) {
	switch v.Type {
	case Null:
		return nil, nil
	case Int:
		return v.IntVal, nil
	case String:
		return v.StringVal, nil
	case Bool:
		return v.BoolVal, nil
	case StringMap:
		return v.StringMapVal, nil
	case Arrary:
		return v.ArrayVal, nil
	default:
		return nil, fmt.Errorf("impossible V.Type: %#v", v.Type)

	}
}

func (v *JsonV) typ(value []byte) (JsonVType, error) {
	start := value[0]
	if start == '"' {
		return String, nil
	}
	if start == '{' {
		return StringMap, nil
	}
	if start == '[' {
		return Arrary, nil
	}

	str := string(value)
	if str == "false" || str == "true" {
		return Bool, nil
	}

	if str == "null" {
		return Null, nil
	}

	_, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return Int, nil
	}

	return Null, fmt.Errorf("UnKnown type of value: %s", string(value))
}

func (v *JsonV) Fuzz(c fuzz.Continue) {
	if v == nil {
		return
	}
	if c.RandBool() {
		v.Type = Bool
		v.IntVal = 0
		c.Fuzz(&v.BoolVal)
	} else {
		v.Type = String
		v.IntVal = 0
		c.Fuzz(&v.StringVal)
	}
}

//OpenAPIDefinition according https://github.com/kubernetes/kube-openapi/tree/release-1.10/pkg/generators
func (_ JsonV) OpenAPIDefinition() openapi.OpenAPIDefinition {
	return openapi.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				// https://github.com/kubernetes/kube-openapi/blob/release-1.10/pkg/util/proto/document.go#L239
				Type: []string{},
			},
		},
	}
}

// func (_ JsonV) OpenAPISchemaType() []string { return []string{} }
// func (_ JsonV) OpenAPISchemaFormat() string { return "" }

// JsonVType represents the stored type of IntOrString.
type JsonVType int

const (
	// Null indicates the type of Object is Null
	Null JsonVType = iota
	Int
	String
	Bool
	StringMap
	Arrary
)
