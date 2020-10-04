module github.com/kouzoh/kube-timeleap/hack/tools

go 1.15

require (
	k8s.io/code-generator v0.19.2
	mvdan.cc/gofumpt v0.0.0-20200927160801-5bfeb2e70dd6
	sigs.k8s.io/controller-tools v0.4.0
	sigs.k8s.io/kind master
	sigs.k8s.io/kubebuilder master
	sigs.k8s.io/kustomize/kustomize/v3 v3.8.4
)

replace (
	golang.org/x/sys => golang.org/x/sys v0.0.0-20191128015809-6d18c012aee9 // sigs.k8s.io/kustomize/kustomize/v3
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.3 // sigs.k8s.io/kustomize/kustomize/v3
	k8s.io/client-go => k8s.io/client-go v0.17.3 // sigs.k8s.io/kustomize/kustomize/v3
)
