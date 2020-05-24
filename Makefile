EXECUTABLE=metal-lookup
PROJECT=github.com/mandelsoft/metal-lookup
VERSION=$(shell cat VERSION)

RELEASE                     := true
NAME                        := metal-lookup
CMD                         := server
REPOSITORY                  := github.com/mandelsoft/metal-lookup
#REGISTRY                    := eu.gcr.io/gardener-project
REGISTRY                    :=
IMAGEORG                    := mandelsoft
IMAGE_PREFIX                := $(REGISTRY)$(IMAGEORG)
REPO_ROOT                   := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
HACK_DIR                    := $(REPO_ROOT)/hack
VERSION                     := $(shell cat "$(REPO_ROOT)/VERSION")
COMMIT                      := $(shell git rev-parse HEAD)
ifeq ($(RELEASE),true)
IMAGE_TAG                   := $(VERSION)
else
IMAGE_TAG                   := $(VERSION)-$(COMMIT)
endif
VERSION_VAR                 := main.Version
.PHONY: all
ifeq ($(RELEASE),true)
all: generate release
else
all: generate dev
endif


.PHONY: revendor
revendor:
	@GO111MODULE=on go mod vendor
	@GO111MODULE=on go mod tidy


.PHONY: check
check:
	@.ci/check

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o $(EXECUTABLE) \
        -mod=vendor \
	    -ldflags "-X $(VERSION_VAR)=$(VERSION)-$(COMMIT)" \
	    ./cmd/$(CMD)

.PHONY: dev
dev:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go install \
        -mod=vendor \
	    -ldflags "-X $(VERSION_VAR)=$(VERSION)-$(COMMIT)" \
	    ./cmd/$(CMD)

.PHONY: build-local
build-local:
	CGO_ENABLED=0 GO111MODULE=on go build -o $(EXECUTABLE) \
	    -mod=vendor \
	    -ldflags "-X $(VERSION_VAR)=$(VERSION)-$(COMMIT)" \
	    ./cmd/$(CMD)

.PHONY: release-all
release-all: generate release

.PHONY: release
release:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go install \
	    -a \
	    -mod=vendor \
	    -ldflags "-w -X $(VERSION_VAR)=$(VERSION)" \
	    ./cmd/$(CMD)

.PHONY: test
test:
	GO111MODULE=on go test -mod=vendor ./pkg/...

.PHONY: generate
generate:
	@go generate ./pkg/...


### Docker commands

.PHONY: docker-login
docker-login:
	@gcloud auth activate-service-account --key-file .kube-secrets/gcr/gcr-readwrite.json

.PHONY: images-dev
images-dev:
	@docker build -t $(IMAGE_PREFIX)/$(NAME):$(VERSION)-$(COMMIT) -t $(IMAGE_PREFIX)/$(NAME):latest -f Dockerfile -m 6g --build-arg TARGETS=dev --target main .
	@docker push $(IMAGE_PREFIX)/$(NAME):latest

.PHONY: images-release
images-release:
	@docker build -t $(IMAGE_PREFIX)/$(NAME):$(VERSION) -t $(IMAGE_PREFIX)/$(NAME):latest -f Dockerfile -m 6g --build-arg TARGETS=release --target main .

.PHONY: images-release-all
images-release-all:
	@docker build -t $(IMAGE_PREFIX)/$(NAME):$(VERSION) -t $(IMAGE_PREFIX)/$(NAME):latest -f Dockerfile -m 6g --build-arg TARGETS=release-all --target main .
