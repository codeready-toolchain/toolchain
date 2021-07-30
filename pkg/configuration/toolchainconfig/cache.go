package toolchainconfig

import (
	"context"
	"sync"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	commonconfig "github.com/codeready-toolchain/toolchain-common/pkg/configuration"
	errs "github.com/pkg/errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var configCache = &cache{}

var cacheLog = logf.Log.WithName("cache_toolchainconfig")

type cache struct {
	sync.RWMutex
	config  *toolchainv1alpha1.ToolchainConfig
	secrets map[string]map[string]string // map of secret key-value pairs indexed by secret name
}

func (c *cache) set(config *toolchainv1alpha1.ToolchainConfig, secrets map[string]map[string]string) {
	c.Lock()
	defer c.Unlock()
	c.config = config.DeepCopy()
	c.secrets = commonconfig.CopyOf(secrets)
}

func (c *cache) get() (*toolchainv1alpha1.ToolchainConfig, map[string]map[string]string) {
	c.RLock()
	defer c.RUnlock()
	return c.config.DeepCopy(), commonconfig.CopyOf(c.secrets)
}

func updateConfig(config *toolchainv1alpha1.ToolchainConfig, secrets map[string]map[string]string) {
	configCache.set(config, secrets)
}

func LoadLatest(cl client.Client) (ToolchainConfig, error) {
	namespace, err := commonconfig.GetWatchNamespace()
	if err != nil {
		return ToolchainConfig{cfg: &toolchainv1alpha1.ToolchainConfigSpec{}}, errs.Wrap(err, "Failed to get watch namespace")
	}

	config := &toolchainv1alpha1.ToolchainConfig{}
	if err := cl.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: "config"}, config); err != nil {
		if apierrors.IsNotFound(err) {
			cacheLog.Info("ToolchainConfig resource with the name 'config' wasn't found, default configuration will be used", "namespace", namespace)
			return ToolchainConfig{cfg: &toolchainv1alpha1.ToolchainConfigSpec{}}, nil
		}
		return ToolchainConfig{cfg: &toolchainv1alpha1.ToolchainConfigSpec{}}, err
	}

	allSecrets, err := commonconfig.LoadSecrets(cl, namespace)
	if err != nil {
		return ToolchainConfig{cfg: &toolchainv1alpha1.ToolchainConfigSpec{}}, err
	}

	configCache.set(config, allSecrets)
	return getConfigOrDefault(), nil
}

// GetConfig returns a cached toolchain config.
// If no config is stored in the cache, then it retrieves it from the cluster and stores in the cache.
// If the resource is not found, then returns the default config.
// If any failure happens while getting the ToolchainConfig resource, then returns an error.
func GetConfig(cl client.Client) (ToolchainConfig, error) {
	config, _ := configCache.get()
	if config == nil {
		return LoadLatest(cl)
	}
	return getConfigOrDefault(), nil
}

func getConfigOrDefault() ToolchainConfig {
	config, secrets := configCache.get()
	if config == nil {
		return ToolchainConfig{cfg: &toolchainv1alpha1.ToolchainConfigSpec{}, secrets: secrets}
	}
	return ToolchainConfig{cfg: &config.Spec, secrets: secrets}
}

// GetCachedConfig returns the cached toolchain config or a toolchainconfig with default values
func GetCachedConfig() ToolchainConfig {
	config, secrets := configCache.get()
	if config == nil {
		return ToolchainConfig{cfg: &toolchainv1alpha1.ToolchainConfigSpec{}, secrets: secrets}
	}
	return ToolchainConfig{cfg: &config.Spec, secrets: secrets}
}

// Reset resets the cache.
// Should be used only in tests, but since it has to be used in other packages,
// then the function has to be exported and placed here.
func Reset() {
	configCache = &cache{}
}
