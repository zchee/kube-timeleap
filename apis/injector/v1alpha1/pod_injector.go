// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

// Package v1alpha1 contains API Schema definitions for the injector v1alpha1 API group.
package v1alpha1

import (
	"context"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	// Path path of admission webhook for injection.
	Path = "/inject-v1-pod"
)

var log = logf.Log.WithName("injector-resource")

// +kubebuilder:webhook:webhookVersions=v1beta1,verbs=create;update,path=/inject-v1-pod,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups="",resources=pods,versions=v1,name=ipod.kb.io,sideEffects=NoneOnDryRun

// Pod represents a injecting pod.
type Pod struct {
	Client  client.Client
	decoder *admission.Decoder
}

// compile time check whether the Pod implements admission.DecoderInjector and admission.Handler interfaces.
var (
	_ admission.DecoderInjector = (*Pod)(nil)
	_ admission.Handler         = (*Pod)(nil)
)

// InjectDecoder implements admission.DecoderInjector.
func (r *Pod) InjectDecoder(d *admission.Decoder) error {
	r.decoder = d
	return nil
}

// Handle implements admission.Handler.
func (r *Pod) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := r.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	log.Info("Get request from pod", "pod", pod)

	return admission.Response{}
}
