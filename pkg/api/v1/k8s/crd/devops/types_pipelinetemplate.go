package devops

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineTaskTemplate specified a task for a pipeline
type PipelineTaskTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Specification of the PipelineTaskTemplate
	Spec PipelineTaskTemplateSpec `json:"spec"`
	// Status
	// +optional
	Status TemplateStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineTaskTemplateList is list of PipelineTaskTemplate
type PipelineTaskTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata"`

	Items []PipelineTaskTemplate `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterPipelineTaskTemplate is kind of cluster PipelineTaskTemplate
type ClusterPipelineTaskTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Specification of the PipelineTaskTemplate
	Spec PipelineTaskTemplateSpec `json:"spec"`
	// Status
	// +optional
	Status TemplateStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterPipelineTaskTemplateList is a list of ClusterPipelineTaskTemplate
type ClusterPipelineTaskTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata"`

	Items []ClusterPipelineTaskTemplate `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineExportedVariables represent the exports for the pipelinetempalte exports
type PipelineExportedVariables struct {
	metav1.TypeMeta `json:",inline"`
	Values          []GlobalParameter `json:"values"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ExportShowOptions struct {
	metav1.TypeMeta `json:",inline"`
	// SourceType specifies the type of the pipeline scm.
	TaskName    string `json:"taskName"`
	FormatValue string `json:"formatvalue"`
}

// PipelineDependency for PipelineTaskTemplate dependency on jenkins plugins
type PipelineDependency struct {
	// Plugins hold Jenkins plugins dependency for the specific task
	// +optional
	Plugins []JenkinsPlugin `json:"plugins"`
}

// JenkinsPlugin is Jenkins plugin info
type JenkinsPlugin struct {
	// Name is the name of plugin
	Name string `json:"name"`
	// Version is the version of plugin
	Version string `json:"version"`
}

// PipelineTemplateTaskEngine describe the kind of engine used for the PipelineTaskTemplate
type PipelineTemplateTaskEngine string

const (
	// PipelineTaskTemplateEngineGoTemplate go template for rendering
	PipelineTaskTemplateEngineGoTemplate = "gotpl"
)

// PipelineTaskTemplateSpec represents PipelineTaskTemplate's specs
type PipelineTaskTemplateSpec struct {
	// Engine the way of how to render taskTemplate
	// +optinal
	Engine PipelineTemplateTaskEngine `json:"engine"`
	// Agent indicates where the task should be running
	// +optional
	Agent *JenkinsAgent `json:"agent,omitempty"`
	// Body task template body
	Body string `json:"body"`
	// Exports all envrionments will be exports
	// +optional
	Exports []GlobalParameter `json:"exports,omitempty"`
	// Parameters that will be use in running
	// +optional
	Parameters []PipelineParameter `json:"parameters,omitempty"`
	// Arguments the task template's arguments
	// +optional
	Arguments []PipelineTaskArgument `json:"arguments,omitempty"`
	// Dependencies indicates plugins denpendencies of task
	// +optional
	Dependencies *PipelineDependency `json:"dependencies,omitempty"`
}

// GlobalParameter for export
type GlobalParameter struct {
	// Name the name of parameter
	Name string `json:"name"`
	// Description description of parameter
	// +optional
	Description *I18nName `json:"description"`
}

// PipelineTaskArgument sepcified a arugment for PipelineTaskTemplate
type PipelineTaskArgument struct {
	// Name the name of task
	Name string `json:"name"`
	// Schema schema of task
	Schema PipelineTaskArgumentSchema `json:"schema"`
	// Display display of task
	Display PipelineTaskArgumentDisplay `json:"display"`
	// Required indicate whether required
	// +optional
	Required bool `json:"required"`
	// Default default value of arugment
	// +optional
	Default string `json:"default"`
	// Validation validation of arugment
	// +optional
	Validation *PipelineTaskArgumentValidation `json:"validation,omitempty"`
	// Relation relation between arguments
	// +optional
	Relation []PipelineTaskArgumentAction `json:"relation,omitempty"`
}

// PipelineTaskArgumentValidation for task arument validation
type PipelineTaskArgumentValidation struct {
	// Pattern pattern of validation
	// +optional
	Pattern string `json:"pattern"`
	// MaxLength maxLength of this field
	// +optional
	MaxLength int `json:"maxLength"`
}

// PipelineTaskArgumentAction action for task argument
type PipelineTaskArgumentAction struct {
	// Action action for task argument
	Action string `json:"action"`
	// When time condition for task execution
	When PipelineTaskArgumentWhen `json:"when"`
}

