package test

import (
	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	"gopkg.in/h2non/gock.v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NameHost   = "dsaas"
	NameMember = "east"
)

func NewToolchainCluster(name, ons, secName string, status toolchainv1alpha1.ToolchainClusterStatus, labels map[string]string) (*toolchainv1alpha1.ToolchainCluster, *corev1.Secret) {
	return NewToolchainClusterWithEndpoint(name, ons, secName, "http://cluster.com", status, labels)
}

func NewToolchainClusterWithEndpoint(name, ons, secName, apiEndpoint string, status toolchainv1alpha1.ToolchainClusterStatus, labels map[string]string) (*toolchainv1alpha1.ToolchainCluster, *corev1.Secret) {
	gock.New(apiEndpoint).
		Get("api").
		Persist().
		Reply(200).
		BodyString("{}")
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      secName,
			Namespace: ons,
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"token": []byte("mycooltoken"),
		},
	}
	//labels["namespace"] = ons
	return &toolchainv1alpha1.ToolchainCluster{
		Spec: toolchainv1alpha1.ToolchainClusterSpec{
			SecretRef: toolchainv1alpha1.LocalSecretReference{
				Name: secName,
			},
			APIEndpoint: apiEndpoint,
			CABundle:    "",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: ons,
			Labels:    labels,
		},
		Status: status,
	}, secret
}

func NewClusterStatus(conType toolchainv1alpha1.ToolchainClusterConditionType, conStatus corev1.ConditionStatus) toolchainv1alpha1.ToolchainClusterStatus {
	return toolchainv1alpha1.ToolchainClusterStatus{
		Conditions: []toolchainv1alpha1.ToolchainClusterCondition{{
			Type:   conType,
			Status: conStatus,
		}},
	}
}
