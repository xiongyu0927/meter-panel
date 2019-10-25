package devops

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

const (
	// ProductName product name
	ProductName = "Alauda DevOps"
	// ProjectNamespaceType is a relationship type between both resources
	// currently only "owned" is supported

	// ProjectNamespaceTypeOwned ownership type of relationship
	ProjectNamespaceTypeOwned = "owned"
	// ProjectNamespaceStatusUnknown status unknown
	ProjectNamespaceStatusUnknown = "Unknown"
	// ProjectNamespaceStatusReady status ready
	ProjectNamespaceStatusReady = "Ready"
	// ProjectNamespaceStatusError status error
	ProjectNamespaceStatusError = "Error"

	// StatusCreating creating
	StatusCreating = "Creating"
	// StatusReady can be reached
	StatusReady = "Ready"
	// StatusError service cannot be reached
	StatusError = "Error"
	// StatusSyncing Syncing
	StatusSyncing = "Syncing"
	// StatusDisabled Disabled
	StatusDisabled = "Disabled"
	//StatusWaitingToDelete WaitingToDelete
	StatusWaitingToDelete = "WaitingToDelete"
	//StatusListRepoError ListRepoError
	StatusListRepoError = "ListRepoError"
	//StatusListTagError ListTagError
	StatusListTagError = "ListTagError"
	// StatusNeedsAuthorization NeedsAuthorization
	StatusNeedsAuthorization = "NeedsAuthorization"

	// JenkinsBindingStatus has a slice of Conditions (like pods)
	// that are used to store different parts of required data
	// like Jenkins and Secret

	ProjectManagementBindingStatusTypeProjectmanage = "ProjectManagement"

	ProjectmanagementStatusTypeUserCount                  = "UserCount"
	ProjectManagementBindingStatusConditionStatusNotValid = "NotValid"

	ProjectManagementBindingStatusConditionStatusNotFound = "NotFound"

	ProjectManagementBindingStatusConditionStatusReady = "Ready"

	DocumentManagementBindingStatusTypeDocumentManagement = "DocumentManagement"

	DocumentManagementBindingStatusConditionStatusNotValid = "NotValid"

	DocumentManagementBindingStatusConditionStatusNotFound = "NotFound"

	DocumentManagementBindingStatusConditionStatusReady = "Ready"

	// JenkinsBindingStatusTypeJenkins condition type for Jenkinsbinding
	JenkinsBindingStatusTypeJenkins = "Jenkins"
	// JenkinsBindingStatusTypeSecret condition type for Secret
	JenkinsBindingStatusTypeSecret = "Secret"
	// JenkinsBindingStatusTypeRepository condition type for Secret
	JenkinsBindingStatusTypeRepository = "CodeRepository"
	// JenkinsBindingStatusTypeImageRepository condition type for ImageRepository
	JenkinsBindingStatusTypeImageRepository = "ImageRepository"
	// JenkinsBindingStatusTypeCodeQualityProject condition type for CodeQualityProject
	JenkinsBindingStatusTypeCodeQualityProject = "CodeQualityProject"

	// For each condition type we have a list of valid reasons:

	// JenkinsBindingStatusConditionStatusNotFound when a condition type is NotFound
	JenkinsBindingStatusConditionStatusNotFound = "NotFound"
	// JenkinsBindingStatusConditionStatusNotValid when a condition type is not filled
	JenkinsBindingStatusConditionStatusNotValid = "NotValid"
	// JenkinsBindingStatusConditionStatusReady when a condition type is Found
	JenkinsBindingStatusConditionStatusReady = "Ready"

	// region Annotation: Used as key in Annotations

	// AnnotationsKeyDisplayName displayName key for annotations
	AnnotationsKeyDisplayName = "alauda.io/displayName"
	// AnnotationsKeyDisplayNameEn english displayName key for annotations
	AnnotationsKeyDisplayNameEn = "alauda.io/displayNameEn"
	// AnnotationsKeyProduct product name key for annotations
	AnnotationsKeyProduct = "alauda.io/product"
	// AnnotationsKeyProductVersion product version key for annotations
	AnnotationsKeyProductVersion = "alauda.io/productVersion"
	// AnnotationsKeyProject project key for annotations
	AnnotationsKeyProject = "alauda.io/project"
	// AnnotationsKeySubProject subProject key for annotations
	AnnotationsKeySubProject = "alauda.io/subProject"
	// AnnotationsKeyPipelineLastNumber last number used to generate pipeline store in PipelineConfig
	AnnotationsKeyPipelineLastNumber = "alauda.io/pipeline.last.number"
	// AnnotationsKeyPipelineNumber key for specific pipeline Number
	AnnotationsKeyPipelineNumber = "alauda.io/pipeline.number"
	// AnnotationsKeyPipelineConfig key for specific pipeline config
	AnnotationsKeyPipelineConfig = "alauda.io/pipelineConfig.name"
	// AnnotationsKeyPipelineConfigScanLog key for multi-branch log url
	AnnotationsKeyPipelineConfigScanLog = "alauda.io/multi-branch-scan-log"
	// AnnotationsJenkinsBuildURI annotations key that contains the build uri
	AnnotationsJenkinsBuildURI = "alauda.io/jenkins-build-uri"
	// AnnotationsJenkinsStagesLogURI annotations key that contains the stages log uri
	AnnotationsJenkinsStagesLogURI = "alauda.io/jenkins-stages-log"
	// AnnotationsJenkinsStagesURI annotations key that contains the stages uri
	AnnotationsJenkinsStagesURI = "alauda.io/jenkins-stages"
	// AnnotationsJenkinsStepsLogURI annotations key that contains the steps log uri
	AnnotationsJenkinsStepsLogURI = "alauda.io/jenkins-steps-log"
	// AnnotationsJenkinsStepsURI annotatios key that contains steps uri
	AnnotationsJenkinsStepsURI = "alauda.io/jenkins-steps"
	// AnnotationsJenkinsViewLogURI annotations key that contains view log uri
	AnnotationsJenkinsViewLogURI = "alauda.io/jenkins-view-log"
	// AnnotationsJenkinsProgressiveLogURI annotations key that contains progressive log uri
	AnnotationsJenkinsProgressiveLogURI = "alauda.io/jenkins-progressive-log"
	//AnnotationsJenkinsMultiBranchName annotations key that contains branch name of pipeline if pipeline is a multi-branch job
	AnnotationsJenkinsMultiBranchName = "alauda.io/multiBranchName"

	// AnnotationsSecretType annotations key for specific secret type
	AnnotationsSecretType = "alauda.io/secretType"
	// AnnotationsCreateAppUrl annotations key for specific create app url
	AnnotationsCreateAppUrl = "alauda.io/createAppUrl"
	// AnnotationsToolHttpHost annotations key for http host url of tool
	AnnotationsToolHttpHost = "alauda.io/toolHttpHost"
	// AnnotationsToolAccessURL annotations key for http access url of tool
	AnnotationsToolAccessURL = "alauda.io/toolAccessUrl"
	// AnnotationsToolSubscription allow subscription for tool in project
	AnnotationsToolSubscription = "alauda.io/subscription"
	// AnnotationsToolItemType Type for specific tool, e.g Github, Bitbucket, Harbor etc
	AnnotationsToolItemType = "alauda.io/toolItemType"
	// AnnotationsToolItemKind Kind of the tool, uses a an API kind
	AnnotationsToolItemKind = "alauda.io/toolItemKind"
	AnnotationsToolName     = "alauda.io/toolName"
	// AnnotationsToolItemPublic stablishes if the integrated tool is a public instance, e.g Github
	AnnotationsToolItemPublic  = "alauda.io/toolItemPublic"
	AnnotationsToolItemProject = "alauda.io/toolItemProject"
	// AnnotationsToolType used as a catagory for each tool kind, example: continuousDelivery, codeRepository etc.
	AnnotationsToolType = "alauda.io/toolType"
	// AnnotationsProjectDataAvatarURL the avatar of the "project" in tool
	AnnotationsProjectDataAvatarURL = "avatarURL"
	// AnnotationsProjectDataAccessPath the access path of the "project" in tool
	AnnotationsProjectDataAccessPath = "accessPath"
	// AnnotationsProjectDataDescription the description of the "project" in tool
	AnnotationsProjectDataDescription = "description"
	// AnnotationsProjectDataType the type of the "project" in tool, tg. Org, Team ,Group, SubGroup, enum by tool
	AnnotationsProjectDataType    = "type"
	AnnotationsSecretProductACE   = "ACE"
	AnnotationsGeneratedBy        = "alauda.io/generatedBy"
	AnnotationsGeneratorNamespace = "alauda.io/generatorNamespace"
	AnnotationsGeneratorName      = "alauda.io/generatorName"
	AnnotationsUsername           = "alauda.io/username"

	AnnotationsSonarQubeProjectLink = "alauda.io/sonarqubeProjectLink"
	AnnotationsTemplateVersion      = "alauda.io/version"
	// AnnotationsTemplateName the original name of template
	AnnotationsTemplateName          = "templateName"
	AnnotationsTemplateLatest        = "alauda.io/latest"
	AnnotationsTemplateLatestVersion = "alauda.io/templateLatestVersion"

	AnnotationsHeartbeatTouchTime = "alauda.io/heartbeatTouchTime"

	// AnnotationsImageRegistryEndpoint for imageRegistryEndpoint
	AnnotationsImageRegistryEndpoint = "imageRegistryEndpoint"

	AnnotationsTemplateMoldHash = "alauda.io/templateMoldHash"

	// endregion

	// region Type: Pascal nomenclature
	TypeProject                     = "Project"
	TypeNamespace                   = "Namespace"
	TypeJenkins                     = "Jenkins"
	TypeJenkinsBinding              = "JenkinsBinding"
	TypePipelineConfig              = "PipelineConfig"
	TypePipeline                    = "Pipeline"
	TypePipelineTemplate            = "PipelineTemplate"
	TypePipelineTaskTemplate        = "PipelineTaskTemplate"
	TypePipelineTemplateSync        = "PipelineTemplateSync"
	TypeCodeRepoService             = "CodeRepoService"
	TypeCodeRepoBinding             = "CodeRepoBinding"
	TypeCodeRepository              = "CodeRepository"
	TypeImageRegistry               = "ImageRegistry"
	TypeImageRegistryBinding        = "ImageRegistryBinding"
	TypeImageRepository             = "ImageRepository"
	TypeProjectManagement           = "ProjectManagement"
	TypeProjectManagementBinding    = "ProjectManagementBinding"
	TypeDocumentManagement          = "DocumentManagement"
	TypeDocumentManagementBinding   = "DocumentManagementBinding"
	TypeTestTool                    = "TestTool"
	TypeTestToolBinding             = "TestToolBinding"
	TypeToolBindingReplica          = "ToolBindingReplica"
	TypeSecret                      = "Secret"
	TypeConfigMap                   = "ConfigMap"
	TypeCodeQualityTool             = "CodeQualityTool"
	TypeCodeQualityBinding          = "CodeQualityBinding"
	TypeCodeQualityProject          = "CodeQualityProject"
	TypeClusterPipelineTemplate     = "ClusterPipelineTemplate"
	TypeClusterPipelineTaskTemplate = "ClusterPipelineTaskTemplate"
	TypeToolBindingMap              = "ToolBindingMap"

	TypeToolBinding = "ToolBinding"

	TypeArtifactRegistry        = "ArtifactRegistry"
	TypeArtifactRegistryBinding = "ArtifactRegistryBinding"
	TypeArtifactRegistryManager = "ArtifactRegistryManager"

	// endregion

	// region Kind: all words were lower case
	ResourceKindProject                   = "project"
	ResourceKindJenkins                   = "jenkins"
	ResourceKindJenkinsBinding            = "jenkinsbinding"
	ResourceKindPipelineConfig            = "pipelineconfig"
	ResourceKindPipeline                  = "pipeline"
	ResourceKindCodeRepoService           = "codereposervice"
	ResourceKindCodeRepoBinding           = "coderepobinding"
	ResourceKindCodeRepository            = "coderepository"
	ResourceKindImageRegistry             = "imageregistry"
	ResourceKindImageRegistryBinding      = "imageregistrybinding"
	ResourceKindImageRepository           = "imagerepository"
	ResourceKindProjectManagement         = "projectmanagement"
	ResourceKindProjectManagementBinding  = "projectmanagementbinding"
	ResourceKindDocumentManagement        = "documentmanagement"
	ResourceKindDocumentManagementBinding = "documentmanagementbinding"
	ResourceKindTestTool                  = "testtool"
	ResourceKindTestToolBinding           = "testtoolbinding"
	ResourceKindCodeQualityTool           = "codequalitytool"
	ResourceKindCodeQualityToolBinding    = "codequalitytoolbinding"
	ResourceKindCodeQualityBinding        = "codequalitybinding"
	ResourceKindToolBindingReplica        = "toolbindingreplica"
	ResourceKindProjectData               = "projectdata"
	ResourceKindProjectDataList           = "projectdatalist"
	// endregion

	// region Label: Camel nomenclature, used as the key in labels or other place
	LabelProject                  = ResourceKindProject
	LabelJenkins                  = ResourceKindJenkins
	LabelJenkinsBinding           = "jenkinsBinding"
	LabelPipelineConfig           = "pipelineConfig"
	LabelMultiBranch              = "multiBranchName"
	LabelTemplateKind             = "templateKind"
	LabelTemplateName             = "templateName"
	LabelTemplateVersion          = "templateVersion"
	LabelTemplateCategory         = "category"
	LabelPipeline                 = ResourceKindPipeline
	LabelPipelineKind             = "pipeline.kind"
	LabelPipelineKindMultiBranch  = "multi-branch"
	LabelCodeRepoService          = "codeRepoService"
	LabelCodeRepoBinding          = "codeRepoBinding"
	LabelCodeRepository           = "codeRepository"
	LabelProjectManagement        = "projectManagement"
	LabelProjectManagementBinding = "projectManagementBinding"
	LabelTestTool                 = "testTool"
	LabelTestToolBinding          = "testToolBinding"
	LabelImageRegistry            = "imageRegistry"
	LabelImageRegistryBinding     = "imageRegistryBinding"
	LabelImageRepository          = "imageRepository"
	LabelToolItemType             = "alauda.io/toolItemType"
	// LabelToolItemKind Kind of the tool, uses a an API kind
	LabelToolItemKind = "alauda.io/toolItemKind"
	// LabelToolItemPublic stablishes if the integrated tool is a public instance, e.g Github
	LabelToolItemPublic              = "alauda.io/toolItemPublic"
	LabelToolName                    = "alauda.io/toolName"
	LabelToolItemProject             = "alauda.io/toolItemProject"
	LabelToolBindingReplica          = "alauda.io/toolBindingReplica"
	LabelToolBindingReplicaNamespace = "alauda.io/toolBindingReplicaNamespace"
	LabelCodeQualityTool             = "codeQualityTool"
	LabelCodeQualityBinding          = "codeQualityBinding"
	LabelCodeQualityProject          = "codeQualityProject"
	LabelTemplateSource              = "source"
	LabelTemplateSourceOfficial      = "official"
	LabelTemplateSourceCustomer      = "customer"
	LabelTemplateLatest              = "alauda.io/latest"
	// LabelConfigTemplateCategory indicate category of template that created by graph template
	LabelConfigTemplateCategory = "PipelineConfigTemplate"

	// LabelDevopsAlaudaIOKey key used for specific Labels
	LabelDevopsAlaudaIOKey = "devops.alauda.io"
	// LabelDevopsAlaudaIOProjectKey key used for roles that are using in a project
	LabelDevopsAlaudaIOProjectKey = "devops.alauda.io/project"
	// LabelAlaudaIOProjectKey key used for project name
	LabelAlaudaIOProjectKey = "alauda.io/project"
	// LabelCodeRepoServiceType label key for codeRepoServiceType
	LabelCodeRepoServiceType = "codeRepoServiceType"
	// LabelCodeRepoServicePublic label key for codeRepoServicePublic
	LabelCodeRepoServicePublic = "codeRepoServicePublic"
	// LabelImageRegistryType label key for imageRegistryType
	LabelImageRegistryType = "imageRegistryType"
	// LabelImageRegistryEndpoint label key for imageRegistryEndpoint
	LabelImageRegistryEndpoint = "imageRegistryEndpoint"
	// LabelImageRepositoryLink label key for imageRepositoryLink
	LabelImageRepositoryLink = "imageRepositoryLink"
	// LabelsSecretName label key for specific secret name
	LabelsSecretName = "secretName"
	// LabelsSecretNamespace label key for specific secret namespace
	LabelsSecretNamespace = "secretNamespace"
	// LabelDevopsAlaudaIOGlobalKey key used for global secret
	LabelDevopsAlaudaIOGlobalKey = "devops.alauda.io/global"
	// LabelCodeQualityToolType label key for codeQualityToolType
	LabelCodeQualityToolType = "codeQualityToolType"
	// endregion

	// region ToolChainItem
	// configmap
	ConfigMapKindName         = "ConfigMap"
	ConfigMapAPIVersion       = "v1"
	SettingsConfigMapName     = "devops-config"
	SettingsKeyDomain         = "_domain"
	SettingsKeyGithubCreated  = "githubCreated"
	SettingsKeyToolChains     = "toolChains"
	SettingsKeyVersionGate    = "versionGate"
	SettingsKeyProduct        = "product"
	SettingsKeyACEEndpoint    = "ace_ui_endpoint"
	SettingsKeyACEAPIEndpoint = "ace_api_endpoint"
	SettingsKeyACEToken       = "ace_token"
	SettingsKeyACERootAccount = "ace_root_account"
	SettingsKeyRoleMapping    = "role_mapping"
	ACEEndpointDefault        = "https://cloud.alauda.cn/console/"
	ACEAPIEndpointDefault     = "https://api-cloud.alauda.cn/"

	// github
	GithubName          = "github"
	GithubHost          = "https://api.github.com"
	GithubHTML          = "https://github.com"
	GithubDisplayNameCN = "Github"
	GithubDisplayNameEN = "Github"

	// gitlab
	GitlabName          = "gitlab"
	GitlabHost          = "https://gitlab.com"
	GitlabHTML          = "https://gitlab.com"
	GitlabDisplayNameCN = "Gitlab"
	GitlabDisplayNameEN = "Gitlab"

	// gitlab-private
	GitlabPrivateName          = "gitlab-enterprise"
	GitlabPrivateHost          = ""
	GitlabPrivateHTML          = ""
	GitlabPrivateDisplayNameCN = "Gitlab 企业版"
	GitlabPrivateDisplayNameEN = "Gitlab Enterprise"

	// gitee
	GiteeName          = "gitee"
	GiteeHost          = "https://gitee.com"
	GiteeHTML          = "https://gitee.com"
	GiteeDisplayNameCN = "码云"
	GiteeDisplayNameEN = "Gitee"

	// gitee-private
	GiteePrivateName          = "gitee-enterprise"
	GiteePrivateHost          = ""
	GiteePrivateHTML          = ""
	GiteePrivateDisplayNameCN = "码云企业版"
	GiteePrivateDisplayNameEN = "Gitee Enterprise"

	// bitbucket
	BitbucketName          = "bitbucket"
	BitbucketHost          = "https://api.bitbucket.org"
	BitbucketHTML          = "https://bitbucket.org"
	BitbucketDisplayNameCN = "Bitbucket"
	BitbucketDisplayNameEN = "Bitbucket"

	// jira
	JiraName          = "jira"
	JiraDisplayNameCN = "Jira"
	JiraDisplayNameEN = "Jira"

	// taiga
	TaigaName          = "taiga"
	TaigaDisplayNameCN = "Taiga"
	TaigaDisplayNameEN = "Taiga"

	// redwoodhq
	RedwoodHQName          = "redwoodhq"
	RedwoodHQDisplayNameCN = "RedwoodHQ"
	RedwoodHQDisplayNameEN = "RedwoodHQ"

	// docker registry
	DockerRegistryName          = "docker-registry"
	DockerRegistryDisplayNameCN = "Docker Registry"
	DockerRegistryDisplayNameEN = "Docker Registry"

	// harbor registry
	HarborRegistryName          = "harbor-registry"
	HarborRegistryDisplayNameCN = "Harbor Registry"
	HarborRegistryDisplayNameEN = "Harbor Registry"

	// alauda registry
	AlaudaRegistryName          = "alauda-registry"
	AlaudaRegistryDisplayNameCN = "Alauda Registry"
	AlaudaRegistryDisplayNameEN = "Alauda Registry"

	// dockerhub registry
	DockerHubRegistryName          = "dockerhub-registry"
	DockerHubRegistryDisplayNameCN = "DockerHub Registry"
	DockerHubRegistryDisplayNameEN = "DockerHub Registry"
	DockerHubHTML                  = "https://hub.docker.com"
	DockerHubHost                  = "https://hub.docker.com"
	DockerHubRegistry              = "index.docker.io"

	// jenkins
	JenkinsName             = "jenkins"
	JenkinsDisplayNameCN    = "Jenkins"
	JenkinsDisplayNameEN    = "Jenkins"
	ToolChainJenkinsAPIPath = "jenkinses"

	// sonarqube
	SonarQubeName          = "sonarqube"
	SonarQubeDisplayNameCN = "SonarQube"
	SonarQubeDisplayNameEN = "SonarQube"

	// confluence
	ConfluenceName          = "confluence"
	ConfluenceDisplayNameCN = "confluence"
	ConfluenceDisplayNameEN = "confluence"

	// endregion

	// region ToolChainElement

	// codeRepository
	ToolChainCodeRepositoryName          = "codeRepository"
	ToolChainCodeRepositoryDisplayNameCN = "代码仓库"
	ToolChainCodeRepositoryDisplayNameEN = "Code Repository"
	ToolChainCodeRepositoryAPIPath       = "codereposervices"

	// continuousIntegration
	ToolChainContinuousIntegrationName          = "continuousIntegration"
	ToolChainContinuousIntegrationDisplayNameCN = "持续集成"
	ToolChainContinuousIntegrationDisplayNameEN = "Continuous Integration"
	ToolChainContinuousIntegrationAPIPath       = ""

	// artifactRepository
	ToolChainArtifactRepositoryName          = "artifactRepository"
	ToolChainArtifactRepositoryDisplayNameCN = "制品仓库"
	ToolChainArtifactRepositoryDisplayNameEN = "Artifact Repository"
	ToolChainArtifactRepositoryAPIPath       = ""
	ToolChainImageRegistryAPIPath            = "imageregistries"

	// testTool
	ToolChainTestToolName          = "testTool"
	ToolChainTestToolDisplayNameCN = "测试工具"
	ToolChainTestToolDisplayNameEN = "Test tool"
	ToolChainTestToolAPIPath       = "testtools"

	// projectManagement
	ToolChainProjectManagementName          = "projectManagement"
	ToolChainProjectManagementDisplayNameCN = "项目管理"
	ToolChainProjectManagementDisplayNameEN = "Project Management"
	ToolChainProjectManagementAPIPath       = "projectmanagements"

	// documentManagement
	ToolChainDocumentManagementName          = "documentManagement"
	ToolChainDocumentManagementDisplayNameCN = "文档管理"
	ToolChainDocumentManagementDisplayNameEN = "Document Management"
	ToolChainDocumentManagementAPIPath       = "documentmanagements"

	// codeQualityTool
	ToolChainCodeQualityToolName          = "codeQualityTool"
	ToolChainCodeQualityToolDisplayNameCN = "代码检查"
	ToolChainCodeQualityToolDisplayNameEN = "Code Quality"
	ToolChainCodeQualityToolAPIPath       = "codequalitytools"

	// endregion

	// region version gate options

	VersionGateGA    = "ga"
	VersionGateAlpha = "alpha"
	VersionGateBeta  = "beta"

	// endregion

	// TrueString true as string
	TrueString = "true"
	// FalseString false as string
	FalseString = "false"

	// APIVersionV1Alpha1 version of the api
	APIVersionV1Alpha1 = "devops.alauda.io/v1alpha1"

	// region TTL
	TTLSession              = 5 * time.Minute
	TTLServiceCheck         = 5 * time.Minute
	TTLDockerSecretSync     = 1 * time.Minute
	TTLCheckCodeRepoService = 5 * time.Minute
	TTLCheckCodeRepoBinding = 5 * time.Minute
	TTLCheckCodeRepository  = 5 * time.Minute
	TTLRoleSyncSession      = 20 * time.Minute

	TTLCheckCodeQualityTool    = 3 * time.Minute
	TTLCheckCodeQualityBinding = 3 * time.Minute

	TTLDevOpsGC = 60 * time.Minute
	// endregion

	// imageRegistry
	SettingsKeyImageRegistryTypes = "imageRegistryTypes"
	SettingsKeyDockerCreated      = "dockerCreated"

	// region OAuth2
	// SecretTypeOAuth2 contains data needed for oauth2 authentication.
	//
	// Required fields:
	// - Secret.Data["clientID"] - client id used for authentication
	// - Secret.Data["clientSecret"] - client secret used for authentication
	SecretTypeOAuth2 corev1.SecretType = "devops.alauda.io/oauth2"

	// OAuth2ClientIDKey is the key of the clientID for SecretTypeOAuth2 secrets
	OAuth2ClientIDKey = "clientID"
	// OAuth2ClientSecretKey is the key of the clientSecret for SecretTypeOAuth2 secrets
	OAuth2ClientSecretKey = "clientSecret"
	// OAuth2CodeKey is the key of the code for SecretTypeOAuth2 secrets
	OAuth2CodeKey = "code"
	// OAuth2AccessTokenKeyKey is the key of the accessTokenKey for SecretTypeOAuth2 secrets
	OAuth2AccessTokenKeyKey = "accessTokenKey"
	// OAuth2AccessTokenKey is the key of the accessToken for SecretTypeOAuth2 secrets
	OAuth2AccessTokenKey = "accessToken"
	// OAuth2ScopeKey is the key of the scope for SecretTypeOAuth2 secrets
	OAuth2ScopeKey = "scope"
	// OAuth2RefreshTokenKey is the key of the refreshToken for SecretTypeOAuth2 secrets
	OAuth2RefreshTokenKey = "refreshToken"
	// OAuth2ExpiresInKey is the key of the expiresIn for SecretTypeOAuth2 secrets
	OAuth2CreatedAtKey = "createdAt"
	// OAuth2ExpiresInKey is the key of the expiresIn for SecretTypeOAuth2 secrets
	OAuth2ExpiresInKey = "expiresIn"
	// OAuth2RedirectURLKey is the key of the redirectURL for SecretTypeOAuth2 secrets
	OAuth2RedirectURLKey = "redirectURL"
	// endregion

	NamespaceAlaudaSystem      = "alauda-system"
	NamespaceDefault           = "default"
	NamespaceKubeSystem        = "kube-system"
	NamespaceGlobalCredentials = "global-credentials"

	EventReasonSuccessUpdated = "Updated"
	EventMessageFmtUpdated    = "%s was updated"

	// resource name
	ResourceNamespaces      = "namespaces"
	ResourceServiceAccounts = "serviceaccounts"

	TypeCodeQualityReport = "CodeQualityReport"

	// SonarQube Metric Names
	SonarQubeBugs               = "bugs"
	SonarQubeVulnerabilities    = "vulnerabilities"
	SonarQubeCodeSmells         = "codeSmells"
	SonarQubeDuplication        = "duplications"
	SonarQubeCoverage           = "coverage"
	SonarQubeLanguages          = "languages"
	SonarQubeNewBugs            = "newBugs"
	SonarQubeNewCodeSmells      = "newCodeSmells"
	SonarQubeNewDuplication     = "newDuplications"
	SonarQubeNewVulnerabilities = "newVulnerabilities"
	SonarQubeNewCoverage        = "newCoverage"

	ClusterTaskTemplateTypePrefix = "public/"

	// PipelineConfigTemplateDefaultVersion is default template version when the config is create by graph,
	// and template is  auto created by pipelineconfig
	PipelineConfigTemplateDefaultVersion = "0.1"

	FinalizerPipelineConfigReferenced   = "pipelineconfig-referenced"
	FinalizerPipelineTemplateReferenced = "pipelinetemplate-referenced"
)
