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
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Module struct {
	// the module source
	Source string `json:"source"`
	// module version
	// +optional
	Version string `json:"version,omitempty"`
}

type VariableFile struct {
	// Variable name
	Key string `json:"key"`

	// Source for the variable's value. Cannot be used if value is not empty.
	ValueFrom *corev1.VolumeSource `json:"valueFrom"`
}

type Variable struct {
	// Variable name
	Key string `json:"key"`
	// Variable value
	// +optional
	Value string `json:"value"`
	// Source for the variable's value. Cannot be used if value is not empty.
	// +optional
	ValueFrom *corev1.EnvVarSource `json:"valueFrom,omitempty"`
	// EnvironmentVariable denotes if this variable should be created as environment variable
	// +optional
	EnvironmentVariable bool `json:"environmentVariable,omitempty"`
}

// OutputSpec specifies which values need to be output
type OutputSpec struct {
	// Output name
	// +optional
	Key string `json:"key"`
	// Attribute name in module
	// +optional
	ModuleOutputName string `json:"moduleOutputName"`
}

// GitSSHKey config
type GitSSHKeySpec struct {
	ValueFrom *corev1.VolumeSource `json:"valueFrom"`
}

type TerraformRunStatus string

const (
	RunStarted   TerraformRunStatus = "Started"
	RunRunning   TerraformRunStatus = "Running"
	RunCompleted TerraformRunStatus = "Completed"
	RunDestroyed TerraformRunStatus = "Destroyed"
	RunFailed    TerraformRunStatus = "Failed"
)

// PreviousRuns stores the previous run information in case the current run object was modified
type PreviousRunStatus struct {
	// Attribute name in module
	// +optional
	RunId string `json:"id"`
	// Value
	// +optional
	Status TerraformRunStatus `json:"status"`
}

// TerraformSpec defines the desired state of Terraform
type TerraformSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Terraform. Edit terraform_types.go to remove/update
	TerraformVersion string `json:"terraformVersion"`
	// The module information to be provided
	Module Module `json:"module"`
	// A custom terrafor backend configuration
	// +optional
	Backend string `json:"backend,omitempty"`
	// A custom terrafor providers configuration
	// +optional
	ProvidersConfig string `json:"providersConfig,omitempty"`
	// The terraform workspae
	// +optional
	Workspace string `json:"workspace,omitempty"`
	// dependencies on other modules
	// +optional
	DependsOn []string `json:"dependsOn,omitempty"`
	// Variables as inputs to module
	// +optional
	Variables []Variable `json:"variables,omitempty"`
	// Variables as inputs to module
	// +optional
	VariableFiles []VariableFile `json:"variableFiles,omitempty"`
	// Outputs denote outputs wanted
	// +optional
	Outputs []*OutputSpec `json:"outputs,omitempty"`
	// Indicates whether a destroy job should run
	// +optional
	Destroy bool `json:"destroy,omitempty"`
	// Indicates whether to keep the jobs/pods after the run is successful
	// +optional
	DeleteCompletedJobs bool `json:"deleteCompletedJobs,omitempty"`
	// A retry limit to be set on the Job as a backOffLimit
	// +optional
	RetryLimit int32 `json:"retryLimit,omitempty"`
	// An SSH key to be able to run terraform on private git repositories
	// +optional
	GitSSHKey *GitSSHKeySpec `json:"gitSSHKey,omitempty"`
}

// TerraformStatus defines the observed state of Terraform
type TerraformStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	RunId        string              `json:"currentRunId"`
	PreviousRuns []PreviousRunStatus `json:"previousRuns,omitempty"`
	Generation   int64               `json:"generation"`
	RunStatus    TerraformRunStatus  `json:"runStatus"`
	Message      string              `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Terraform is the Schema for the terraforms API
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

// this evaluate the first time the object was created
func (t *Terraform) IsSubmitted() bool {
	return t.Status.Generation == 0 && t.Status.RunId == ""
}

// the run is either started or running
func (t *Terraform) IsStarted() bool {
	allowedStatuses := map[TerraformRunStatus]bool{
		RunStarted: true,
		RunRunning: true,
	}

	return allowedStatuses[t.Status.RunStatus]
}

// check if the status is running
func (t *Terraform) IsRunning() bool {
	return t.Status.RunStatus == RunRunning
}

// check if the object was updated
func (t *Terraform) IsUpdated() bool {
	return t.Generation > 0 && t.Generation > t.Status.Generation
}

func (t *Terraform) HasErrored() bool {
	return t.Status.RunStatus == RunFailed
}

func (r *Terraform) SetRunId() {
	r.Status.RunId = fmt.Sprint(random(8))
}

func (t *Terraform) PrepareForUpdate() {
	if len(t.Status.PreviousRuns) == 0 {
		t.Status.PreviousRuns = []PreviousRunStatus{}
	}

	t.Status.PreviousRuns = append(t.Status.PreviousRuns, PreviousRunStatus{
		RunId:  t.Status.RunId,
		Status: t.Status.RunStatus,
	})
}

// returns the owner reference
func (t *Terraform) GetOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: fmt.Sprintf("%s/%s", GroupVersion.Group, GroupVersion.Version),
		Kind:       t.Kind,
		Name:       t.Name,
		UID:        t.GetUID(),
	}
}

const runnerRBACName string = "terraform-runner"

// Creates a terraform Run as a Kubernetes job
func (t *Terraform) CreateTerraformRun(namespacedName types.NamespacedName) (*batchv1.Job, error) {
	if err := createRbacConfigIfNotExist(runnerRBACName, namespacedName.Namespace); err != nil {
		return nil, err
	}

	configMap, err := createConfigMapForModule(namespacedName, t)

	if err != nil {
		return nil, err
	}

	secret, err := createSecretForOutputs(namespacedName, t)

	if err != nil {
		return nil, err
	}

	job, err := createJobForRun(t, configMap, secret)

	if err != nil {
		return nil, err
	}

	return job, nil
}

func (t *Terraform) DeleteAfterCompletion() error {
	if err := deleteJobByRun(t.Name, t.Namespace, t.Status.RunId); err != nil {
		return err
	}

	return nil
}

func (t *Terraform) GetJobByRun() (*batchv1.Job, error) {
	job, err := getJobForRun(t.Name, t.Namespace, t.Status.RunId)

	if err != nil {
		return nil, err
	}

	return job, err
}

func init() {
	SchemeBuilder.Register(&Terraform{}, &TerraformList{})
}
