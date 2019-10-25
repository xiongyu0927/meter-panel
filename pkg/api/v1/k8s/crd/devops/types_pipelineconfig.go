package devops

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JenkinsfilePreview jenkinsfile's preview
type JenkinsfilePreview struct {
	metav1.TypeMeta `json:",inline"`

	// Jenkinsfile generated by template or other stuff
	Jenkinsfile string `json:"jenkinsfile"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JenkinsfilePreviewOptions request option for jenkinsfile's preview
type JenkinsfilePreviewOptions struct {
	metav1.TypeMeta `json:",inline"`

	// PipelineConfigSpec specification of PipelineConfig
	PipelineConfigSpec *PipelineConfigSpec `json:"template"`
	// Source git source
	// +optional
	Source *PipelineSource `json:"source"`
	// Values arguments values
	// +optional
	Values map[string]string `json:"values"`
	// +optional
	Environments []PipelineEnvironment `json:"environments"`
}

// endregion

// region Pipeline

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineConfigLogOptions for config log
type PipelineConfigLogOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Start is the start number to fetch the log.
	// +optional
	Start int64 `json:"start"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineConfigScanOptions for scan mutli-branch
type PipelineConfigScanOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Delay for scan multi-branch pipeline
	// +optional
	Delay int `json:"delay"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineConfigScanResult represent the multi-branch scan result
type PipelineConfigScanResult struct {
	metav1.TypeMeta `json:",inline"`

	// Code scan result code, 0 represent sucess
	// +optional
	Code int `json:"code"`
	// Success scan result, true represents sucess
	// +optional
	Success bool `json:"success"`
	// Message scan result details
	// +optional
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineConfigLog log for config log
type PipelineConfigLog PipelineLog

// +genclient
// +genclient:method=Preview,verb=create,subresource=preview,input=JenkinsfilePreviewOptions,result=JenkinsfilePreview
// +genclient:method=Scan,verb=create,subresource=scan,input=PipelineConfigScanOptions,result=PipelineConfigScanResult
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineConfig struct holds a reference to a specific pipeline configuration
// and some user data for access
type PipelineConfig struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Specification of the desired behavior of the PipelineConfig.
	// +optional
	Spec PipelineConfigSpec `json:"spec"`
	// Most recently observed status of the PipelineConfig.
	// Populated by the system.
	// Read-only.
	// +optional
	Status PipelineConfigStatus `json:"status"`
}

var _ Phaser = &PipelineConfig{}

// GetPhase returns its phase as string
func (p *PipelineConfig) GetPhase() string {
	return string(p.Status.Phase)
}

// PipelineConfigSpec defines PipelineConfig's specs
type PipelineConfigSpec struct {
	// JenkinsBinding is the jenkinsBinding of the pipelineConfig.
	// PipelineConfig will be synced between k8s and jenkins defined in this jenkinsBinding.
	// +optional
	JenkinsBinding LocalObjectReference `json:"jenkinsBinding"`
	// RunPolicy is the runPolicy of the pipelineConfig.
	// One of the "Serial" and "Parallel".
	// +optional
	RunPolicy PipelineRunPolicy `json:"runPolicy"`
	// RunLimits is the runLimits of the pipelineConfig.
	// Limit the max number of the pipelines stored.
	// +optional
	RunLimits PipelineRunLimits `json:"runLimits"`
	// Parameters is the parameters of the pipelineConfig.
	// +optional
	Parameters []PipelineParameter `json:"parameters"`
	// Triggers is the triggers of the pipelineConfig.
	// +optional
	Triggers []PipelineTrigger `json:"triggers"`
	// Strategy is the strategy of the pipelineConfig.
	// +optional
	Strategy PipelineStrategy `json:"strategy"`
	// Hooks is the hooks of the pipelineConfig.
	// +optional
	Hooks []PipelineHook `json:"hooks"`
	// Source is the source of the pipelineConfig.
	// +optional
	Source PipelineSource `json:"source"`

	// Template is template info when created pipeline by template
	// +optional
	Template PipelineTemplateWithValue `json:"template"`

	//+optional
	Environments []PipelineEnvironment `json:"environments"`
}

// PipelineTemplateWithValue is Pipeline Template with values
type PipelineTemplateWithValue struct {
	PipelineTemplateSource `json:",inline"`
	// Values arguments values when create pipelineconfig by template
	// +optional
	Values map[string]string `json:"values"`

	// GraphValues arguments values when create pipelineconfig by template
	// +optional
	GraphValues map[string]KeyValueSet `json:"graphValues"`
}

// PipelineTemplateSource is template reference or template definition
type PipelineTemplateSource struct {
	// PipelineTemplateRef is reference of pipeline template
	// +optional
	PipelineTemplateRef PipelineTemplateRef `json:"pipelineTemplateRef"`

	// Mold of creating PipelineTemplate
	PipelineTemplate *PipelineTemplateMold `json:"pipelineTemplate"`
}

// PipelineTemplateMold is mold of pipelinetemplate
// it would be interpret as a template for creating pipelinetemplate
type PipelineTemplateMold struct {
	// +optional
	metav1.ObjectMeta `json:"metadata"`
	Spec              PipelineTemplateSpec `json:"spec"`
}

// PipelineTemplateRef is reference of pipeline template
type PipelineTemplateRef struct {
	// Kind is the kind of PipelineTemplate
	Kind string `json:"kind"`
	// Name is the name of PipelineTemplate
	Name string `json:"name"`
	// Namespace is the namespace of PipelineTemplate
	// +optional
	Namespace string `json:"namespace"`
}

// PipelineRunPolicy pipeline run policy
type PipelineRunPolicy string

const (
	// PipelinePolicySerial Serial run policy
	// used to have sequential pipelines for the same PipelineConfig
	PipelinePolicySerial PipelineRunPolicy = "Serial"
	// PipelinePolicyParallel Parallel run policy
	// used to run multiple pipelines of the same PipelineConfig in parallel
	PipelinePolicyParallel PipelineRunPolicy = "Parallel"
)

// PipelineRunLimits specifies a limited number of stored Pipeline runs
type PipelineRunLimits struct {
	// If set, the number of success pipeline stored will be limited.
	// +optional
	SuccessCount int64 `json:"successCount"`
	// If set, the number of failure pipeline stored will be limited.
	// +optional
	FailureCount int64 `json:"failureCount"`
}

// PipelineParameter specifies a parameter for a pipeline
type PipelineParameter struct {
	// Name is the name of the parameter.
	// +optional
	Name string `json:"name"`
	// Type is the type of the parameter.
	// +optional
	Type PipelineParameterType `json:"type"`
	// Value is the value of the parameter.
	// +optional
	Value string `json:"value"`
	// Description is the description of the parameter.
	// +optional
	Description string `json:"description"`
}

// PipelineParameterType type of parameter for pipeline
type PipelineParameterType string

const (
	// PipelineParameterTypeString parameter type string
	PipelineParameterTypeString PipelineParameterType = "string"
	// PipelineParameterTypeBoolean parameter type boolean
	PipelineParameterTypeBoolean PipelineParameterType = "boolean"
)

// PipelineTrigger specifies a trigger for a pipeline
type PipelineTrigger struct {
	// Type is the type of the pipeline trigger.
	// One of "cron" or "codeChange".
	// +optional
	Type PipelineTriggerType `json:"type"`
	// Cron is one trigger type of pipeline.
	// The pipeline will be triggered in accordance with the cron rule.
	// +optional
	Cron *PipelineTriggerCron `json:"cron,omitempty"`
	// CodeChange is one trigger type of pipeline.
	// The pipeline will be triggered once the code was pushed to the repository.
	// +optional
	CodeChange *PipelineTriggerCodeChange `json:"codeChange,omitempty"`
}

// PipelineTriggerType trigger type for pipelines
type PipelineTriggerType string

const (
	// PipelineTriggerTypeCron cron trigger type
	PipelineTriggerTypeCron PipelineTriggerType = "cron"
	// PipelineTriggerTypeCodeChange code change trigger type
	PipelineTriggerTypeCodeChange PipelineTriggerType = "codeChange"
)

// PipelineTriggerCron cron trigger type
type PipelineTriggerCron struct {
	// Enabled timing trigger pipeline or not.
	// +optional
	Enabled bool `json:"enabled"`
	// Rule is the rule of cron trigger.
	// +optional
	Rule string `json:"rule"`
	// Schedule base on weeks and times
	// +optional
	Schedule *PipelineTriggeSchedule `json:"schedule"`
}

// PipelineTriggeSchedule schedule base on weeks and times
type PipelineTriggeSchedule struct {
	// Weeks mon, tue, wed, thu, fri, sat, sun
	// +optional
	Weeks []Week `json:"weeks"`
	// Times like 1:00, 12:00, 16:00
	// +optional
	Times []string `json:"times"`
}

// Week represent week
type Week string

const (
	// WeekMonday Monday
	WeekMonday Week = "mon"
	// WeekTuesday Tuesday
	WeekTuesday Week = "tue"
	// WeekWednesday Wednesday
	WeekWednesday Week = "wed"
	// WeekThursday Thursday
	WeekThursday Week = "thu"
	// WeekFriday Friday
	WeekFriday Week = "fri"
	// WeekSaturday Saturday
	WeekSaturday Week = "sat"
	// WeekSunday Sunday
	WeekSunday Week = "sun"
)

// IsValid check whether Week is valid
func (w Week) IsValid() bool {
	switch w {
	case WeekMonday, "1", WeekTuesday, "2", WeekWednesday, "3",
		WeekThursday, "4", WeekFriday, "5", WeekSaturday, "6",
		WeekSunday, "7":
		return true
	}
	return false
}

// PipelineTriggerCodeChange code change trigger type
type PipelineTriggerCodeChange struct {
	// Enabled trigger pipeline by code change or not.
	// +optional
	Enabled bool `json:"enabled"`
	// PeriodicCheck specifies how often check code changes.
	// +optional
	PeriodicCheck string `json:"periodicCheck"`
}

// PipelineStrategy pipeline execution strategy
type PipelineStrategy struct {
	// Template is a template for pipeline
	// +optional
	Template *PipelineConfigTemplate `json:"template"`
	// Jenkins specific jenkinsfile path
	// +optional
	Jenkins PipelineStrategyJenkins `json:"jenkins"`
}

// PipelineStrategyJenkins jenkins execution strategy
type PipelineStrategyJenkins struct {
	// Jenkinsfile approves the jenkinsfile script will be executed.
	// +optional
	Jenkinsfile string `json:"jenkinsfile,omitempty"`
	// JenkinsfilePath is the jenkinsfile path in the code repository.
	// +optional
	JenkinsfilePath string `json:"jenkinsfilePath,omitempty"`
	// MultiBranch hold multi-branch pipeline configuration
	// +optional
	MultiBranch MultiBranchPipeline `json:"multiBranch"`
}

// MultiBranchPipeline represent multi-branch pipeline
type MultiBranchPipeline struct {
	// Orphaned orphan strategy
	// +optional
	Orphaned MultiBranchOrphan `json:"orphaned"`
	// Behaviours discover strategy for multi-branch pipeline
	// +optional
	Behaviours MultiBranchBehaviours `json:"behaviours"`
}

// MultiBranchOrphan orphan strategy for multi-branch pipeline
type MultiBranchOrphan struct {
	// Days max days for keeping stale pipeline
	DaysAfterClosed int `json:"days"`
	// Max number of max to keep stale pipeline
	MaxNumberToKeep int `json:"max"`
}

// MultiBranchBehaviours discover strategy for multi-branch pipeline
type MultiBranchBehaviours struct {
	// FilterExpression expression for filter the branches or PRs or tags
	// +optional
	FilterExpression string `json:"filterExpression"`
	// DiscoverTags whether discovering tags
	// +optional
	DiscoverTags bool `json:"discoverTags"`
	// DiscoverBranches indicate how to discover branches
	// +optional
	DiscoverBranches string `json:"discoverBranches"`
	// DiscoverPRFromOrigin indicate how to discover PRs from the origin
	// +optional
	DiscoverPRFromOrigin string `json:"discoverPRFromOrigin"`
	// DiscoverPRFromForks indicate how to discover PRs from the forks
	// +optional
	DiscoverPRFromForks string `json:"discoverPRFromForks"`
	// ForksTrust indicate how to trust the Jenkinsfile from forks
	// +optional
	ForksTrust string `json:"forksTrust"`
}

// PipelineHook pipeline hook definition
type PipelineHook struct {
	// Type is the type of the pipeline hook.
	// Now just supports "httpRequest".
	// +optional
	Type PipelineHookType `json:"type"`
	// Events is a list of events of the pipeline changed.
	// +optional
	Events []PipelineEvent `json:"events"`
	// HTTPRequest is the httpRequest of the pipeline hook.
	// +optional
	HTTPRequest *PipelineHookHTTPRequest `json:"httpRequest,omitempty"`
}

// PipelineHookType pipeline hook types
type PipelineHookType string

const (
	// PipelineHookTypeHTTPRequest httpRequest type for sending requests
	PipelineHookTypeHTTPRequest PipelineHookType = "httpRequest"
)

// PipelineEvent pipeline event types
type PipelineEvent string

const (
	// PipelineEventConfigCreated event for creating PipelineConfig
	PipelineEventConfigCreated PipelineEvent = "PipelineConfigCreated"
	// PipelineEventConfigUpdated event for updating PipelineConfig
	PipelineEventConfigUpdated PipelineEvent = "PipelineConfigUpdated"
	// PipelineEventConfigDeleted event for deleting PipelineConfig
	PipelineEventConfigDeleted PipelineEvent = "PipelineConfigDeleted"
	// PipelineEventPipelineStarted event for starting PipelineConfig
	PipelineEventPipelineStarted PipelineEvent = "PipelineStarted"
	// PipelineEventPipelineStopped event for stoping PipelineConfig
	PipelineEventPipelineStopped PipelineEvent = "PipelineStopped"
)

// PipelineHookHTTPRequest HTTP request hook type
type PipelineHookHTTPRequest struct {
	// URI is the uri of the http request.
	// +optional
	URI string `json:"uri"`
	// Method is the method of the http request.
	// +optional
	Method string `json:"method"`
	// Headers is the header of the http request.
	// +optional
	Headers map[string]string `json:"headers"`
}

type PipelineSourceType string

const (
	PipelineSourceTypeGit PipelineSourceType = "GIT"
	PipelineSourceTypeSvn PipelineSourceType = "SVN"
)

// PipelineSource source code specification for Pipeline
type PipelineSource struct {
	// CodeRepository contains git url and user info
	// +optional
	CodeRepository *CodeRepositoryRef `json:"codeRepository,omitempty"`
	// Git is the git code repository settings of the pipeline source.
	// +optional
	Git *PipelineSourceGit `json:"git,omitempty"`
	// Svn is the svn code repository settings of the pipeline source.
	// +optional
	Svn *PipelineSourceSvn `json:"svn,omitempty"`
	// Source Type is the type of Pipeline Source
	SourceType PipelineSourceType `json:"sourceType"`
	// Secret is the secret to access the code repository.
	// +optional
	Secret *SecretKeySetRef `json:"secret,omitempty"`
}

// PipelineSourceGit generic git implementation for PipelineSource
type PipelineSourceGit struct {
	// URI is the uri of the git code repository.
	// +optional
	URI string `json:"uri"`
	// URI is the branch of the git code repository.
	// +optional
	Ref string `json:"ref"`
}

type PipelineSourceSvn struct {
	// URI is the uri of the git code repository.
	// +optional
	URI string `json:"uri"`
}

// CodeRepositoryRef is reference of CodeRepostiory
type CodeRepositoryRef struct {
	LocalObjectReference `json:",inline"`
	// Ref is branch of git repo
	Ref string `json:"ref"`
}

// PipelineConfigStatus represents PipelineConfig's status
// PipelineConfigStatus defines PipelineConfig's status
type PipelineConfigStatus struct {
	// Current condition of the pipelineConfig.
	// +optional
	Phase PipelineConfigPhase `json:"phase"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
	// Last time we update the pipelineConfig.
	// +optional
	LastUpdate *metav1.Time `json:"lastUpdated"`
	// A list of condition referenced to the pipelineConfig.
	// +optional
	Conditions []Condition `json:"conditions,omitempty"`
}

