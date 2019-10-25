package devops

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:method=Input,verb=create,subresource=input,input=PipelineInputOptions,result=PipelineInputResponse
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Pipeline struct holds a reference to a specific pipeline run
type Pipeline struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Specification of the desired behavior of the Pipeline.
	// +optional
	Spec PipelineSpec `json:"spec"`
	// Most recently observed status of the Pipeline.
	// Populated by the system.
	// Read-only.
	// +optional
	Status PipelineStatus `json:"status"`
}

// GetObjectMeta object meta
func (p Pipeline) GetObjectMeta() metav1.ObjectMeta {
	return p.ObjectMeta
}

// PipelineSpec specifications for a PipelineConfig
type PipelineSpec struct {
	// JenkinsBinding is the jenkinsBinding of the pipeline.
	// +optional
	JenkinsBinding LocalObjectReference `json:"jenkinsBinding"`
	// PipelineConfig is the pipelineConfig of the pipeline.
	// +optional
	PipelineConfig LocalObjectReference `json:"pipelineConfig"`
	// Cause is the cause of the pipeline.
	// +optional
	Cause PipelineCause `json:"cause"`
	// RunPolicy is the runPolicy of the pipeline.
	// +optional
	RunPolicy PipelineRunPolicy `json:"runPolicy"`
	// Parameters is the parameters of the pipeline.
	// +optional
	Parameters []PipelineParameter `json:"parameters"`
	// Triggers is the triggers of the pipeline.
	// +optional
	Triggers []PipelineTrigger `json:"triggers"`
	// Strategy is the strategy of the pipeline.
	// +optional
	Strategy PipelineStrategy `json:"strategy"`
	// Hooks is the hooks of the pipeline.
	// +optional
	Hooks []PipelineHook `json:"hooks"`
	// Source is the source of the pipeline.
	// +optional
	Source PipelineSource `json:"source"`
}

// PipelineCause describe the cause for a pipeline trigger
type PipelineCause struct {
	// Type is the type of the pipeline pipelineCause.
	// One of "manual"、"cron"、"codeChange".
	// +optional
	Type PipelineCauseType `json:"type"`
	// Human-readable message indicating details about a pipeline cause.
	// +optional
	Message string `json:"message"`
}

// PipelineCauseType pipeline run start cause
type PipelineCauseType string

const (
	// PipelineCauseTypeManual manual execution by user
	PipelineCauseTypeManual PipelineCauseType = "manual"
	// PipelineCauseTypeCron cron timer execution
	PipelineCauseTypeCron PipelineCauseType = "cron"
	// PipelineCauseTypeCodeChange code change execution
	PipelineCauseTypeCodeChange PipelineCauseType = "codeChange"
)

// PipelineStatus pipeline status
type PipelineStatus struct {
	// Current condition of the pipeline.
	// +optional
	Phase PipelinePhase `json:"phase"`
	// StartedAt is the start time of the pipeline.
	// +optional
	StartedAt *metav1.Time `json:"startedAt"`
	// FinishedAt is finish time of the pipeline.
	// +optional
	FinishedAt *metav1.Time `json:"finishedAt"`
	// UpdatedAt is the update time of the pipeline.
	// +optional
	UpdatedAt *metav1.Time `json:"updatedAt"`
	// Jenkins is the status of the jenkins this pipeline used.
	// +optional
	Jenkins *PipelineStatusJenkins `json:"jenkins,omitempty"`
	// Aborted is aborted status of the pipeline trigger.
	// +optional
	Aborted bool `json:"aborted"`
}

// PipelinePhase a phase for PipelineStatus
type PipelinePhase string

// IsValid check whether the pipeline is valid or not.
func (phase PipelinePhase) IsValid() bool {
	switch phase {
	case PipelinePhasePending, PipelinePhaseQueued, PipelinePhaseRunning,
		PipelinePhaseComplete, PipelinePhaseFailed,
		PipelinePhaseError, PipelinePhaseCancelled, PipelinePhaseAborted:
		return true
	}
	return false
}

// IsFinalPhase check whether the pipeline is finished.
func (phase PipelinePhase) IsFinalPhase() bool {
	switch phase {
	case PipelinePhaseComplete, PipelinePhaseFailed, PipelinePhaseError,
		PipelinePhaseCancelled, PipelinePhaseAborted:
		return true
	}
	return false
}

const (
	// PipelinePhasePending created but not yet sinced
	PipelinePhasePending PipelinePhase = "Pending"
	// PipelinePhaseQueued entered in the jenkins queue
	PipelinePhaseQueued PipelinePhase = "Queued"
	// PipelinePhaseRunning started execution
	PipelinePhaseRunning PipelinePhase = "Running"
	// PipelinePhaseComplete finished execution
	PipelinePhaseComplete PipelinePhase = "Complete"
	// PipelinePhaseFailed finished execution but failed
	PipelinePhaseFailed PipelinePhase = "Failed"
	// PipelinePhaseError finished execution but failed
	PipelinePhaseError PipelinePhase = "Error"
	// PipelinePhaseCancelled paused execution
	PipelinePhaseCancelled PipelinePhase = "Cancelled"
	// PipelinePhaseAborted when user aborts a pipeline while in queue
	PipelinePhaseAborted PipelinePhase = "Aborted"
)

