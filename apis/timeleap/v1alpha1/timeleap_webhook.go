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

func (r *TimeLeap) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:webhookVersions=v1beta1,path=/mutate-timeleap-x-k8s-io-v1alpha1-timeleap,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=timeleap.x-k8s.io,resources=timeleaps,verbs=create;update,versions=v1alpha1,name=mtimeleap.kb.io,sideEffects=NoneOnDryRun

var _ webhook.Defaulter = &TimeLeap{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *TimeLeap) Default() {
	timeleaplog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:webhookVersions=v1beta1,verbs=create;update,path=/validate-timeleap-x-k8s-io-v1alpha1-timeleap,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=timeleap.x-k8s.io,resources=timeleaps,versions=v1alpha1,name=vtimeleap.kb.io,sideEffects=None

var _ webhook.Validator = &TimeLeap{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *TimeLeap) ValidateCreate() error {
	timeleaplog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *TimeLeap) ValidateUpdate(old runtime.Object) error {
	timeleaplog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *TimeLeap) ValidateDelete() error {
	timeleaplog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
