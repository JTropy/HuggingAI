package controller

import (
	"testing"

	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/oauth"
)

func TestRegistrationPolicyAllowsOnlyBaseQCustomOAuth(t *testing.T) {
	baseQProvider := oauth.NewGenericOAuthProvider(&model.CustomOAuthProvider{Slug: "baseq"})
	if !isBaseQRegistrationProvider(baseQProvider) {
		t.Fatal("baseQ custom OAuth provider should be allowed for registration")
	}

	otherProvider := oauth.NewGenericOAuthProvider(&model.CustomOAuthProvider{Slug: "other"})
	if isBaseQRegistrationProvider(otherProvider) {
		t.Fatal("non-baseQ custom OAuth provider should not be allowed for registration")
	}

	if isBaseQRegistrationProvider(&oauth.GitHubProvider{}) {
		t.Fatal("built-in OAuth providers should not be allowed for registration")
	}
}

func TestPasswordRegistrationPolicyDisabled(t *testing.T) {
	if isPasswordRegistrationAllowed() {
		t.Fatal("password registration should remain disabled")
	}
}
