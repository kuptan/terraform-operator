/*
Copyright 2022.

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

package v1alpha1

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// TerraformFinalizer is the finalizer name
const TerraformFinalizer string = "finalizers.terraform-operator.io"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Module holds the Terraform module source and version information
type Module struct {
	// module source, must be a valid Terraform module source
	Source string `json:"source"`
	// module version
	// +optional
	Version string `json:"version,omitempty"`
}

// VariableFile holds the information of the Terraform variable files to include
type VariableFile struct {
	// The module variable name
	Key string `json:"key"`

	// The source of the variable file
	ValueFrom *corev1.VolumeSource `json:"valueFrom"`
}

// TerraformDependencyRef holds the information of the Terraform dependency name and key for the module
// to use as a variable
type TerraformDependencyRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// Variable holds the information of the Terraform variable
type Variable struct {
	// Terraform module variable name
	Key string `json:"key"`
	// The value of the variable
	// +optional
	Value string `json:"value"`
	// The variable value from a key source (secret or configmap)
	// +optional
	ValueFrom *corev1.EnvVarSource `json:"valueFrom,omitempty"`
	// EnvironmentVariable denotes if this variable should be created as environment variable
	// +optional
	EnvironmentVariable bool `json:"environmentVariable,omitempty"`
	// DependencyRef denotes if this variable should be fetched from the output of a dependency
	// +optional
	DependencyRef *TerraformDependencyRef `json:"dependencyRef,omitempty"`
}

// Output holds the information of the Terraform output information
// that will be written to a Kubernetes secret
type Output struct {
	// Output key specifies the Kubernetes secret key
	// +optional
	Key string `json:"key"`
	// The output name as defined in the source Terraform module
	// +optional
	ModuleOutputName string `json:"moduleOutputName"`
}

// DependsOn holds the information of the Terraform dependency
type DependsOn struct {
	// The Terraform object metadata.name
	Name string `json:"name"`
	// The namespace where the Terraform run exist
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// GitSSHKey holds the information of the Git SSH key
type GitSSHKey struct {
	// The source of the value where the private SSH key exist
	ValueFrom *corev1.VolumeSource `json:"valueFrom"`
}

// TerraformRunStatus is the status of the workflow/run
type TerraformRunStatus string

// workflow/run statuses
const (
	RunStarted              TerraformRunStatus = "Started"
	RunRunning              TerraformRunStatus = "Running"
	RunCompleted            TerraformRunStatus = "Completed"
	RunFailed               TerraformRunStatus = "Failed"
	RunWaitingForDependency TerraformRunStatus = "WaitingForDependency"
	RunDeleted              TerraformRunStatus = "Deleted"
)

// PreviousRunStatus stores the previous workflows/runs information
// in case the current workflow/run object was modified
type PreviousRunStatus struct {
	// Attribute name in module
	// +optional
	RunID string `json:"id"`
	// Value
	// +optional
	Status TerraformRunStatus `json:"status"`
}

// TerraformSpec defines the desired state of Terraform object
type TerraformSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The terraform version to use
	TerraformVersion string `json:"terraformVersion"`
	// The module information (source & version)
	Module Module `json:"module"`
	// A custom terraform backend configuration
	// +optional
	Backend string `json:"backend,omitempty"`
	// A custom terraform providers configuration
	// +optional
	ProvidersConfig string `json:"providersConfig,omitempty"`
	// The terraform workspae. Defaults to `default`
	// +optional
	Workspace string `json:"workspace,omitempty"`
	// A list of dependencies on other Terraform runs
	// +optional
	DependsOn []*DependsOn `json:"dependsOn,omitempty"`
	// Variables as inputs to the Terraform module
	// +optional
	Variables []Variable `json:"variables,omitempty"`
	// Terraform variable files
	// +optional
	VariableFiles []VariableFile `json:"variableFiles,omitempty"`
	// Terraform outputs will be written to a Kubernetes secret
	// +optional
	Outputs []*Output `json:"outputs,omitempty"`
	// Indicates whether a destroy job should run
	// +optional
	Destroy bool `json:"destroy,omitempty"`
	// Indicates whether to keep the jobs/pods after the run is successful/completed
	// +optional
	DeleteCompletedJobs bool `json:"deleteCompletedJobs,omitempty"`
	// A retry limit to be set on the Job as a backOffLimit
	// +optional
	RetryLimit int32 `json:"retryLimit,omitempty"`
	// An SSH key to be able to pull modules from private git repositories
	// +optional
	GitSSHKey *GitSSHKey `json:"gitSSHKey,omitempty"`
}

// TerraformStatus defines the observed state of Terraform
type TerraformStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	RunID              string             `json:"currentRunId"`
	PreviousRunID      string             `json:"previousRunId,omitempty"`
	OutputSecretName   string             `json:"outputSecretName,omitempty"`
	ObservedGeneration int64              `json:"observedGeneration"`
	RunStatus          TerraformRunStatus `json:"runStatus"`
	Message            string             `json:"message,omitempty"`
	StartedTime        string             `json:"startTime,omitempty"`
	CompletionTime     string             `json:"completionTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Terraform is the Schema for the terraforms API
// +kubebuilder:resource:shortName=tf,path=terraforms
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.runStatus"
// +kubebuilder:printcolumn:name="Secret",type="string",JSONPath=".status.outputSecretName"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Terraform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TerraformSpec   `json:"spec,omitempty"`
	Status TerraformStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TerraformList contains a list of Terraform
type TerraformList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Terraform `json:"items"`
}

// IsSubmitted evaluates if the workflow/run is created for the first time
func (t *Terraform) IsSubmitted() bool {
	return t.Status.RunID == ""
}

// IsStarted evaluates that the workflow/run is started
func (t *Terraform) IsStarted() bool {
	allowedStatuses := map[TerraformRunStatus]bool{
		RunStarted: true,
		RunRunning: true,
	}

	return allowedStatuses[t.Status.RunStatus]
}

// IsRunning evaluates that the workflow/run is running
func (t *Terraform) IsRunning() bool {
	return t.Status.RunStatus == RunRunning
}

// IsUpdated evaluates if the workflow/run was updated
func (t *Terraform) IsUpdated() bool {
	return t.Generation > 0 && t.Generation > t.Status.ObservedGeneration
}

// IsWaiting evaluates if the workflow/run is waiting for a dependency
func (t *Terraform) IsWaiting() bool {
	return t.Status.RunStatus == RunWaitingForDependency
}

// HasErrored evaluates if the workflow/run failed
func (t *Terraform) HasErrored() bool {
	return t.Status.RunStatus == RunFailed
}

// SetRunID sets a new value for the run ID
func (t *Terraform) SetRunID() {
	if t.Status.RunID != "" {
		t.Status.PreviousRunID = t.Status.RunID
	}
	t.Status.RunID = random(6)
}

// GetOwnerReference returns the Kubernetes owner reference meta
func (t *Terraform) GetOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: fmt.Sprintf("%s/%s", GroupVersion.Group, GroupVersion.Version),
		Kind:       t.Kind,
		Name:       t.Name,
		UID:        t.GetUID(),
	}
}

// setBackendCfgIfNotExist sets the default backend to Kunernetes if not provided
func setBackendCfgIfNotExist(run *Terraform) {
	if run.Spec.Backend == "" {
		run.Spec.Backend = fmt.Sprintf(`backend "kubernetes" {
  secret_suffix     = "%s"
  in_cluster_config = true
	namespace					= "%s"
}
`, run.ObjectMeta.Name, run.ObjectMeta.Namespace)
	}
}

// runnerRBACName is the RBAC name that will be used in the role and service account creation
// if they're not found
const runnerRBACName string = "terraform-runner"

// CreateTerraformRun creates the Kubernetes objects to start the workflow/run
//
// (RBAC (service account & Role), ConfigMap for the terraform module file,
// Secret to store the outputs if any, will be empty if no outputs are defined,
// Job to execute the workflow/run)
func (t *Terraform) CreateTerraformRun(ctx context.Context, namespacedName types.NamespacedName) (*batchv1.Job, error) {
	setBackendCfgIfNotExist(t)

	if err := createRbacConfigIfNotExist(ctx, runnerRBACName, namespacedName.Namespace); err != nil {
		return nil, err
	}

	_, err := createConfigMapForModule(ctx, namespacedName, t)

	if err != nil {
		return nil, err
	}

	_, err = createSecretForOutputs(ctx, namespacedName, t)

	if err != nil {
		return nil, err
	}

	job, err := createJobForRun(ctx, t)

	if err != nil {
		return nil, err
	}

	return job, nil
}

// DeleteAfterCompletion removes the Kubernetes of the workflow/run once completed
func (t *Terraform) DeleteAfterCompletion(ctx context.Context) error {
	if err := deleteJobByRun(ctx, t.Name, t.Namespace, t.Status.RunID); err != nil {
		return err
	}

	return nil
}

// GetOutputSecretName returns the secret name of the Terraform outputs
func (t *Terraform) GetOutputSecretName() string {
	return getOutputSecretname(t.Name)
}

// CleanupResources cleans up old resources (secrets & configmaps)
func (t *Terraform) CleanupResources(ctx context.Context) error {
	previousRunID := t.Status.PreviousRunID

	if previousRunID == "" {
		return nil
	}

	// delete the older job
	if err := deleteJobByRun(ctx, t.Name, t.Namespace, previousRunID); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	// delete the older configmap that holds the module
	if err := deleteConfigMapByRun(ctx, t.Name, t.Namespace, previousRunID); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}

// GetJobByRun returns the Kubernetes job of the workflow/run
func (t *Terraform) GetJobByRun(ctx context.Context) (*batchv1.Job, error) {
	job, err := getJobForRun(ctx, t.Name, t.Namespace, t.Status.RunID)

	if err != nil {
		return nil, err
	}

	return job, err
}

// Init initializes the scheme builder
func init() {
	SchemeBuilder.Register(&Terraform{}, &TerraformList{})
}