// PipelineConfigPhase a phase for PipelineConfigStatus
type PipelineConfigPhase string

const (
	// PipelineConfigPhaseCreating creating phase of PipelineConfig
	PipelineConfigPhaseCreating PipelineConfigPhase = StatusCreating
	// PipelineConfigPhaseSyncing syncing phase of PipelineConfig
	PipelineConfigPhaseSyncing PipelineConfigPhase = StatusSyncing
	// PipelineConfigPhaseReady ready phase of PipelineConfig
	// this means that the PipelineConfig was already synced
	// to the destinatination service
	PipelineConfigPhaseReady PipelineConfigPhase = StatusReady
	// PipelineConfigPhaseError error phase for the PipelineConfig
	// this means that some error happened during the syncing phase
	// and added the related errors to the Reason and Message
	PipelineConfigPhaseError PipelineConfigPhase = StatusError
	// PipelineConfigPhaseDisabled error phase for the PipelineConfig
	// this means the pipelineconfig is disable for the coderepository deleted
	PipelineConfigPhaseDisabled PipelineConfigPhase = StatusDisabled
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineConfigList is a list of PipelineConfig objects.
type PipelineConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata"`

	// Items is a list of PipelineConfig objects.
	Items []PipelineConfig `json:"items"`
}
