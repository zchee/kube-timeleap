// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var timeleaplog = logf.Log.WithName("timeleap-resource")

// SetupWebhookWithManager setup TimeLeap webhook with manager.
func (r *TimeLeap) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:webhookVersions=v1,path=/mutate-timeleap-x-k8s-io-v1alpha1-timeleap,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=timeleap.x-k8s.io,resources=timeleaps,verbs=create;update,versions=v1alpha1,name=mtimeleap.kb.io,sideEffects=NoneOnDryRun

// compile time check whether the TimeLeap implements webhook.Defaulter interface.
var _ webhook.Defaulter = (*TimeLeap)(nil)

// Default implements webhook.Defaulter.
//
// Default is a webhook will be registered for the type.
func (r *TimeLeap) Default() {
	timeleaplog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// +kubebuilder:webhook:webhookVersions=v1,verbs=create;update,path=/validate-timeleap-x-k8s-io-v1alpha1-timeleap,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=timeleap.x-k8s.io,resources=timeleaps,versions=v1alpha1,name=vtimeleap.kb.io,sideEffects=None

// compile time check whether the TimeLeap implements webhook.Validator interface.
var _ webhook.Validator = (*TimeLeap)(nil)

// ValidateCreate implements webhook.Validator.
//
// ValidateCreate is a webhook will be registered for the type.
func (r *TimeLeap) ValidateCreate() error {
	timeleaplog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator.
//
// ValidateUpdate is a webhook will be registered for the type.
func (r *TimeLeap) ValidateUpdate(old runtime.Object) error {
	timeleaplog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator.
//
// ValidateDelete is a webhook will be registered for the type.
func (r *TimeLeap) ValidateDelete() error {
	timeleaplog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
