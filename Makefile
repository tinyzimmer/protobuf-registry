GOLANGCI_VERSION = 1.18.0

IMG ?= protobuf-registry:latest
NUM ?= 1
REGISTRY_HOST ?= http://localhost:8080

NUM_PACKAGES = $(shell find pkg -type d | wc -l)

build:
	docker build . -t ${IMG}

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint

bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ./bin v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

license-check:
	go run hack/util.go license-check

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	@bin/golangci-lint run -v --skip-dirs pkg/util/client

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Run tests
test: fmt vet
	go test ./pkg/... -coverprofile cover.out -covermode=atomic -race

# Runs tests and outputs total coverage
coverage: test
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	$(shell go env GOPATH)/bin/goveralls \
		-show \
		-covermode=atomic \
		-coverprofile=cover.out \
		-service=github-actions \
		-repotoken ${COVERALLS_TOKEN} \
		-race \
		-package ./pkg/...

ui/node_modules:
	cd ui && npm install

ui/build: ui/node_modules
	cd ui && npm run build

build-ui: ui/build

run: build
	docker run --rm -p 8080:8080 ${IMG}

run_persistent: build
	mkdir -p data
	docker run \
		--rm \
		-p 8080:8080 \
		-v "`pwd`/data:/opt/proto-registry/data" \
		${IMG} --persist-memory --ui-redirect-all

test_data:
	cd hack && NUM=${NUM} REGISTRY_HOST=${REGISTRY_HOST} bash add_test_data.sh

run_dev_ui: ui/node_modules
	cd ui && npm start

clean:
	rm -rf bin/
	rm -rf data/
	rm -rf ui/build
	rm -rf ui/node_modules
	rm -f cover.out
