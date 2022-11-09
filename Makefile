.DEFAULT_GOAL := test
LINT_VERSION = v1.45.2
.PHONY: download
download:
	@echo Download go.mod dependencies
	@go mod download -x

.PHONY: install-tools
install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: code-format-checks
code-format-check:
	@unformatted_files="$$(gofmt -l .)" \
	&& test -z "$$unformatted_files" || ( printf "Unformatted files: \n$${unformatted_files}\nRun make code-format\n"; exit 1 )

lint:
	golangci-lint run --build-tags=musl --config golangci.yml --verbose

lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(LINT_VERSION)
.PHONY: code-format
code-format:
	goimports -l -w .
	gofmt -l -w .

.PHONY: test
test:
	go test ./... -short -count=1   -tags musl -timeout 300ms

.PHONY: docker-compose
docker-compose:
	docker-compose build && docker-compose up

.PHONY: docker-build

.PHONY: docker-run
docker-run:
	docker run -it -p 8000:8000 --rm --name=risk-rules risk-rules

dist/%: cmd/%/main/main.go $(shell find ./internal -name '*.go') ./dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ cmd/$(@F)/main/main.go

integration-tests:
	go test ./test/integration -v -coverpkg=./... -coverprofile integration.out -count=1 -tags musl

acceptance-tests:
	go test ./test/acceptance -v -coverpkg=./... -coverprofile acceptance.out -tags musl

build-image:
	DOCKER_BUILDKIT=1 docker build -t risk-rules -f Dockerfile --ssh default  .
	