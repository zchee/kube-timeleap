// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(&TimeLeap{}, &TimeLeapList{})
}

// TimeLeapSpec defines the desired state of TimeLeap.
type TimeLeapSpec struct {
	// Foo is an example field of TimeLeap. Edit TimeLeap_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// TimeLeapStatus defines the observed state of TimeLeap.
type TimeLeapStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TimeLeap is the Schema for the timeleaps API.
type TimeLeap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TimeLeapSpec   `json:"spec,omitempty"`
	Status TimeLeapStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TimeLeapList contains a list of TimeLeap.
type TimeLeapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TimeLeap `json:"items"`
}
