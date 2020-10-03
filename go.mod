module github.com/zchee/kube-timeleap

go 1.15

require (
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.0-alpha.2
)

// pin
replace (
	k8s.io/api => k8s.io/api v0.19.2 // k8s.io/client-go@v0.19.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.2 // k8s.io/client-go@v0.19.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.2 // k8s.io/client-go@v0.19.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.2 // k8s.io/client-go@v0.19.2
	k8s.io/client-go => k8s.io/client-go v0.19.2 // k8s.io/client-go@v0.19.2
	k8s.io/utils => k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // sigs.k8s.io/controller-runtime@v0.7.0-alpha.2
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.2.0
)
