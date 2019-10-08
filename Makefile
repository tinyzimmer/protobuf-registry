GOLANGCI_VERSION = 1.18.0

IMG ?= protobuf-registry:latest
NUM ?= 10

build:
	docker build . -t ${IMG}

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint

bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ./bin v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

license-check:
	go run hack/licensecheck.go

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	@bin/golangci-lint run -v

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

clean:
	rm -rf bin/
	rm -rf data/
	rm -rf ui/build
	rm -rf ui/node_modules
	rm -f cover.out

ui-deps:
	cd ui && npm install

build-ui: ui-deps
	cd ui && npm run build

run: build
	docker run --rm -p 8080:8080 ${IMG}

run_persistent: build
	mkdir -p data
	docker run \
		--rm \
		-p 8080:8080 \
		-v "`pwd`/data:/data" \
		-e PROTO_REGISTRY_PERSIST_MEMORY=true \
		-e PROTO_REGISTRY_ENABLE_CORS=true \
		${IMG}

test_data:
	NUM=${NUM} bash hack/add_test_data.sh

run_dev_ui: ui-deps
	cd ui && npm start
