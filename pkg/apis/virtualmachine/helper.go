package virtualmachine

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewVM(name, namespace string, options ...VMOption) *VirtualMachine {
	vm := &VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: VirtualMachineSpec{
			Template: &VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
				Spec: VirtualMachineInstanceSpec{
					Domain: DomainSpec{},
				},
			},
		},
	}

	for _, option := range options {
		option(vm)
	}

	return vm
}

type VMOption func(*VirtualMachine)

func WithRequests(requests corev1.ResourceList) VMOption {
	return func(vm *VirtualMachine) {
		vm.Spec.Template.Spec.Domain.Resources.Requests = requests
	}
}

func WithLimits(limits corev1.ResourceList) VMOption {
	return func(vm *VirtualMachine) {
		vm.Spec.Template.Spec.Domain.Resources.Limits = limits
	}
}
