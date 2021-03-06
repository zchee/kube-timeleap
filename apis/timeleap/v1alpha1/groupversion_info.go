// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

// Package v1alpha1 contains API Schema definitions for the timeleap v1alpha1 API group.

// +kubebuilder:object:generate=true
// +groupName=timeleap.x-k8s.io
package v1alpha1

import (
	runtimeschema "k8s.io/apimachinery/pkg/runtime/schema"
	ctrlscheme "sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects.
	GroupVersion = runtimeschema.GroupVersion{Group: "timeleap.x-k8s.io", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme.
	SchemeBuilder = &ctrlscheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
