package template_test

import (
	"embed"
	"testing"

	"github.com/codeready-toolchain/toolchain-common/pkg/template"
	"github.com/codeready-toolchain/toolchain-common/pkg/test"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

//go:embed testdata/*
var EFS embed.FS

//go:embed testdata/host/*
var hostFS embed.FS

//go:embed testdata/member/*
var memberFS embed.FS

func TestLoadObjectsFromEmbedFS(t *testing.T) {
	t.Run("loads objects recursively from all subdirectories", func(t *testing.T) {
		// when
		allObjects, err := template.LoadObjectsFromEmbedFS(&EFS, &template.Variables{Namespace: test.HostOperatorNs})
		require.NoError(t, err)
		hostFolderObjects, err := template.LoadObjectsFromEmbedFS(&hostFS, &template.Variables{Namespace: test.HostOperatorNs})
		require.NoError(t, err)
		memberFolderObjects, err := template.LoadObjectsFromEmbedFS(&memberFS, nil)
		require.NoError(t, err)
		// then
		require.NotNil(t, allObjects)
		require.NotNil(t, hostFolderObjects)
		require.NotNil(t, memberFolderObjects)
		require.Len(t, allObjects, 4, "invalid number of expected total objects")
		require.Len(t, hostFolderObjects, 3, "invalid number of expected objects from host folder")
		require.Len(t, memberFolderObjects, 1, "invalid number of expected objects from member folder")
		// check match for the expected objects
		checkExpectedObjects(t, allObjects)
	})

	t.Run("error - when variables are not provided", func(t *testing.T) {
		// when
		// we do not pass required variables for the templates that requires variables
		objects, err := template.LoadObjectsFromEmbedFS(&hostFS, nil)
		// then
		// we should get back an error
		require.Error(t, err)
		require.Nil(t, objects)
	})
}

func checkExpectedObjects(t *testing.T, objects []*unstructured.Unstructured) {
	sa := &v1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(objects[0].Object, sa)
	require.NoError(t, err)
	require.Equal(t, "toolchaincluster-host", sa.GetName())
	require.Equal(t, "toolchain-host-operator", sa.GetNamespace())
	role := &rbac.Role{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(objects[1].Object, role)
	require.NoError(t, err)
	require.Equal(t, "toolchaincluster-host", role.GetName())
	require.Equal(t, "toolchain-host-operator", role.GetNamespace())
	require.Equal(t, []rbac.PolicyRule{
		{
			APIGroups: []string{"toolchain.dev.openshift.com"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		},
	}, role.Rules)
	roleBinding := &rbac.RoleBinding{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(objects[2].Object, roleBinding)
	require.NoError(t, err)
	require.Equal(t, "toolchaincluster-host", roleBinding.GetName())
	require.Equal(t, "toolchain-host-operator", roleBinding.GetNamespace())
	require.Equal(t, rbac.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "toolchaincluster-host",
	}, roleBinding.RoleRef)
	require.Len(t, roleBinding.Subjects, 1)
	require.Equal(t, rbac.Subject{
		Kind: "ServiceAccount",
		Name: "toolchaincluster-host",
	}, roleBinding.Subjects[0])
	clusterRole := &rbac.ClusterRole{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(objects[3].Object, clusterRole)
	require.NoError(t, err)
	require.Equal(t, "member-toolchaincluster-cr", clusterRole.GetName())
	require.Equal(t, []rbac.PolicyRule{
		{
			APIGroups: []string{"authentication.k8s.io"},
			Resources: []string{"tokenreviews"},
			Verbs:     []string{"create"},
		},
	}, clusterRole.Rules)
}
