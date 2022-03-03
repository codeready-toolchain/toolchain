package nstemplateset

import (
	"sort"
	"time"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	"github.com/codeready-toolchain/toolchain-common/pkg/test"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	clusterResourcesTemplateRef = "basic-clusterresources-abcde00"
	devTemplateRef              = "basic-dev-abcde11"
	codeTemplateRef             = "basic-code-abcde21"
)

type Option func(*toolchainv1alpha1.NSTemplateSet)

func WithReferencesFor(nstemplateTier *toolchainv1alpha1.NSTemplateTier) Option {
	return func(nstmplSet *toolchainv1alpha1.NSTemplateSet) {
		nstmplSet.Spec.TierName = nstemplateTier.Name

		// cluster resources
		if nstemplateTier.Spec.ClusterResources != nil {
			nstmplSet.Spec.ClusterResources = &toolchainv1alpha1.NSTemplateSetClusterResources{
				TemplateRef: nstemplateTier.Spec.ClusterResources.TemplateRef,
			}
		}

		// namespace resources
		if len(nstemplateTier.Spec.Namespaces) > 0 {
			nstmplSet.Spec.Namespaces = make([]toolchainv1alpha1.NSTemplateSetNamespace, len(nstemplateTier.Spec.Namespaces))
			for i, ns := range nstemplateTier.Spec.Namespaces {
				nstmplSet.Spec.Namespaces[i] = toolchainv1alpha1.NSTemplateSetNamespace(ns)
			}
		}

		// space roles
		// append by alphabetical order of role names
		if len(nstemplateTier.Spec.SpaceRoles) > 0 {
			roles := []string{}
			for r := range nstemplateTier.Spec.SpaceRoles {
				roles = append(roles, r)
			}
			sort.Strings(roles)
			for _, r := range roles {
				nstmplSet.Spec.SpaceRoles = append(nstmplSet.Spec.SpaceRoles, toolchainv1alpha1.NSTemplateSetSpaceRole{
					TemplateRef: nstemplateTier.Spec.SpaceRoles[r].TemplateRef,
					// TODO: include usernames from SpaceBindings
				})
			}
		}
	}
}

func WithReadyCondition() Option {
	return func(nstmplSet *toolchainv1alpha1.NSTemplateSet) {
		nstmplSet.Status.Conditions = []toolchainv1alpha1.Condition{
			{
				Type:   toolchainv1alpha1.ConditionReady,
				Status: corev1.ConditionTrue,
				Reason: toolchainv1alpha1.NSTemplateSetProvisionedReason,
			},
		}
	}
}

func WithNotReadyCondition(reason, message string) Option {
	return func(nstmplSet *toolchainv1alpha1.NSTemplateSet) {
		nstmplSet.Status.Conditions = []toolchainv1alpha1.Condition{
			{
				Type:    toolchainv1alpha1.ConditionReady,
				Status:  corev1.ConditionFalse,
				Reason:  reason,
				Message: message,
			},
		}
	}
}

func WithDeletionTimestamp(ts time.Time) Option {
	return func(nstmplSet *toolchainv1alpha1.NSTemplateSet) {
		nstmplSet.DeletionTimestamp = &metav1.Time{Time: ts}
	}
}

func NewNSTemplateSet(name string, options ...Option) *toolchainv1alpha1.NSTemplateSet {
	nstmplSet := &toolchainv1alpha1.NSTemplateSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: test.MemberOperatorNs,
			Name:      name,
		},
		Spec: toolchainv1alpha1.NSTemplateSetSpec{
			TierName: "basic",
			ClusterResources: &toolchainv1alpha1.NSTemplateSetClusterResources{
				TemplateRef: clusterResourcesTemplateRef,
			},
			Namespaces: []toolchainv1alpha1.NSTemplateSetNamespace{
				{
					TemplateRef: devTemplateRef,
				},
				{
					TemplateRef: codeTemplateRef,
				},
			},
		},
		Status: toolchainv1alpha1.NSTemplateSetStatus{},
	}
	for _, apply := range options {
		apply(nstmplSet)
	}
	return nstmplSet
}
