//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021.

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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirGCSConnector) DeepCopyInto(out *FhirGCSConnector) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirGCSConnector.
func (in *FhirGCSConnector) DeepCopy() *FhirGCSConnector {
	if in == nil {
		return nil
	}
	out := new(FhirGCSConnector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FhirGCSConnector) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirGCSConnectorList) DeepCopyInto(out *FhirGCSConnectorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FhirGCSConnector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirGCSConnectorList.
func (in *FhirGCSConnectorList) DeepCopy() *FhirGCSConnectorList {
	if in == nil {
		return nil
	}
	out := new(FhirGCSConnectorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FhirGCSConnectorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirGCSConnectorSpec) DeepCopyInto(out *FhirGCSConnectorSpec) {
	*out = *in
	out.FhirStoreSelector = in.FhirStoreSelector
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirGCSConnectorSpec.
func (in *FhirGCSConnectorSpec) DeepCopy() *FhirGCSConnectorSpec {
	if in == nil {
		return nil
	}
	out := new(FhirGCSConnectorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirGCSConnectorSpecFHIRSelector) DeepCopyInto(out *FhirGCSConnectorSpecFHIRSelector) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirGCSConnectorSpecFHIRSelector.
func (in *FhirGCSConnectorSpecFHIRSelector) DeepCopy() *FhirGCSConnectorSpecFHIRSelector {
	if in == nil {
		return nil
	}
	out := new(FhirGCSConnectorSpecFHIRSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirGCSConnectorStatus) DeepCopyInto(out *FhirGCSConnectorStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirGCSConnectorStatus.
func (in *FhirGCSConnectorStatus) DeepCopy() *FhirGCSConnectorStatus {
	if in == nil {
		return nil
	}
	out := new(FhirGCSConnectorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirResource) DeepCopyInto(out *FhirResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirResource.
func (in *FhirResource) DeepCopy() *FhirResource {
	if in == nil {
		return nil
	}
	out := new(FhirResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FhirResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirResourceList) DeepCopyInto(out *FhirResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FhirResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirResourceList.
func (in *FhirResourceList) DeepCopy() *FhirResourceList {
	if in == nil {
		return nil
	}
	out := new(FhirResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FhirResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirResourceSpec) DeepCopyInto(out *FhirResourceSpec) {
	*out = *in
	out.Selector = in.Selector
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirResourceSpec.
func (in *FhirResourceSpec) DeepCopy() *FhirResourceSpec {
	if in == nil {
		return nil
	}
	out := new(FhirResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirResourceSpecFhirStoreSelector) DeepCopyInto(out *FhirResourceSpecFhirStoreSelector) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirResourceSpecFhirStoreSelector.
func (in *FhirResourceSpecFhirStoreSelector) DeepCopy() *FhirResourceSpecFhirStoreSelector {
	if in == nil {
		return nil
	}
	out := new(FhirResourceSpecFhirStoreSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirResourceStatus) DeepCopyInto(out *FhirResourceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirResourceStatus.
func (in *FhirResourceStatus) DeepCopy() *FhirResourceStatus {
	if in == nil {
		return nil
	}
	out := new(FhirResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirStore) DeepCopyInto(out *FhirStore) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirStore.
func (in *FhirStore) DeepCopy() *FhirStore {
	if in == nil {
		return nil
	}
	out := new(FhirStore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FhirStore) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirStoreList) DeepCopyInto(out *FhirStoreList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FhirStore, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirStoreList.
func (in *FhirStoreList) DeepCopy() *FhirStoreList {
	if in == nil {
		return nil
	}
	out := new(FhirStoreList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FhirStoreList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirStoreSpec) DeepCopyInto(out *FhirStoreSpec) {
	*out = *in
	if in.Auth != nil {
		in, out := &in.Auth, &out.Auth
		*out = make(map[string]FhirStoreSpecAuth, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	in.Options.DeepCopyInto(&out.Options)
	out.ExportOptions = in.ExportOptions
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirStoreSpec.
func (in *FhirStoreSpec) DeepCopy() *FhirStoreSpec {
	if in == nil {
		return nil
	}
	out := new(FhirStoreSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirStoreSpecAuth) DeepCopyInto(out *FhirStoreSpecAuth) {
	*out = *in
	if in.Members != nil {
		in, out := &in.Members, &out.Members
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirStoreSpecAuth.
func (in *FhirStoreSpecAuth) DeepCopy() *FhirStoreSpecAuth {
	if in == nil {
		return nil
	}
	out := new(FhirStoreSpecAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirStoreSpecExportOptions) DeepCopyInto(out *FhirStoreSpecExportOptions) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirStoreSpecExportOptions.
func (in *FhirStoreSpecExportOptions) DeepCopy() *FhirStoreSpecExportOptions {
	if in == nil {
		return nil
	}
	out := new(FhirStoreSpecExportOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirStoreSpecOptions) DeepCopyInto(out *FhirStoreSpecOptions) {
	*out = *in
	if in.Bigquery != nil {
		in, out := &in.Bigquery, &out.Bigquery
		*out = make([]FhirStoreSpecOptionsBigquery, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirStoreSpecOptions.
func (in *FhirStoreSpecOptions) DeepCopy() *FhirStoreSpecOptions {
	if in == nil {
		return nil
	}
	out := new(FhirStoreSpecOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirStoreSpecOptionsBigquery) DeepCopyInto(out *FhirStoreSpecOptionsBigquery) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirStoreSpecOptionsBigquery.
func (in *FhirStoreSpecOptionsBigquery) DeepCopy() *FhirStoreSpecOptionsBigquery {
	if in == nil {
		return nil
	}
	out := new(FhirStoreSpecOptionsBigquery)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FhirStoreStatus) DeepCopyInto(out *FhirStoreStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FhirStoreStatus.
func (in *FhirStoreStatus) DeepCopy() *FhirStoreStatus {
	if in == nil {
		return nil
	}
	out := new(FhirStoreStatus)
	in.DeepCopyInto(out)
	return out
}
