// +build tools

package tools

import (
	_ "k8s.io/code-generator/cmd/conversion-gen"
	_ "mvdan.cc/gofumpt"
	_ "mvdan.cc/gofumpt/gofumports"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
	_ "sigs.k8s.io/kind"
	_ "sigs.k8s.io/kubebuilder/cmd"
	_ "sigs.k8s.io/kustomize/kustomize/v3"
)
