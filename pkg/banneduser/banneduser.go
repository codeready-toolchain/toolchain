package banneduser

import (
	"context"
	"fmt"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewBannedUser creates a bannedUser resource
func NewBannedUser(userSignup *toolchainv1alpha1.UserSignup, bannedBy string) (*toolchainv1alpha1.BannedUser, error) {
	var emailHashLbl, phoneHashLbl string
	var exists bool

	if emailHashLbl, exists = userSignup.Labels[toolchainv1alpha1.UserSignupUserEmailHashLabelKey]; !exists {
		return nil, fmt.Errorf("the UserSignup %s doesn't have the label '%s' set", userSignup.Name, toolchainv1alpha1.UserSignupUserEmailHashLabelKey) // nolint:loggercheck
	}

	bannedUser := &toolchainv1alpha1.BannedUser{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: userSignup.Namespace,
			Name:      fmt.Sprintf("banneduser-%s", emailHashLbl),
			Labels: map[string]string{
				toolchainv1alpha1.BannedUserEmailHashLabelKey: emailHashLbl,
				toolchainv1alpha1.BannedByLabelKey:            bannedBy,
			},
		},
		Spec: toolchainv1alpha1.BannedUserSpec{
			Email: userSignup.Spec.IdentityClaims.Email,
		},
	}

	if phoneHashLbl, exists = userSignup.Labels[toolchainv1alpha1.UserSignupUserPhoneHashLabelKey]; exists {
		bannedUser.Labels[toolchainv1alpha1.BannedUserPhoneNumberHashLabelKey] = phoneHashLbl
	}
	return bannedUser, nil
}

// GetBannedUser returns BannedUser with the provided user email hash if found. Otherwise it returns nil.
func GetBannedUser(ctx context.Context, userEmailHash string, hostClient client.Client, hostNamespace string) (*toolchainv1alpha1.BannedUser, error) {
	emailHashLabelMatch := client.MatchingLabels(map[string]string{
		toolchainv1alpha1.BannedUserEmailHashLabelKey: userEmailHash,
	})
	bannedUsers := &toolchainv1alpha1.BannedUserList{}

	if err := hostClient.List(ctx, bannedUsers, emailHashLabelMatch, client.InNamespace(hostNamespace)); err != nil {
		return nil, err
	}

	if len(bannedUsers.Items) > 0 {
		return &bannedUsers.Items[0], nil
	}

	return nil, nil
}
