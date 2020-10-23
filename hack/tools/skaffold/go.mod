module github.com/kouzoh/kube-timeleap/hack/tools

go 1.15

require github.com/GoogleContainerTools/skaffold master

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.3.4 // github.com/GoogleContainerTools/skaffold@v1.15.0
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200221181110-62bd5a33f707 // github.com/GoogleContainerTools/skaffold@v1.15.0
	golang.org/x/sys => golang.org/x/sys v0.0.0-20191128015809-6d18c012aee9 // github.com/docker/docker@v1.4.2-0.20200221181110-62bd5a33f707
	k8s.io/api => k8s.io/api v0.17.4 // github.com/GoogleContainerTools/skaffold@v1.15.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.4 // github.com/GoogleContainerTools/skaffold@v1.15.0
	k8s.io/client-go => k8s.io/client-go v0.17.4 // github.com/GoogleContainerTools/skaffold@v1.15.0
	k8s.io/kubectl => k8s.io/kubectl v0.17.4 // github.com/GoogleContainerTools/skaffold@v1.15.0
	k8s.io/kubernetes => k8s.io/kubernetes v1.14.10 // github.com/GoogleContainerTools/skaffold@v1.15.0
)
