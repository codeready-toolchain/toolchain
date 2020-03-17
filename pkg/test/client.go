package test

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/kubefed/pkg/apis"
)

// NewFakeClient creates a fake K8s client with ability to override specific Get/List/Create/Update/StatusUpdate/Delete functions
func NewFakeClient(t *testing.T, initObjs ...runtime.Object) *FakeClient {
	s := scheme.Scheme
	err := apis.AddToScheme(s)
	require.NoError(t, err)
	client := fake.NewFakeClientWithScheme(s, initObjs...)
	return &FakeClient{Client: client, T: t}
}

type FakeClient struct {
	client.Client
	T                *testing.T
	MockGet          func(ctx context.Context, key client.ObjectKey, obj runtime.Object) error
	MockList         func(ctx context.Context, list runtime.Object, opts ...client.ListOption) error
	MockCreate       func(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error
	MockUpdate       func(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error
	MockPatch        func(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error
	MockStatusUpdate func(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error
	MockStatusPatch  func(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error
	MockDelete       func(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error
	MockDeleteAllOf  func(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error
}

type mockStatusUpdate struct {
	mockUpdate func(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error
	mockPatch  func(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error
}

func (m *mockStatusUpdate) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	return m.mockUpdate(ctx, obj, opts...)
}

func (m *mockStatusUpdate) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	return m.mockPatch(ctx, obj, patch, opts...)
}

func (c *FakeClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	if c.MockGet != nil {
		return c.MockGet(ctx, key, obj)
	}
	return c.Client.Get(ctx, key, obj)
}

func (c *FakeClient) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	if c.MockList != nil {
		return c.MockList(ctx, list, opts...)
	}
	return c.Client.List(ctx, list, opts...)
}

func (c *FakeClient) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	if c.MockCreate != nil {
		return c.MockCreate(ctx, obj, opts...)
	}

	// Set Generation to `1` for newly created objects since the kube fake client doesn't set it
	mt, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	mt.SetGeneration(1)
	return c.Client.Create(ctx, obj, opts...)
}

func (c *FakeClient) Status() client.StatusWriter {
	m := mockStatusUpdate{}
	if c.MockStatusUpdate == nil && c.MockStatusPatch == nil {
		return c.Client.Status()
	}
	if c.MockStatusUpdate != nil {
		m.mockUpdate = c.MockStatusUpdate
	}
	if c.MockStatusPatch != nil {
		m.mockPatch = c.MockStatusPatch
	}
	return &m
}

func (c *FakeClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	if c.MockUpdate != nil {
		return c.MockUpdate(ctx, obj, opts...)
	}

	// Update Generation if needed since the kube fake client doesn't update generations.
	// Compare the specs (only) and only increment the generation if something changed
	// (the server will check the object metadata, but we're skipping this here)
	if svc, ok := obj.(*corev1.Service); ok {
		existing := corev1.Service{}
		if err := c.Client.Get(ctx, types.NamespacedName{Namespace: svc.GetNamespace(), Name: svc.GetName()}, &existing); err != nil {
			return err
		}
		if !reflect.DeepEqual(existing.Spec, svc.Spec) { // Service has a `spec` field
			svc.SetGeneration(existing.GetGeneration() + 1)
		}
	} else if cm, ok := obj.(*corev1.ConfigMap); ok {
		existing := corev1.ConfigMap{}
		if err := c.Client.Get(ctx, types.NamespacedName{Namespace: cm.GetNamespace(), Name: cm.GetName()}, &existing); err != nil {
			return err
		}
		if !reflect.DeepEqual(existing.Data, cm.Data) { // ConfigMap has a `data` field
			cm.SetGeneration(existing.GetGeneration() + 1)
		}
	}
	return c.Client.Update(ctx, obj, opts...)
}

func (c *FakeClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	if c.MockDelete != nil {
		return c.MockDelete(ctx, obj, opts...)
	}
	return c.Client.Delete(ctx, obj, opts...)
}

func (c *FakeClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	if c.MockDeleteAllOf != nil {
		return c.MockDeleteAllOf(ctx, obj, opts...)
	}
	return c.Client.DeleteAllOf(ctx, obj, opts...)
}

func (c *FakeClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	if c.MockPatch != nil {
		return c.MockPatch(ctx, obj, patch, opts...)
	}
	return c.Client.Patch(ctx, obj, patch, opts...)
}
