package toolchainconfig

import (
	"strings"
	"time"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ToolchainStatusName = "toolchain-status"

	// NotificationDeliveryServiceMailgun is the notification delivery service to use during production
	NotificationDeliveryServiceMailgun = "mailgun"
)

var logger = logf.Log.WithName("configuration")

type ToolchainConfig struct {
	cfg     *toolchainv1alpha1.ToolchainConfigSpec
	secrets map[string]map[string]string
}

func NewToolchainConfig(cfg *toolchainv1alpha1.ToolchainConfigSpec, secrets map[string]map[string]string) ToolchainConfig {
	return ToolchainConfig{
		cfg:     cfg,
		secrets: secrets,
	}
}

func (c *ToolchainConfig) Print() {
	logger.Info("Toolchain configuration variables", "ToolchainConfigSpec", c.cfg)
}

func (c *ToolchainConfig) Environment() string {
	return getString(c.cfg.Host.Environment, "prod")
}

func (c *ToolchainConfig) AutomaticApproval() AutoApprovalConfig {
	return AutoApprovalConfig{c.cfg.Host.AutomaticApproval}
}

func (c *ToolchainConfig) Deactivation() DeactivationConfig {
	return DeactivationConfig{c.cfg.Host.Deactivation}
}

func (c *ToolchainConfig) Metrics() MetricsConfig {
	return MetricsConfig{c.cfg.Host.Metrics}
}

func (c *ToolchainConfig) Notifications() NotificationsConfig {
	return NotificationsConfig{
		c:       c.cfg.Host.Notifications,
		secrets: c.secrets,
	}
}

func (c *ToolchainConfig) RegistrationService() RegistrationServiceConfig {
	return RegistrationServiceConfig{
		c:       c.cfg.Host.RegistrationService,
		secrets: c.secrets,
	}
}

func (c *ToolchainConfig) Tiers() TiersConfig {
	return TiersConfig{c.cfg.Host.Tiers}
}

func (c *ToolchainConfig) ToolchainStatus() ToolchainStatusConfig {
	return ToolchainStatusConfig{c.cfg.Host.ToolchainStatus}
}

func (c *ToolchainConfig) Users() UsersConfig {
	return UsersConfig{c.cfg.Host.Users}
}

type AutoApprovalConfig struct {
	approval toolchainv1alpha1.AutomaticApprovalConfig
}

func (a AutoApprovalConfig) IsEnabled() bool {
	return getBool(a.approval.Enabled, false)
}

func (a AutoApprovalConfig) ResourceCapacityThresholdDefault() int {
	return getInt(a.approval.ResourceCapacityThreshold.DefaultThreshold, 80)
}

func (a AutoApprovalConfig) ResourceCapacityThresholdSpecificPerMemberCluster() map[string]int {
	return a.approval.ResourceCapacityThreshold.SpecificPerMemberCluster
}

func (a AutoApprovalConfig) MaxNumberOfUsersOverall() int {
	return getInt(a.approval.MaxNumberOfUsers.Overall, 1000)
}

func (a AutoApprovalConfig) MaxNumberOfUsersSpecificPerMemberCluster() map[string]int {
	return a.approval.MaxNumberOfUsers.SpecificPerMemberCluster
}

type DeactivationConfig struct {
	dctv toolchainv1alpha1.DeactivationConfig
}

func (d DeactivationConfig) DeactivatingNotificationDays() int {
	return getInt(d.dctv.DeactivatingNotificationDays, 3)
}

func (d DeactivationConfig) DeactivationDomainsExcluded() []string {
	excluded := getString(d.dctv.DeactivationDomainsExcluded, "")
	v := strings.FieldsFunc(excluded, func(c rune) bool {
		return c == ','
	})
	return v
}

func (d DeactivationConfig) UserSignupDeactivatedRetentionDays() int {
	return getInt(d.dctv.UserSignupDeactivatedRetentionDays, 365)
}

func (d DeactivationConfig) UserSignupUnverifiedRetentionDays() int {
	return getInt(d.dctv.UserSignupUnverifiedRetentionDays, 7)
}

type MetricsConfig struct {
	metrics toolchainv1alpha1.MetricsConfig
}

func (d MetricsConfig) ForceSynchronization() bool {
	return getBool(d.metrics.ForceSynchronization, false)
}

