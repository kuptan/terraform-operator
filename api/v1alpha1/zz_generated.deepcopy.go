//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Module) DeepCopyInto(out *Module) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Module.
func (in *Module) DeepCopy() *Module {
	if in == nil {
		return nil
	}
	out := new(Module)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OutputSpec) DeepCopyInto(out *OutputSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutputSpec.
func (in *OutputSpec) DeepCopy() *OutputSpec {
	if in == nil {
		return nil
	}
	out := new(OutputSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreviousRunStatus) DeepCopyInto(out *PreviousRunStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreviousRunStatus.
func (in *PreviousRunStatus) DeepCopy() *PreviousRunStatus {
	if in == nil {
		return nil
	}
	out := new(PreviousRunStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Terraform) DeepCopyInto(out *Terraform) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Terraform.
func (in *Terraform) DeepCopy() *Terraform {
	if in == nil {
		return nil
	}
	out := new(Terraform)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Terraform) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TerraformList) DeepCopyInto(out *TerraformList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Terraform, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TerraformList.
func (in *TerraformList) DeepCopy() *TerraformList {
	if in == nil {
		return nil
	}
	out := new(TerraformList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TerraformList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TerraformSpec) DeepCopyInto(out *TerraformSpec) {
	*out = *in
	out.Module = in.Module
	if in.DependsOn != nil {
		in, out := &in.DependsOn, &out.DependsOn
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Variables != nil {
		in, out := &in.Variables, &out.Variables
		*out = make([]Variable, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VariableFiles != nil {
		in, out := &in.VariableFiles, &out.VariableFiles
		*out = make([]VariableFile, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Outputs != nil {
		in, out := &in.Outputs, &out.Outputs
		*out = make([]*OutputSpec, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(OutputSpec)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TerraformSpec.
func (in *TerraformSpec) DeepCopy() *TerraformSpec {
	if in == nil {
		return nil
	}
	out := new(TerraformSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TerraformStatus) DeepCopyInto(out *TerraformStatus) {
	*out = *in
	if in.PreviousRuns != nil {
		in, out := &in.PreviousRuns, &out.PreviousRuns
		*out = make([]PreviousRunStatus, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TerraformStatus.
func (in *TerraformStatus) DeepCopy() *TerraformStatus {
	if in == nil {
		return nil
	}
	out := new(TerraformStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Variable) DeepCopyInto(out *Variable) {
	*out = *in
	if in.ValueFrom != nil {
		in, out := &in.ValueFrom, &out.ValueFrom
		*out = new(v1.EnvVarSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Variable.
func (in *Variable) DeepCopy() *Variable {
	if in == nil {
		return nil
	}
	out := new(Variable)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VariableFile) DeepCopyInto(out *VariableFile) {
	*out = *in
	if in.ValueFrom != nil {
		in, out := &in.ValueFrom, &out.ValueFrom
		*out = new(v1.VolumeSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VariableFile.
func (in *VariableFile) DeepCopy() *VariableFile {
	if in == nil {
		return nil
	}
	out := new(VariableFile)
	in.DeepCopyInto(out)
	return out
}
