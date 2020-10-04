-include Makefile.local

GOBIN=${CURDIR}/bin
GO_PATH ?= $(shell go env GOPATH)
GO_OS ?= $(shell go env GOOS)
GO_ARCH ?= $(shell go env GOARCH)

PKG=$(subst $(GO_PATH)/src/,,$(CURDIR))
GO?=go
GO_PKGS=$(shell go list ./... | grep -v -e '.pb.go')
GO_PKGS_DIR=$(shell go list -f='{{ if or .GoFiles .CgoFiles }}{{ .Dir }}{{ end }}' ./... | grep -v -e '.pb.go' -e 'zz_generated.*.go')
TOOLS_DIR=${CURDIR}/hack/tools
TOOLS_GOBIN=${TOOLS_DIR}/bin
CONTROLLER_GEN=${TOOLS_GOBIN}/controller-gen
KUSTOMIZE=${TOOLS_GOBIN}/kustomize

IMG ?= gcr.io/containerz/kube-timeleap/controller:latest
CRD_OPTIONS ?= "crd:trivialVersions=true"

# ----------------------------------------------------------------------------
# defines

define target
@printf "+ \\x1b[1;32m$(patsubst ,$@,$(1))\\x1b[0m\\n" >&2
endef

# ----------------------------------------------------------------------------
# target

all: manager

mod:
	$(call target)
	@rm -f go.sum
	@${GO} mod tidy
	@${GO} mod vendor

##@ test

fmt: gofumpt gofumports
fmt:  ## Run go fmt against code
	$(call target)
	@-rm -rf ${TOOLS_DIR}/vendor
	@${TOOLS_GOBIN}/gofumpt -s -w -extra ${GO_PKGS_DIR}
	@${TOOLS_GOBIN}/gofumports -w -local=${PKG} ${GO_PKGS_DIR}

vet:  # Run go vet against code
	$(call target)
	@GOOS=linux ${GO} vet $(shell GOOS=linux go list ./... | grep -v 'pkg/vdso')

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test:  ## Run tests
test: generate fmt vet manifests
	$(call target)
	@mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/master/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); ${GO} test -v -race -coverprofile cover.out ./...

##@ build, run

manager: generate manifests fmt vet
manager:  ## Build manager binary
	$(call target)
	@${GO} build -o bin/manager cmd/manager/main.go

run: generate fmt vet manifests
run:  ## Run against the configured Kubernetes cluster in ~/.kube/config
	$(call target)
	@${GO} run ./cmd/manager/main.go

generate: mod controller-gen
generate:  ## Generate code
	$(call target)
	@${CONTROLLER_GEN} object:headerFile="hack/boilerplate/boilerplate.go.txt" paths="./..."

manifests: mod controller-gen
manifests:  ## Generate manifests e.g. CRD, RBAC etc.
	$(call target)
	@$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases


##@ deploy

install: manifests kustomize
install:  ## Install CRDs into a cluster
	$(call target)
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

uninstall: manifests kustomize
uninstall:  ## Uninstall CRDs from a cluster
	$(call target)
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

deploy: manifests kustomize
deploy:  ## Deploy controller in the configured Kubernetes cluster in ~/.kube/config
	$(call target)
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

undeploy:  ## UnDeploy controller from the configured Kubernetes cluster in ~/.kube/config
	$(call target)
	$(KUSTOMIZE) build config/default | kubectl delete -f -


##@ tools

${TOOLS_GOBIN}/%:
	@pushd ${TOOLS_DIR} > /dev/null 2>&1; \
		${GO} mod edit -require=sigs.k8s.io/kind@master -require=sigs.k8s.io/kubebuilder@master > /dev/null 2>&1 || true; \
		rm -f go.sum; ${GO} mod tidy -v
	@pushd ${TOOLS_DIR} > /dev/null 2>&1; \
		./install-tools $*
	@pushd ${TOOLS_DIR} > /dev/null 2>&1; \
		${GO} mod edit -require=sigs.k8s.io/kind@master -require=sigs.k8s.io/kubebuilder@master

.PHONY: tools
tools: ${TOOLS_GOBIN}/''
tools:  ## install tools

tools/%:
	$(call target)
	@${MAKE} ${TOOLS_GOBIN}/$* > /dev/null

gofumpt: ${TOOLS_GOBIN}/gofumpt
gofumports: ${TOOLS_GOBIN}/gofumports
controller-gen: ${TOOLS_GOBIN}/controller-gen
kustomize: ${TOOLS_GOBIN}/kustomize


##@ container

docker/build: test
docker/build:  ## Build the docker image
	$(call target)
	docker image build . -t ${IMG}

docker/push:  ## Push the docker image
	$(call target)
	docker image push ${IMG}

.PHONY: clean
clean:  ## Clean workspace
	$(call target)
	-@rm -rf ./bin ./vendor ${TOOLS_DIR}/bin ${TOOLS_DIR}/vendor

.PHONY: help
help:  ## Show make target help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[33m<target>\033[0m\n"} /^[a-zA-Z_0-9\/_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: env/%
env/%:  ## Print the value of MAKEFILE_VARIABLE. Use `make env/MAKEFILE_VARIABLE`.
	@echo $($*)