// PipelineTaskArgumentWhen action time config
type PipelineTaskArgumentWhen struct {
	// // +optional
	// *RelationWhenItem `json:",inline"`

	// Name name of when
	// +optional
	Name string `json:"name,omitempty"`
	// Value value of when
	// +optional
	Value bool `json:"value,omitempty"`
	// +optional
	All []RelationWhenItem `json:"all,omitempty"`
	// +optional
	Any []RelationWhenItem `json:"any,omitempty"`
}

// RelationWhenItem all condition in pipelinetaskarguementwhtn
type RelationWhenItem struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

// PipelineTemplateEngine describe the kind of engine used for the PipelineTemplate
type PipelineTemplateEngine string

const (
	// PipelineTemplateEngineGraph render template as the graph form
	PipelineTemplateEngineGraph PipelineTemplateEngine = "graph"
)

// +genclient
// +genclient:method=Preview,verb=create,subresource=preview,input=JenkinsfilePreviewOptions,result=JenkinsfilePreview
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineTemplate specified a jenkinsFile template
type PipelineTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Spec specification of PipelineTemplate
	Spec PipelineTemplateSpec `json:"spec"`
	// Status
	// +optional
	Status TemplateStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineTemplateList is a list of PipelineTemplate
type PipelineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items items of PipelineTemplates
	Items []PipelineTemplate `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterPipelineTemplate specified a cluster kind of PipelineTemplate
type ClusterPipelineTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Spec specification of PipelineTemplate
	Spec PipelineTemplateSpec `json:"spec"`
	// Status
	// +optional
	Status TemplateStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterPipelineTemplateList is a list of ClusterPipelineTemplate
type ClusterPipelineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items items of ClusterPipelineTemplate
	Items []ClusterPipelineTemplate `json:"items"`
}

type TemplatePhase string

const (
	TemplateReady       TemplatePhase = "Ready"
	TemplateTerminating TemplatePhase = "Terminating"
)

// TemplateStatus defines template status
// includes pipelinetemplate and pipelinetasktemplate
type TemplateStatus struct {
	// +optional
	Phase TemplatePhase `json:"phase"`
}

// PipelineTemplateSpec represents PipelineTemplate's specs
type PipelineTemplateSpec struct {
	// Engine the way how to render PipelineTemplate
	// +optinal
	Engine PipelineTemplateEngine `json:"engine"`
	// WithSCM indicate if we use scm in Pipeline
	// +optional
	WithSCM bool `json:"withSCM"`
	// Agent indicate where the pipeline will running
	// +optional
	Agent *JenkinsAgent `json:"agent,omitempty"`
	// Stages contains all stages of a pipeline script
	Stages []PipelineStage `json:"stages"`
	// Parameters will need before pipeline run
	// +optional
	Parameters []PipelineParameter `json:"parameters,omitempty"`
	// Arguments is arguments for templates
	// +optional
	Arguments []PipelineTemplateArgumentGroup `json:"arguments"`
	// Environments is environment for jenkinsfile
	// +optional
	Environments []PipelineEnvironment `json:"environments,omitempty"`
	// +optional
	Options PipelineOptions `json:"options,omitempty"`
	// +optional
	Triggers PipelineTriggers `json:"triggers,omitempty"`
	// +optional
	Post map[string][]PipelineTemplateTask `json:"post,omitempty"`
	// +optional
	ConstValues ConstValues `json:"values,omitempty"`
}

type PipelineTriggers struct {
	Raw string `json:"raw"`
}

type ConstValues struct {
	// +optional
	Tasks map[string]TaskConstValue `json:"tasks,omitempty"`
}

type TaskConstValue struct {
	Args map[string]string `json:"args,omitempty"`
	// +optional
	Options PipelineOptions `json:"options"`
	// +optional
	Approve PipelineTaskApprove `json:"approve"`
}

// PipelineOptions  specifed options for templatespec
type PipelineOptions struct {
	Timeout int `json:"timeout"`
}

// JenkinsAgent specifed agent for PipelineTemplate
type JenkinsAgent struct {
	// Label is the label for agents
	// +optional
	Label string `json:"label"`
	// Raw is the text plain for Jenkins agent
	// +optional
	Raw string `json:"raw"`
}

// PipelineStage specifed stage for pipeline
type PipelineStage struct {
	// Name is name for a stage
	Name string `json:"name"`
	// Display is display for a stage
	// +optional
	Display I18nName `json:"display,omitempty"`
	// Tasks contains all tasks which will running
	Tasks []PipelineTemplateTask `json:"tasks"`
	// +optional
	Conditions map[string][]string `json:"conditions,omitempty"`
}

// PipelineTemplateTask spcifed task template for PipelineTemplate
type PipelineTemplateTask struct {
	// ID is the id of task
	// +optional
	ID string `json:"id,omitempty"`
	// Name is the name of a task template reference
	Name string `json:"name"`
	// Display is display for a task If display is null it will use Name as Display
	// +optional
	Display I18nName `json:"display"`
	// Agent indicate that where the current task will running
	// +optional
	Agent *JenkinsAgent `json:"agent,omitempty"`
	// Type is type of a task
	// +optional
	Type string `json:"type"`
	// Kind is kind of a task template reference
	Kind string `json:"kind"`
	// Options is some options for a task
	// +optional
	Options *PipelineTaskOption `json:"options,omitempty"`
	// Approve is a option for maual confirm
	// +optional
	Approve *PipelineTaskApprove `json:"approve,omitempty"`
	// Environments contains custom define variables
	// +optional
	Environments []PipelineEnvironment `json:"environments,omitempty"`
	// Relation relation between task and arguments
	// +optional
	Relation []PipelineTaskArgumentAction `json:"relation,omitempty"`
	// +optional
	Conditions map[string][]string `json:"conditions,omitempty"`
}

// PipelineEnvironment specifed environment for Pipeline
type PipelineEnvironment struct {
	// Name is a key for environment map
	Name string `json:"name"`
	// Value is a value for environment map
	Value string `json:"value"`
}

// PipelineTaskApprove specfied approve option for pipeline
type PipelineTaskApprove struct {
	// Message is the message show to users
	Message string `json:"message"`
	// Timeout is timeout for waiting
	// +optional
	Timeout int64 `json:"timeout"`
}

// PipelineTaskOption spcified task option for task template
type PipelineTaskOption struct {
	// Timeout is timeout for a operation
	// +optional
	Timeout int64 `json:"timeout"`
}

// PipelineTemplateArgumentGroup specifed argument group for PipelineTemplate
type PipelineTemplateArgumentGroup struct {
	// DisplayName is used to display
	DisplayName I18nName `json:"displayName"`
	// Items contains all argument for templates
	Items []PipelineTemplateArgumentValue `json:"items"`
}

// I18nName spcified name for Piepline's stage or aruguments
type I18nName struct {
	// Zh is the Chinese name
	Zh string `json:"zh-CN,omitempty"`
	// EN is the English name
	En string `json:"en,omitempty"`
}

// PipelineTemplateArgument specified arugment for PipelineTemplate
type PipelineTemplateArgument struct {
	// Name is the name of a PipelineTemplate
	Name string `json:"name"`
	// Schema is the schema of a template
	Schema PipelineTaskArgumentSchema `json:"schema"`
	// Binding mean bind argument to task
	// +optional
	Binding []string `json:"binding"`
	// Display is used to display
	Display PipelineTaskArgumentDisplay `json:"display"`
	// Required specific argument is required
	// +optional
	Required bool `json:"required"`
	// Default specific argument has default value
	// +optional
	Default string `json:"default"`
	// Validation specific validation for argument
	// +optional
	Validation *PipelineTaskArgumentValidation `json:"validation,omitempty"`
	// Relation relation between arguments
	// +optional
	Relation []PipelineTaskArgumentAction `json:"relation,omitempty"`
}

// PipelineTemplateArgumentValue hold argument and value
type PipelineTemplateArgumentValue struct {
	PipelineTemplateArgument `json:",inline"`
	// Value is the value of a argument
	// +optional
	Value string `json:"value"`
}

// PipelineTaskArgumentSchema specifed arugment schema
type PipelineTaskArgumentSchema struct {
	// Type is the type of a argument
	Type string `json:"type"`
}

// PipelineTaskArgumentDisplay specifed the way of dipslay
type PipelineTaskArgumentDisplay struct {
	// Type is the type of a argument
	Type string `json:"type"`
	// Name contains multi-languages name
	Name I18nName `json:"name"`
	// Related related to other arguments
	// +optional
	Related string `json:"related,omitempty"`
	// Description is used to describe the arugments
	// +optional
	Description I18nName `json:"description,omitempty"`
	// Advanced field has default value
	// +optional
	Advanced bool `json:"advanced,omitempty"`
	// Args is used to add extra data to this argument
	// +optional
	Args map[string]string `json:"args,omitempty"`
}

// PipelineConfigTemplate is instance of template
type PipelineConfigTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// Labels add some marks
	// +optional
	Labels map[string]string `json:"labels"`
	// Spec specification of PipelineConfigTemplate
	Spec PipelineConfigTemplateSpec `json:"spec"`
}

// PipelineConfigTemplateSpec specified for  PipelineConfigTemplate
type PipelineConfigTemplateSpec struct {
	// Engine is a engine for render PipelineTemplate
	// +optional
	Engine PipelineTemplateEngine `json:"engine"`
	// WithSCM indicate if pipeline needs a scm
	// +optional
	WithSCM bool `json:"withSCM"`
	// Agent agent will indicates where the pipeline will running
	// +optional
	Agent *JenkinsAgent `json:"agent"`
	// Stages contains all stages for a pipeline
	Stages []PipelineStageInstance `json:"stages"`
	// Parameters is for execute process usage
	// +optional
	Parameters []PipelineParameter `json:"parameters"`
	// Arguments contains all arguments need by template
	// +optional
	Arguments []PipelineTemplateArgumentGroup `json:"arguments"`
	// Dependencies indicates plugins denpendencies of task
	// +optional
	Dependencies *PipelineDependency `json:"dependencies"`
	// Environments contains env config for jenkinsfile
	// +optional
	Environments []PipelineEnvironment `json:"environments"`
}

// PipelineStageInstance is a instance of PipelineStage
type PipelineStageInstance struct {
	// Name is the name a stage
	Name string `json:"name"`
	// Tasks contains all task include in a stage
	Tasks []PipelineTemplateTaskInstance `json:"tasks"`
}

// PipelineTemplateTaskInstance is a instance of PipelineTemplateTask
type PipelineTemplateTaskInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec PipelineTemplateTaskInstanceSpec `json:"spec"`
}

// PipelineTemplateTaskInstanceSpec specified PipelineTemplateTaskInstance
type PipelineTemplateTaskInstanceSpec struct {
	// Agent specific task running agent
	// +optional
	Agent *JenkinsAgent `json:"agent"`
	// Engine will render task template
	// +optional
	Engine PipelineTemplateTaskEngine `json:"engine"`
	// Type is type of task
	// +optional
	Type string `json:"type"`
	// Body is the body of a pipeline
	// +optional
	Body string `json:"body"`
	// Options is option for task
	// +optional
	Options *PipelineTaskOption `json:"options"`
	// Approve mean task need to approve
	// +optional
	Approve *PipelineTaskApprove `json:"approve"`
	// Environments is variables for pipeline
	// +optional
	Environments []PipelineEnvironment `json:"environments"`
	// Exports mean task will export some variables
	// +optional
	Exports []GlobalParameter `json:"exports"`
	// Arguments contains all arguments include in templates
	// +optional
	Arguments []PipelineTemplateArgument `json:"arguments"`
	// Relation relation between task and arguments
	// +optional
	Relation []PipelineTaskArgumentAction `json:"relation"`
}

type versionable interface {
	GetVersionedName() (string, error)
}

// PipelineTemplateInterface interface of PipelineTemplate
type PipelineTemplateInterface interface {
	runtime.Object
	GetPiplineTempateSpec() *PipelineTemplateSpec
	// GetTypeMeta() *metav1.TypeMeta{ // TypeMeta is not always has kind
	GetKind() string
	GetStatus() *TemplateStatus
	metav1.ObjectMetaAccessor
	metav1.Object
	versionable
	fmt.Stringer
}

// PipelineTaskTemplateInterface interface of PipelineTaskTemplate
type PipelineTaskTemplateInterface interface {
	runtime.Object
	GetPiplineTaskTempateSpec() *PipelineTaskTemplateSpec
	// GetTypeMeta() *metav1.TypeMeta // TypeMeta is not always has kind
	GetKind() string
	GetStatus() *TemplateStatus
	metav1.ObjectMetaAccessor
	metav1.Object
	versionable
	fmt.Stringer
}

var _ PipelineTemplateInterface = &ClusterPipelineTemplate{}
var _ PipelineTemplateInterface = &PipelineTemplate{}
var _ PipelineTaskTemplateInterface = &ClusterPipelineTaskTemplate{}
var _ PipelineTaskTemplateInterface = &PipelineTaskTemplate{}

func (template *ClusterPipelineTemplate) String() string {
	return fmt.Sprintf("%s/%s", template.GetKind(), template.Name)
}

// GetTypeMeta get typemeta
func (template *ClusterPipelineTemplate) GetKind() string {
	return TypeClusterPipelineTemplate
}

// GetStatus get status
func (template *ClusterPipelineTemplate) GetStatus() *TemplateStatus {
	return &template.Status
}

// GetPiplineTempateSpec get PipelineTemplateSpec
func (template *ClusterPipelineTemplate) GetPiplineTempateSpec() *PipelineTemplateSpec {
	return &template.Spec
}

// GetVersionedName get name appended version suffix
func (template *ClusterPipelineTemplate) GetVersionedName() (string, error) {
	version := template.Annotations[AnnotationsTemplateVersion]
	if version == "" {
		return "", fmt.Errorf("ClusterPipelineTemplate '%s/%s' has no version", template.Namespace, template.Name)
	}

	templateName := template.Annotations[AnnotationsTemplateName]
	if templateName == "" {
		templateName = strings.TrimSuffix(template.Name, "."+version)
	}

	return fmt.Sprintf("%s.%s", templateName, version), nil
}

func (template *PipelineTemplate) String() string {
	return fmt.Sprintf("%s/%s/%s", template.GetKind(), template.Namespace, template.Name)
}

// GetKind get kind
func (template *PipelineTemplate) GetKind() string {
	return TypePipelineTemplate
}

// GetStatus get status
func (template *PipelineTemplate) GetStatus() *TemplateStatus {
	return &template.Status
}

// GetPiplineTempateSpec get PipelineTemplateSpec
func (template *PipelineTemplate) GetPiplineTempateSpec() *PipelineTemplateSpec {
	return &template.Spec
}

// GetVersionedName get name appended version suffix
func (template *PipelineTemplate) GetVersionedName() (string, error) {
	version := template.Annotations[AnnotationsTemplateVersion]
	if version == "" {
		return "", fmt.Errorf("PipelineTemplate '%s/%s' has no version", template.Namespace, template.Name)
	}
	templateName := template.Annotations[AnnotationsTemplateName]
	if templateName == "" {
		templateName = strings.TrimSuffix(template.Name, "."+version)
	}
	return fmt.Sprintf("%s.%s", template.Name, version), nil
}

func (template *ClusterPipelineTaskTemplate) String() string {
	return fmt.Sprintf("%s/%s", template.GetKind(), template.Name)
}

// GetKind get kind
func (template *ClusterPipelineTaskTemplate) GetKind() string {
	return TypeClusterPipelineTaskTemplate
}

// GetStatus get status
func (template *ClusterPipelineTaskTemplate) GetStatus() *TemplateStatus {
	return &template.Status
}

// GetPiplineTaskTempateSpec get PipelineTemplateSpec
func (template *ClusterPipelineTaskTemplate) GetPiplineTaskTempateSpec() *PipelineTaskTemplateSpec {
	return &template.Spec
}

// GetVersionedName get name appended version suffix
func (template *ClusterPipelineTaskTemplate) GetVersionedName() (string, error) {
	version := template.Annotations[AnnotationsTemplateVersion]
	if version == "" {
		return "", fmt.Errorf("ClusterPipelineTaskTemplate '%s/%s' has no version", template.Namespace, template.Name)
	}
	templateName := template.Annotations[AnnotationsTemplateName]
	if templateName == "" {
		templateName = strings.TrimSuffix(template.Name, "."+version)
	}
	return fmt.Sprintf("%s.%s", template.Name, version), nil
}

func (template *PipelineTaskTemplate) String() string {
	return fmt.Sprintf("%s/%s/%s", template.GetKind(), template.Namespace, template.Name)
}

// GetKind get kind
func (template *PipelineTaskTemplate) GetKind() string {
	return TypePipelineTaskTemplate
}

// GetStatus get status
func (template *PipelineTaskTemplate) GetStatus() *TemplateStatus {
	return &template.Status
}

// GetPiplineTaskTempateSpec get PipelineTemplateSpec
func (template *PipelineTaskTemplate) GetPiplineTaskTempateSpec() *PipelineTaskTemplateSpec {
	return &template.Spec
}

// GetVersionedName get name appended version suffix
func (template *PipelineTaskTemplate) GetVersionedName() (string, error) {
	version := template.Annotations[AnnotationsTemplateVersion]
	if version == "" {
		return "", fmt.Errorf("PipelineTaskTemplate '%s/%s' has no version", template.Namespace, template.Name)
	}
	templateName := template.Annotations[AnnotationsTemplateName]
	if templateName == "" {
		templateName = strings.TrimSuffix(template.Name, "."+version)
	}
	return fmt.Sprintf("%s.%s", template.Name, version), nil
}

func NewPipelineTemplate(kind string) PipelineTemplateInterface {
	switch kind {
	case TypePipelineTemplate:
		return &PipelineTemplate{}
	default:
		return &ClusterPipelineTemplate{}
	}
}

func NewPipelineTaskTemplate(kind string) PipelineTaskTemplateInterface {
	switch kind {
	case TypePipelineTaskTemplate:
		return &PipelineTaskTemplate{}
	default:
		return &ClusterPipelineTaskTemplate{}
	}
}
