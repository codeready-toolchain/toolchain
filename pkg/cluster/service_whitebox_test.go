package cluster

import (
	"testing"
	"time"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	"github.com/codeready-toolchain/toolchain-common/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func TestRefreshCacheInService(t *testing.T) {
	// given
	defer gock.Off()
	status := test.NewClusterStatus(toolchainv1alpha1.ToolchainClusterReady, corev1.ConditionTrue)
	toolchainCluster, sec := test.NewToolchainCluster("east", test.HostOperatorNs, "secret", status, map[string]string{"ownerClusterName": test.NameMember, "namespace": test.MemberOperatorNs})
	s := scheme.Scheme
	err := toolchainv1alpha1.AddToScheme(s)
	require.NoError(t, err)
	cl := test.NewFakeClient(t, toolchainCluster, sec)
	service := newToolchainClusterService(cl, 0, test.HostOperatorNs)

	t.Run("the member cluster should be retrieved when refreshCache func is called", func(t *testing.T) {
		// given
		_, ok := clusterCache.clusters["east"]
		require.False(t, ok)
		defer service.DeleteToolchainCluster("east")

		// when
		service.refreshCache()

		// then
		cachedCluster, ok := GetCachedToolchainCluster(test.NameMember)
		require.True(t, ok)
		assertMemberCluster(t, cachedCluster, status)
	})

	t.Run("the member cluster should be retrieved when GetCachedToolchainCluster func is called", func(t *testing.T) {
		// given
		_, ok := clusterCache.clusters["east"]
		require.False(t, ok)
		defer service.DeleteToolchainCluster("east")

		// when
		cachedCluster, ok := GetCachedToolchainCluster(test.NameMember)

		// then
		require.True(t, ok)
		assertMemberCluster(t, cachedCluster, status)
	})

	t.Run("the host cluster should not be retrieved", func(t *testing.T) {
		// given
		_, ok := clusterCache.clusters["east"]
		require.False(t, ok)
		defer service.DeleteToolchainCluster("east")

		// when
		cachedCluster, ok := GetCachedToolchainCluster(test.NameHost)

		// then
		require.False(t, ok)
		assert.Nil(t, cachedCluster)
	})
}

func TestUpdateClientBasedOnRestConfig(t *testing.T) {
	// given
	defer gock.Off()
	statusTrue := test.NewClusterStatus(toolchainv1alpha1.ToolchainClusterReady, corev1.ConditionTrue)
	toolchainCluster1, sec1 := test.NewToolchainCluster("east", test.HostOperatorNs, "secret1", statusTrue,
		map[string]string{"namespace": test.HostOperatorNs})

	t.Run("don't update when RestConfig is the same", func(t *testing.T) {
		// given
		cl := test.NewFakeClient(t, sec1)
		service := newToolchainClusterService(cl, 3*time.Second, test.HostOperatorNs)
		defer service.DeleteToolchainCluster("east")

		err := service.AddOrUpdateToolchainCluster(toolchainCluster1)
		require.NoError(t, err)
		originalClient := clusterCache.clusters["east"].Client
		clusterCache.clusters["east"].Client = cl

		// when
		err = service.AddOrUpdateToolchainCluster(toolchainCluster1)
		require.NoError(t, err)

		// then
		require.NoError(t, err)
		cachedToolchainCluster, ok := GetCachedToolchainCluster("east")
		require.True(t, ok)
		assert.NotEqual(t, originalClient, cachedToolchainCluster.Client)
		assert.Equal(t, cl, cachedToolchainCluster.Client)
	})

	t.Run("update when RestConfig is not the same", func(t *testing.T) {
		// given
		cl := test.NewFakeClient(t, sec1)
		service := newToolchainClusterService(cl, 3*time.Second, test.HostOperatorNs)
		defer service.DeleteToolchainCluster("east")

		err := service.AddOrUpdateToolchainCluster(toolchainCluster1)
		require.NoError(t, err)
		clusterCache.clusters["east"].Client = cl
		clusterCache.clusters["east"].RestConfig.BearerToken = "old-token"

		// when
		err = service.AddOrUpdateToolchainCluster(toolchainCluster1)
		require.NoError(t, err)

		// then
		require.NoError(t, err)
		cachedToolchainCluster, ok := GetCachedToolchainCluster("east")
		require.True(t, ok)
		assert.NotEqual(t, cl, cachedToolchainCluster.Client)
	})
}

func newToolchainClusterService(cl client.Client, timeout time.Duration, tcNs string) ToolchainClusterService {
	return NewToolchainClusterServiceWithClient(cl, logf.Log, tcNs, timeout, func(config *rest.Config, options client.Options) (client.Client, error) {
		// make sure that insecure is false to make Gock mocking working properly
		// let's use a copy of the config, so it doesn't affect the cache logic
		copiedConfig := rest.CopyConfig(config)
		copiedConfig.Insecure = false
		return client.New(copiedConfig, options)
	})
}

func assertMemberCluster(t *testing.T, cachedCluster *CachedToolchainCluster, status toolchainv1alpha1.ToolchainClusterStatus) {
	assert.Equal(t, status, *cachedCluster.ClusterStatus)
	assert.Equal(t, test.NameMember, cachedCluster.OwnerClusterName)
	assert.Equal(t, "http://cluster.com", cachedCluster.APIEndpoint)
}
