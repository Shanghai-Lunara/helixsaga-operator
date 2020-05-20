// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppNotificationSpec) DeepCopyInto(out *AppNotificationSpec) {
	*out = *in
	if in.DispatchReplicas != nil {
		in, out := &in.DispatchReplicas, &out.DispatchReplicas
		*out = new(int32)
		**out = **in
	}
	if in.LogicReplicas != nil {
		in, out := &in.LogicReplicas, &out.LogicReplicas
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppNotificationSpec.
func (in *AppNotificationSpec) DeepCopy() *AppNotificationSpec {
	if in == nil {
		return nil
	}
	out := new(AppNotificationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CampaignSpec) DeepCopyInto(out *CampaignSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CampaignSpec.
func (in *CampaignSpec) DeepCopy() *CampaignSpec {
	if in == nil {
		return nil
	}
	out := new(CampaignSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonStatus) DeepCopyInto(out *CommonStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonStatus.
func (in *CommonStatus) DeepCopy() *CommonStatus {
	if in == nil {
		return nil
	}
	out := new(CommonStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GuildWarSpec) DeepCopyInto(out *GuildWarSpec) {
	*out = *in
	if in.RegisterReplicas != nil {
		in, out := &in.RegisterReplicas, &out.RegisterReplicas
		*out = new(int32)
		**out = **in
	}
	if in.GatewayReplicas != nil {
		in, out := &in.GatewayReplicas, &out.GatewayReplicas
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GuildWarSpec.
func (in *GuildWarSpec) DeepCopy() *GuildWarSpec {
	if in == nil {
		return nil
	}
	out := new(GuildWarSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelixSaga) DeepCopyInto(out *HelixSaga) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelixSaga.
func (in *HelixSaga) DeepCopy() *HelixSaga {
	if in == nil {
		return nil
	}
	out := new(HelixSaga)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HelixSaga) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelixSagaList) DeepCopyInto(out *HelixSagaList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HelixSaga, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelixSagaList.
func (in *HelixSagaList) DeepCopy() *HelixSagaList {
	if in == nil {
		return nil
	}
	out := new(HelixSagaList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HelixSagaList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelixSagaSpec) DeepCopyInto(out *HelixSagaSpec) {
	*out = *in
	in.VersionSpec.DeepCopyInto(&out.VersionSpec)
	in.ApiSpec.DeepCopyInto(&out.ApiSpec)
	in.GameSpec.DeepCopyInto(&out.GameSpec)
	in.PayNotifySpec.DeepCopyInto(&out.PayNotifySpec)
	in.GmtSpec.DeepCopyInto(&out.GmtSpec)
	in.FriendSpec.DeepCopyInto(&out.FriendSpec)
	in.QueueSpec.DeepCopyInto(&out.QueueSpec)
	in.RankSpec.DeepCopyInto(&out.RankSpec)
	in.ChatSpec.DeepCopyInto(&out.ChatSpec)
	in.HeartSpec.DeepCopyInto(&out.HeartSpec)
	in.CampaignSpec.DeepCopyInto(&out.CampaignSpec)
	in.GuildWarSpec.DeepCopyInto(&out.GuildWarSpec)
	in.AppNotificationSpec.DeepCopyInto(&out.AppNotificationSpec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelixSagaSpec.
func (in *HelixSagaSpec) DeepCopy() *HelixSagaSpec {
	if in == nil {
		return nil
	}
	out := new(HelixSagaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HelixSagaSpecStatus) DeepCopyInto(out *HelixSagaSpecStatus) {
	*out = *in
	out.VersionStatus = in.VersionStatus
	out.ApiStatus = in.ApiStatus
	out.GameStatus = in.GameStatus
	out.PayNotifyStatus = in.PayNotifyStatus
	out.GmtStatus = in.GmtStatus
	out.FriendStatus = in.FriendStatus
	out.QueueStatus = in.QueueStatus
	out.RankStatus = in.RankStatus
	out.ChatStatus = in.ChatStatus
	out.HeartStatus = in.HeartStatus
	out.CampaignStatus = in.CampaignStatus
	out.GuildWarStatus = in.GuildWarStatus
	out.AppNotificationStatus = in.AppNotificationStatus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HelixSagaSpecStatus.
func (in *HelixSagaSpecStatus) DeepCopy() *HelixSagaSpecStatus {
	if in == nil {
		return nil
	}
	out := new(HelixSagaSpecStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NginxPhpFpmSpec) DeepCopyInto(out *NginxPhpFpmSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NginxPhpFpmSpec.
func (in *NginxPhpFpmSpec) DeepCopy() *NginxPhpFpmSpec {
	if in == nil {
		return nil
	}
	out := new(NginxPhpFpmSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PhpSwooleSpec) DeepCopyInto(out *PhpSwooleSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PhpSwooleSpec.
func (in *PhpSwooleSpec) DeepCopy() *PhpSwooleSpec {
	if in == nil {
		return nil
	}
	out := new(PhpSwooleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PhpWorkermanSpec) DeepCopyInto(out *PhpWorkermanSpec) {
	*out = *in
	if in.RegisterReplicas != nil {
		in, out := &in.RegisterReplicas, &out.RegisterReplicas
		*out = new(int32)
		**out = **in
	}
	if in.GatewayReplicas != nil {
		in, out := &in.GatewayReplicas, &out.GatewayReplicas
		*out = new(int32)
		**out = **in
	}
	if in.BusinessWorkerReplicas != nil {
		in, out := &in.BusinessWorkerReplicas, &out.BusinessWorkerReplicas
		*out = new(int32)
		**out = **in
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PhpWorkermanSpec.
func (in *PhpWorkermanSpec) DeepCopy() *PhpWorkermanSpec {
	if in == nil {
		return nil
	}
	out := new(PhpWorkermanSpec)
	in.DeepCopyInto(out)
	return out
}