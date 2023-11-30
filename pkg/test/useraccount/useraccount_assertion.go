package useraccount

import (
	"context"
	"testing"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	"github.com/codeready-toolchain/toolchain-common/pkg/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Assertion struct {
	userAccount    *toolchainv1alpha1.UserAccount
	client         client.Client
	namespacedName types.NamespacedName
	t              *testing.T
}

func (a *Assertion) loadUaAssertion() error {
	if a.userAccount != nil {
		return nil
	}
	ua := &toolchainv1alpha1.UserAccount{}
	err := a.client.Get(context.TODO(), a.namespacedName, ua)
	a.userAccount = ua
	return err
}

func AssertThatUserAccount(t *testing.T, name string, client client.Client) *Assertion {
	return &Assertion{
		client:         client,
		namespacedName: test.NamespacedName(test.MemberOperatorNs, name),
		t:              t,
	}
}

func (a *Assertion) DoesNotExist() *Assertion {
	err := a.loadUaAssertion()
	require.Error(a.t, err)
	assert.IsType(a.t, metav1.StatusReasonNotFound, errors.ReasonForError(err))
	return a
}

func (a *Assertion) Exists() *Assertion {
	err := a.loadUaAssertion()
	require.NoError(a.t, err)
	return a
}

func (a *Assertion) Get() *toolchainv1alpha1.UserAccount {
	err := a.loadUaAssertion()
	require.NoError(a.t, err)
	return a.userAccount
}

func (a *Assertion) HasFinalizer(finalizer string) *Assertion {
	err := a.loadUaAssertion()
	require.NoError(a.t, err)
	assert.Contains(a.t, a.userAccount.Finalizers, finalizer)
	return a
}

func (a *Assertion) HasNoFinalizer() *Assertion {
	err := a.loadUaAssertion()
	require.NoError(a.t, err)
	assert.Empty(a.t, a.userAccount.Finalizers)
	return a
}

func (a *Assertion) MatchMasterUserRecord(mur *toolchainv1alpha1.MasterUserRecord) *Assertion {
	err := a.loadUaAssertion()
	require.NoError(a.t, err)
	assert.Equal(a.t, mur.Spec.PropagatedClaims, a.userAccount.Spec.PropagatedClaims)
	assert.Equal(a.t, mur.Spec.Disabled, a.userAccount.Spec.Disabled)
	assert.NotNil(a.t, a.userAccount.Annotations)
	return a
}

// HasLabelWithValue verifies that the UserAccount has
// a label with the given key and value
func (a *Assertion) HasLabelWithValue(key, value string) *Assertion {
	err := a.loadUaAssertion()
	require.NoError(a.t, err)
	require.NotNil(a.t, a.userAccount.Labels)
	assert.Equal(a.t, value, a.userAccount.Labels[key])
	return a
}

// HasAnnotationWithValue verifies that the UserAccount has an annotation with the given key and value
func (a *Assertion) HasAnnotationWithValue(key, value string) *Assertion {
	err := a.loadUaAssertion()
	require.NoError(a.t, err)
	require.NotNil(a.t, a.userAccount.Annotations)
	assert.Equal(a.t, value, a.userAccount.Annotations[key])
	return a
}

func (a *Assertion) HasSpec(spec toolchainv1alpha1.UserAccountSpec) *Assertion {
	err := a.loadUaAssertion()
	require.NoError(a.t, err)
	assert.Equal(a.t, spec, a.userAccount.Spec)
	return a
}

func (a *Assertion) HasConditions(expected ...toolchainv1alpha1.Condition) *Assertion {
	err := a.loadUaAssertion()
	require.NoError(a.t, err)
	test.AssertConditionsMatch(a.t, a.userAccount.Status.Conditions, expected...)
	return a
}

func (a *Assertion) HasNoConditions() *Assertion {
	err := a.loadUaAssertion()
	require.NoError(a.t, err)
	assert.Empty(a.t, a.userAccount.Status.Conditions)
	return a
}
