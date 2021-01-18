APPNAME := s3website-openid-proxy
STAGE ?= dev
BRANCH ?= master
SAR_VERSION ?= 1.0.0
MODULE_PKG := github.com/wolfeidau/s3website-openid-proxy

GOLANGCI_VERSION = 1.31.0

GIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y%m%dT%H%M%S')

DEPLOY_BUCKET_SSM_PATH ?= "/config/$(STAGE)/$(BRANCH)/$(APPNAME)/deploy_bucket"

# This path is used to cache binaries used for development and can be overridden to avoid issues with osx vs linux
# binaries.
BIN_DIR ?= $(shell pwd)/bin

default: clean build archive deploy-bucket package deploy-api
.PHONY: default

deploy: build archive package deploy-api
.PHONY: deploy

ci: clean lint test
.PHONY: ci

LDFLAGS := -ldflags="-s -w -X $(MODULE_PKG)/internal/app.BuildDate=${BUILD_DATE} -X $(MODULE_PKG)/internal/app.Commit=${GIT_HASH}"

$(BIN_DIR)/golangci-lint: $(BIN_DIR)/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} $(BIN_DIR)/golangci-lint
$(BIN_DIR)/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv $(BIN_DIR)/golangci-lint $@

$(BIN_DIR)/mockgen:
	@go get -u github.com/golang/mock/mockgen
	@env GOBIN=$(BIN_DIR) GO111MODULE=on go install github.com/golang/mock/mockgen

mocks: $(BIN_DIR)/mockgen
	@echo "--- build all the mocks"
	@bin/mockgen -destination=mocks/session_store.go -package=mocks github.com/dghubble/sessions Store
.PHONY: mocks

clean:
	@echo "--- clean all the things"
	@rm -rf ./dist
.PHONY: clean

lint: $(BIN_DIR)/golangci-lint
	@echo "--- lint all the things"
	@$(BIN_DIR)/golangci-lint run
.PHONY: lint

lint-fix: $(BIN_DIR)/golangci-lint
	@echo "--- lint all the things"
	@$(BIN_DIR)/golangci-lint run --fix
.PHONY: lint-fix

test:
	@echo "--- test all the things"
	@go test -coverprofile=coverage.txt ./...
	@go tool cover -func=coverage.txt
.PHONY: test

build:
	@echo "--- build all the things"
	@mkdir -p dist
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -trimpath -o dist ./cmd/...
.PHONY: build

archive:
	@echo "--- build an archive"
	@cd dist && zip -X -9 -r ./handler.zip *-lambda
.PHONY: archive

deploy-bucket:
	@sam deploy \
		--no-fail-on-empty-changeset \
		--template-file sam/deploy/bucket.yaml \
		--capabilities CAPABILITY_IAM \
		--stack-name $(APPNAME)-$(STAGE)-$(BRANCH)-deploybucket \
		--tags "environment=$(STAGE)" "branch=$(BRANCH)" "service=$(APPNAME)" \
		--parameter-overrides \
			AppName=$(APPNAME) \
			Stage=$(STAGE) \
			Branch=$(BRANCH)
.PHONY: deploy-bucket

package:
	@echo "--- uploading cloudformation assets to $(S3_BUCKET)"
	@sam package \
		--template-file sam/backend/api.yaml \
		--output-template-file dist/api.out.yaml \
		--s3-bucket $(shell aws ssm get-parameter --name $(DEPLOY_BUCKET_SSM_PATH) --query 'Parameter.Value' --output text) \
		--s3-prefix sam/$(GIT_HASH)
.PHONY: package

publish:
	@echo "--- publish stack $(APPNAME)-$(STAGE)-$(BRANCH)"
	@sam publish \
		--template-file api.out.yaml \
		--semantic-version $(SAR_VERSION)
.PHONY: publish

deploy-api:
	@echo "--- deploy stack $(APPNAME)-$(STAGE)-$(BRANCH)"
	@sam deploy \
		--no-fail-on-empty-changeset \
		--template-file dist/api.out.yaml \
		--capabilities CAPABILITY_IAM \
		--tags "environment=$(STAGE)" "branch=$(BRANCH)" "service=$(APPNAME)" \
		--stack-name $(APPNAME)-$(STAGE)-$(BRANCH) \
		--parameter-overrides AppName=$(APPNAME) Stage=$(STAGE) Branch=$(BRANCH) \
			ClientID=$(CLIENT_ID) ClientSecret=$(CLIENT_SECRET) Issuer=$(ISSUER) \
			HostedZoneId=$(HOSTED_ZONE_ID) HostedZoneName=$(HOSTED_ZONE_NAME) \
			SubDomainName=$(SUBDOMAIN_NAME)
.PHONY: deploy-api
