// copied from https://github.com/kubevirt/api/blob/main/core/v1/deepcopy_generated.go and stripped down to include only what's needed for testing
package virtualmachine

import (
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudInitConfigDriveSource) DeepCopyInto(out *CloudInitConfigDriveSource) {
	*out = *in
	if in.UserDataSecretRef != nil {
		in, out := &in.UserDataSecretRef, &out.UserDataSecretRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.NetworkDataSecretRef != nil {
		in, out := &in.NetworkDataSecretRef, &out.NetworkDataSecretRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudInitConfigDriveSource.
func (in *CloudInitConfigDriveSource) DeepCopy() *CloudInitConfigDriveSource {
	if in == nil {
		return nil
	}
	out := new(CloudInitConfigDriveSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudInitNoCloudSource) DeepCopyInto(out *CloudInitNoCloudSource) {
	*out = *in
	if in.UserDataSecretRef != nil {
		in, out := &in.UserDataSecretRef, &out.UserDataSecretRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.NetworkDataSecretRef != nil {
		in, out := &in.NetworkDataSecretRef, &out.NetworkDataSecretRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudInitNoCloudSource.
func (in *CloudInitNoCloudSource) DeepCopy() *CloudInitNoCloudSource {
	if in == nil {
		return nil
	}
	out := new(CloudInitNoCloudSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DomainSpec) DeepCopyInto(out *DomainSpec) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DomainSpec.
func (in *DomainSpec) DeepCopy() *DomainSpec {
	if in == nil {
		return nil
	}
	out := new(DomainSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRequirements) DeepCopyInto(out *ResourceRequirements) {
	*out = *in
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = make(corev1.ResourceList, len(*in))
		for key, val := range *in {
			(*out)[key] = val.DeepCopy()
		}
	}
	if in.Limits != nil {
		in, out := &in.Limits, &out.Limits
		*out = make(corev1.ResourceList, len(*in))
		for key, val := range *in {
			(*out)[key] = val.DeepCopy()
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRequirements.
func (in *ResourceRequirements) DeepCopy() *ResourceRequirements {
	if in == nil {
		return nil
	}
	out := new(ResourceRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachine) DeepCopyInto(out *VirtualMachine) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachine.
func (in *VirtualMachine) DeepCopy() *VirtualMachine {
	if in == nil {
		return nil
	}
	out := new(VirtualMachine)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachine) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineInstanceSpec) DeepCopyInto(out *VirtualMachineInstanceSpec) {
	*out = *in
	in.Domain.DeepCopyInto(&out.Domain)
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineInstanceSpec.
func (in *VirtualMachineInstanceSpec) DeepCopy() *VirtualMachineInstanceSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineInstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineInstanceTemplateSpec) DeepCopyInto(out *VirtualMachineInstanceTemplateSpec) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineInstanceTemplateSpec.
func (in *VirtualMachineInstanceTemplateSpec) DeepCopy() *VirtualMachineInstanceTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineInstanceTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineSpec) DeepCopyInto(out *VirtualMachineSpec) {
	*out = *in
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = new(VirtualMachineInstanceTemplateSpec)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineSpec.
func (in *VirtualMachineSpec) DeepCopy() *VirtualMachineSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Volume) DeepCopyInto(out *Volume) {
	*out = *in
	in.VolumeSource.DeepCopyInto(&out.VolumeSource)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Volume.
func (in *Volume) DeepCopy() *Volume {
	if in == nil {
		return nil
	}
	out := new(Volume)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeSource) DeepCopyInto(out *VolumeSource) {
	*out = *in
	if in.CloudInitNoCloud != nil {
		in, out := &in.CloudInitNoCloud, &out.CloudInitNoCloud
		*out = new(CloudInitNoCloudSource)
		(*in).DeepCopyInto(*out)
	}
	if in.CloudInitConfigDrive != nil {
		in, out := &in.CloudInitConfigDrive, &out.CloudInitConfigDrive
		*out = new(CloudInitConfigDriveSource)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeSource.
func (in *VolumeSource) DeepCopy() *VolumeSource {
	if in == nil {
		return nil
	}
	out := new(VolumeSource)
	in.DeepCopyInto(out)
	return out
}
