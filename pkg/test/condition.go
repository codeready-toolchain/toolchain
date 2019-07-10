package test

import (
	"fmt"
	"testing"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/pkg/apis/toolchain/v1alpha1"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

// AssertConditionsMatch asserts that the specified list A of conditions is equal to specified
// list B of conditions ignoring the order of the elements. We can't use assert.ElementsMatch
// because the LastTransitionTime of the actual conditions can be modified but the conditions
// still should be treated as matched
func AssertConditionsMatch(t *testing.T, actual []toolchainv1alpha1.Condition, expected ...toolchainv1alpha1.Condition) {
	require.Equal(t, len(expected), len(actual))
	for _, c := range expected {
		AssertContainsCondition(t, actual, c)
	}
}

// AssertContainsCondition asserts that the specified list of conditions contains the specified condition.
// LastTransitionTime is ignored.
func AssertContainsCondition(t *testing.T, conditions []toolchainv1alpha1.Condition, contains toolchainv1alpha1.Condition) {
	for _, c := range conditions {
		if c.Type == contains.Type {
			assert.Equal(t, contains.Status, c.Status)
			assert.Equal(t, contains.Reason, c.Reason)
			assert.Equal(t, contains.Message, c.Message)
			return
		}
	}
	assert.FailNow(t, fmt.Sprintf("the list of conditions %v doesn't contain the expected condition %v", conditions, contains))
}