type NotificationsConfig struct {
	c       toolchainv1alpha1.NotificationsConfig
	secrets map[string]map[string]string
}

func (n NotificationsConfig) notificationSecret(secretKey string) string {
	secret := getString(n.c.Secret.Ref, "")
	return n.secrets[secret][secretKey]
}

func (n NotificationsConfig) NotificationDeliveryService() string {
	return getString(n.c.NotificationDeliveryService, "mailgun")
}

func (n NotificationsConfig) DurationBeforeNotificationDeletion() time.Duration {
	v := getString(n.c.DurationBeforeNotificationDeletion, "24h")
	duration, err := time.ParseDuration(v)
	if err != nil {
		duration = 24 * time.Hour
	}
	return duration
}

func (n NotificationsConfig) AdminEmail() string {
	return getString(n.c.AdminEmail, "")
}

func (n NotificationsConfig) MailgunDomain() string {
	key := getString(n.c.Secret.MailgunDomain, "")
	return n.notificationSecret(key)
}

func (n NotificationsConfig) MailgunAPIKey() string {
	key := getString(n.c.Secret.MailgunAPIKey, "")
	return n.notificationSecret(key)
}

func (n NotificationsConfig) MailgunSenderEmail() string {
	key := getString(n.c.Secret.MailgunSenderEmail, "")
	return n.notificationSecret(key)
}

func (n NotificationsConfig) MailgunReplyToEmail() string {
	key := getString(n.c.Secret.MailgunReplyToEmail, "")
	return n.notificationSecret(key)
}

type RegistrationServiceConfig struct {
	c       toolchainv1alpha1.RegistrationServiceConfig
	secrets map[string]map[string]string
}

func (r RegistrationServiceConfig) Analytics() RegistrationServiceAnalyticsConfig {
	return RegistrationServiceAnalyticsConfig{r.c.Analytics}
}

func (r RegistrationServiceConfig) Auth() RegistrationServiceAuthConfig {
	return RegistrationServiceAuthConfig{r.c.Auth}
}

func (r RegistrationServiceConfig) Environment() string {
	return getString(r.c.Environment, "prod")
}

func (r RegistrationServiceConfig) LogLevel() string {
	return getString(r.c.LogLevel, "info")
}

func (r RegistrationServiceConfig) Namespace() string {
	return getString(r.c.Namespace, "toolchain-host-operator")
}

func (r RegistrationServiceConfig) RegistrationServiceURL() string {
	return getString(r.c.RegistrationServiceURL, "https://registration.crt-placeholder.com")
}

func (r RegistrationServiceConfig) Verification() RegistrationServiceVerificationConfig {
	return RegistrationServiceVerificationConfig{c: r.c.Verification, secrets: r.secrets}
}

type RegistrationServiceAnalyticsConfig struct {
	c toolchainv1alpha1.RegistrationServiceAnalyticsConfig
}

func (r RegistrationServiceAnalyticsConfig) WoopraDomain() string {
	return getString(r.c.WoopraDomain, "")
}

func (r RegistrationServiceAnalyticsConfig) SegmentWriteKey() string {
	return getString(r.c.SegmentWriteKey, "")
}

type RegistrationServiceAuthConfig struct {
	c toolchainv1alpha1.RegistrationServiceAuthConfig
}

func (r RegistrationServiceAuthConfig) AuthClientLibraryURL() string {
	return getString(r.c.AuthClientLibraryURL, "https://sso.prod-preview.openshift.io/auth/js/keycloak.js")
}

func (r RegistrationServiceAuthConfig) AuthClientConfigContentType() string {
	return getString(r.c.AuthClientConfigContentType, "application/json; charset=utf-8")
}

func (r RegistrationServiceAuthConfig) AuthClientConfigRaw() string {
	return getString(r.c.AuthClientConfigRaw, `{"realm": "toolchain-public","auth-server-url": "https://sso.prod-preview.openshift.io/auth","ssl-required": "none","resource": "crt","clientId": "crt","public-client": true}`)
}

func (r RegistrationServiceAuthConfig) AuthClientPublicKeysURL() string {
	return getString(r.c.AuthClientPublicKeysURL, "https://sso.prod-preview.openshift.io/auth/realms/toolchain-public/protocol/openid-connect/certs")
}

