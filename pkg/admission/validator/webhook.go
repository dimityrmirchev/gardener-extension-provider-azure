// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	extensionspredicate "github.com/gardener/gardener/extensions/pkg/predicate"
	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	"github.com/gardener/gardener/pkg/apis/security"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/gardener/gardener-extension-provider-azure/pkg/azure"
)

const (
	// Name is a name for a validation webhook.
	Name = "validator"
	// SecretsValidatorName is the name of the secrets validator.
	SecretsValidatorName = "secrets." + Name
)

var logger = log.Log.WithName("azure-validator-webhook")

// New creates a new webhook that validates Shoot, CloudProfile, SecretBinding and CredentialsBinding resources.
func New(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	logger.Info("Setting up webhook", "name", Name)

	return extensionswebhook.New(mgr, extensionswebhook.Args{
		Provider: azure.Type,
		Name:     Name,
		Path:     "/webhooks/validate",
		// TODO(dimityrmirchev): Uncomment this line once this extension uses a g/g version that contains https://github.com/gardener/gardener/pull/10499
		// Predicates: []predicate.Predicate{predicate.Or(extensionspredicate.GardenCoreProviderType(azure.Type), extensionspredicate.GardenSecurityProviderType(azure.Type))},
		Predicates: []predicate.Predicate{extensionspredicate.GardenCoreProviderType(azure.Type)},
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewShootValidator(mgr):              {{Obj: &core.Shoot{}}},
			NewCloudProfileValidator(mgr):       {{Obj: &core.CloudProfile{}}},
			NewSecretBindingValidator(mgr):      {{Obj: &core.SecretBinding{}}},
			NewCredentialsBindingValidator(mgr): {{Obj: &security.CredentialsBinding{}}},
		},
		Target: extensionswebhook.TargetSeed,
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"provider.extensions.gardener.cloud/azure": "true"},
		},
	})
}

// NewSecretsWebhook creates a new validation webhook for Secrets.
func NewSecretsWebhook(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	logger.Info("Setting up webhook", "name", SecretsValidatorName)

	return extensionswebhook.New(mgr, extensionswebhook.Args{
		Provider: azure.Type,
		Name:     SecretsValidatorName,
		Path:     "/webhooks/validate/secrets",
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewSecretValidator(): {{Obj: &corev1.Secret{}}},
		},
		Target: extensionswebhook.TargetSeed,
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"provider.shoot.gardener.cloud/azure": "true"},
		},
	})
}