// PipelineStatusJenkins used to store jenkins related information
type PipelineStatusJenkins struct {
	// Result is the result of the jenkins.
	// +optional
	Result string `json:"result"`
	// Status is the status of the jenkins.
	// +optional
	Status string `json:"status"`
	// Build is the build of the jenkins.
	// +optional
	Build string `json:"build"`
	// Stages is the stages of the jenkins.
	// +optional
	Stages string `json:"stages"`
	// StartStageID is the startStageID of the jenkins.
	// +optional
	StartStageID string `json:"startStageID"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineList is a list of Pipeline objects.
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of Pipeline objects.
	Items []Pipeline `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineLog used to retrieve logs from a pipeline
type PipelineLog struct {
	metav1.TypeMeta `json:",inline"`

	// True means has more log behind。
	// +optional
	HasMore bool `json:"more"`
	// NextStart is next start number to fetch new log.
	// +optional
	NextStart *int64 `json:"nextStart,omitempty"`
	// Text is the context of the log.
	// +optional
	Text string `json:"text"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineLogOptions used to fetch logs from a pipeline
type PipelineLogOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Start is the start number to fetch the log.
	// +optional
	Start int64 `json:"start"`

	// Stage if given will limit the log to a specific stage
	// +optional
	Stage int64 `json:"stage"`

	// Step if given will limit the log to a specific step
	// +optional
	Step int64 `json:"step"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineTestReport represents the test report from Pipeline
type PipelineTestReport struct {
	metav1.TypeMeta `json:",inline"`

	// Summary is the summary report
	// +optional
	Summary *PipelineTestReportSummary

	// Items for the test report
	// +optional
	Items []PipelineTestReportItem `json:"items"`
}

// PipelineTestReportSummary test report summary
type PipelineTestReportSummary struct {
	// existing failed
	ExistingFailed int64
	// failed
	Failed int64
	// fixed
	Fixed int64
	// passed
	Passed int64
	// regressions
	Regressions int64
	// skipped
	Skipped int64
	// total
	Total int64
}

// PipelineTestReportItem test report item
type PipelineTestReportItem struct {
	Age int `json:"age"`
	// Duration the time of test
	Duration float32 `json:"duration"`
	// ErrorDetails error details of test
	ErrorDetails string `json:"errorDetails"`
	// ErrorStackTrace if the status is erro then error stack trace of test
	ErrorStackTrace string `json:"errorStackTrace"`
	// HasStdLog indicate whether has standard log outpupt
	HasStdLog bool `json:"hasStdLog"`
	// ID id for the test report item
	ID string `json:"id"`
	// Name is the name of test report item
	Name  string `json:"name"`
	State string `json:"state"`
	// Status indicates the status of report item
	Status string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineTestReportOptions for getting the test reports from Jenkins
type PipelineTestReportOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Start indicates the offset of reports
	Start int64 `json:"start"`
	// Limit indicates the num of report items
	Limit int64 `json:"limit"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineTask retrieve steps or stages from a pipeline
type PipelineTask struct {
	metav1.TypeMeta `json:",inline"`

	// Tasks steps/stages for a Pipeline
	Tasks []PipelineBlueOceanTask `json:"tasks"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineTaskOptions options for requesting stage/steps from jenkins blue ocean
type PipelineTaskOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Stage indicates the stage id to fetch the step list
	// if not provided will fetch the stage list
	// +optional
	Stage int64 `json:"stage"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineInputOptions options for pipeline input request
type PipelineInputOptions struct {
	metav1.TypeMeta

	// Stage if given will limit the log to a specific stage
	Stage int64 `json:"stage"`

	// Step if given will limit the log to a specific step
	Step int64 `json:"step"`

	// Approve whether approve this
	Approve bool `json:"approve"`
	// InputID is the id for input dsl step from Jenkinsfile
	InputID string `json:"inputID"`
	// PlatformApprover for who approve or reject this
	// +optional
	PlatformApprover string `json:"platformApprover"`
	// Parameters is the parameters of the pipeline input request
	// +optional
	Parameters []PipelineParameter `json:"parameters"`
}

// PipelineInputRequest represent the input request model
type PipelineInputRequest struct {
	BaseURI   string
	BuildID   string
	ID        string
	Message   string
	Status    string
	Submitter string
}

type PipelineInputRequestList []PipelineInputRequest

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PipelineInputResponse represent the response of input request
type PipelineInputResponse struct {
	metav1.TypeMeta

	// Message is a description about response
	Message string `json:"message"`
	// StatusCode represent the response status code
	StatusCode int `json:"statusCode"`
}

// PipelineBlueOceanTask a task from BlueOcean API
type PipelineBlueOceanTask struct {
	// extends PipelineBlueOceanRef
	PipelineBlueOceanRef
	// DisplayDescription description for step/stage
	DisplayDescription string `json:"displayDescription"`
	// DisplayName is a display name for step/stage
	// +optional
	DisplayName string `json:"displayName"`
	// Duration in milliseconds
	// +optional
	DurationInMillis int64 `json:"durationInMillis"`
	// Input describes a input for Jenkins step
	// +optional
	Input *PipelineBlueOceanInput `json:"input"`

	// Result describes a result for a stage/step in Jenkins
	Result string `json:"result"`
	// Stage describe the current state of the stage/step in Jenkins
	State string `json:"state"`
	// StartTime the starting time for the stage/step
	// +optional
	StartTime string `json:"startTime,omitempty"`

	// Edges edges for a specific stage
	// +optional
	Edges []PipelineBlueOceanRef `json:"edges,omitempty"`
	// Actions
	// +optional
	Actions []PipelineBlueOceanRef `json:"actions,omitempty"`
}

// PipelineBlueOceanRef reference of a class/resource
type PipelineBlueOceanRef struct {
	// Href reference url for resource
	// +optional
	Href string `json:"href,omitempty"`
	// ID unique identifier for step/stage
	// +optional
	ID string `json:"id,omitempty"`
	// Type describes the resource type
	// +optional
	Type string `json:"type,omitempty"`
	// URLName describes a url name for the resource
	// +optional
	URLName string `json:"urlName,omitempty"`

	// Description description for reference
	// +optional
	Description string `json:"description,omitempty"`
	// Name name for reference
	// +optional
	Name string `json:"name,omitempty"`

	// Value for reference
	// +optional
	Value string `json:"value,omitempty"`
}

// ComposeValue represent a compose value
// type ComposeValue struct {
// 	Value string `json:"value,omitempty"`
// }

// UnmarshalJSON convert multi-type into string
// func (f ComposeValue) UnmarshalJSON(data []byte) error {
// 	originStr := strings.Trim(string(data), `"`)
// 	switch str := strings.ToLower(originStr); str {
// 	case "true":
// 		f.Value = "true"
// 	case "false":
// 		f.Value = "false"
// 	default:
// 		f.Value = originStr
// 	}
// 	return nil
// }

// MarshalJSON marshal the compose value
// func (f ComposeValue) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(f.Value)
// }

// func (f ComposeValue) String() string {
// 	return f.Value
// }

// PipelineBlueOceanInput describes a Jenkins input for a step
type PipelineBlueOceanInput struct {
	// extends PipelineBlueOceanRef
	PipelineBlueOceanRef
	// Message describes the message for the input
	Message string `json:"message"`
	// OK describes which option is used for successful submit
	OK string `json:"ok"`

	// Parameters parameters for input
	// +optional
	Parameters []PipelineBlueOceanParameter `json:"parameters,omitempty"`
	// Submitter list of usernames or user ids that can approve
	// +optional
	Submitter string `json:"submitter"`
}

// PipelineBlueOceanParameter one step parameter for Jenkins step
type PipelineBlueOceanParameter struct {
	PipelineBlueOceanRef
	// DefaultParameterValue type and default value for parameter
	// +optional
	DefaultParameterValue PipelineBlueOceanRef `json:"defaultParameterValue"`
}

// PipelineConfigData defines the old and new config info.
type PipelineConfigData struct {
	// Old is the old pipelineConfig info.
	// +optional
	Old *PipelineConfig `json:"old"`
	// New is the old pipelineConfig info.
	// +optional
	New *PipelineConfig `json:"new"`
}

// PipelineConfigPayload defines pipelineConfig payload in event.
type PipelineConfigPayload struct {
	// Event is the event of the payload.
	// +optional
	Event PipelineEvent `json:"event"`
	// Data is the data of the payload.
	// +optional
	Data PipelineConfigData `json:"data"`
}

type PipelineData struct {
	// Old is the old pipeline info.
	// +optional
	Old *Pipeline `json:"old"`
	// New is the new pipeline info.
	// +optional
	New *Pipeline `json:"new"`
}

// PipelinePayload defines pipeline payload in event.
type PipelinePayload struct {
	// Event is the event of the payload.
	// +optional
	Event PipelineEvent `json:"event"`
	// Data is the data of the payload.
	// +optional
	Data PipelineData `json:"data"`
}

// HasPipelineEvent check whether include events
func (ph *PipelineHook) HasPipelineEvent(event PipelineEvent) bool {
	if ph.Events == nil || len(ph.Events) == 0 {
		return false
	}

	for _, e := range ph.Events {
		if e == event {
			return true
		}
	}
	return false
}

// GetLastNumber get the last trigger number
func (p *PipelineConfig) GetLastNumber() (lastNumber string) {
	annotations := p.GetAnnotations()
	if annotations == nil || len(annotations) == 0 {
		return
	}

	lastNumber, _ = annotations[AnnotationsKeyPipelineLastNumber]
	return
}

// endregion