type RegistrationServiceVerificationConfig struct {
	c       toolchainv1alpha1.RegistrationServiceVerificationConfig
	secrets map[string]map[string]string
}

func (r RegistrationServiceVerificationConfig) registrationServiceSecret(secretKey string) string {
	secret := getString(r.c.Secret.Ref, "")
	return r.secrets[secret][secretKey]
}

func (r RegistrationServiceVerificationConfig) Enabled() bool {
	return getBool(r.c.Enabled, false)
}

func (r RegistrationServiceVerificationConfig) DailyLimit() int {
	return getInt(r.c.DailyLimit, 5)
}

func (r RegistrationServiceVerificationConfig) AttemptsAllowed() int {
	return getInt(r.c.AttemptsAllowed, 3)
}

func (r RegistrationServiceVerificationConfig) MessageTemplate() string {
	return getString(r.c.MessageTemplate, "Developer Sandbox for Red Hat OpenShift: Your verification code is %s")
}

func (r RegistrationServiceVerificationConfig) ExcludedEmailDomains() string {
	return getString(r.c.ExcludedEmailDomains, "")
}

func (r RegistrationServiceVerificationConfig) CodeExpiresInMin() int {
	return getInt(r.c.CodeExpiresInMin, 5)
}

func (r RegistrationServiceVerificationConfig) TwilioAccountSID() string {
	key := getString(r.c.Secret.TwilioAccountSID, "")
	return r.registrationServiceSecret(key)
}

func (r RegistrationServiceVerificationConfig) TwilioAuthToken() string {
	key := getString(r.c.Secret.TwilioAuthToken, "")
	return r.registrationServiceSecret(key)
}

func (r RegistrationServiceVerificationConfig) TwilioFromNumber() string {
	key := getString(r.c.Secret.TwilioFromNumber, "")
	return r.registrationServiceSecret(key)
}

type TiersConfig struct {
	tiers toolchainv1alpha1.TiersConfig
}

func (d TiersConfig) DefaultTier() string {
	return getString(d.tiers.DefaultTier, "base")
}

func (d TiersConfig) DurationBeforeChangeTierRequestDeletion() time.Duration {
	v := getString(d.tiers.DurationBeforeChangeTierRequestDeletion, "24h")
	duration, err := time.ParseDuration(v)
	if err != nil {
		duration = 24 * time.Hour
	}
	return duration
}

func (d TiersConfig) TemplateUpdateRequestMaxPoolSize() int {
	return getInt(d.tiers.TemplateUpdateRequestMaxPoolSize, 5)
}

type ToolchainStatusConfig struct {
	t toolchainv1alpha1.ToolchainStatusConfig
}

func (d ToolchainStatusConfig) ToolchainStatusRefreshTime() time.Duration {
	v := getString(d.t.ToolchainStatusRefreshTime, "5s")
	duration, err := time.ParseDuration(v)
	if err != nil {
		duration = 5 * time.Second
	}
	return duration
}

type UsersConfig struct {
	c toolchainv1alpha1.UsersConfig
}

func (d UsersConfig) MasterUserRecordUpdateFailureThreshold() int {
	return getInt(d.c.MasterUserRecordUpdateFailureThreshold, 2) // default: allow 1 failure, try again and then give up if failed again
}

func (d UsersConfig) ForbiddenUsernamePrefixes() []string {
	prefixes := getString(d.c.ForbiddenUsernamePrefixes, "openshift,kube,default,redhat,sandbox")
	v := strings.FieldsFunc(prefixes, func(c rune) bool {
		return c == ','
	})
	return v
}

func (d UsersConfig) ForbiddenUsernameSuffixes() []string {
	suffixes := getString(d.c.ForbiddenUsernameSuffixes, "admin")
	v := strings.FieldsFunc(suffixes, func(c rune) bool {
		return c == ','
	})
	return v
}

func getBool(value *bool, defaultValue bool) bool {
	if value != nil {
		return *value
	}
	return defaultValue
}

func getInt(value *int, defaultValue int) int {
	if value != nil {
		return *value
	}
	return defaultValue
}

func getString(value *string, defaultValue string) string {
	if value != nil {
		return *value
	}
	return defaultValue
}
