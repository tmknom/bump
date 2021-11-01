# If this variable is not set, the program /bin/sh is used as the shell.
#
# The -e flag causes bash with qualifications to exit immediately if a command it executes fails.
# The -u flag causes bash to exit with an error message if a variable is accessed without being defined.
# The -o pipefail option causes bash to exit if any of the commands in a pipeline fail.
# The -c flag is  how make passes the script to be executed to bash.
SHELL := /bin/bash -eu -o pipefail -c

# This option causes make to display a warning whenever an undefined variable is expanded.
MAKEFLAGS += --warn-undefined-variables

# Disable any builtin rules and variables, so speedup a bit.
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables
.SUFFIXES:

# Sets the default goal to be used if no targets were specified on the command line.
.DEFAULT_GOAL := help

#
# Variables to be used by commands
#
GIT ?= $(shell which git)
DOCKER ?= $(shell which docker)
DOCKER_RUN ?= $(DOCKER) run -i --rm -v $(CURDIR):/work -w /work

#
# Build and run
#
.PHONY: build
build: ## build executable binary
	go build -o bin/bump ./cmd/bump

.PHONY: run
run: build ## run bump command for test
	bin/bump init
	bin/bump show
	bin/bump patch
	bin/bump minor
	bin/bump major
	@$(MAKE) clean

#
# Development
#
.PHONY: install-tools
install-tools: ## install tools for development
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: check
check: lint test ## check all

.PHONY: test
test: ## test
	go test ./...
	rm -rf ./VERSION

.PHONY: lint
lint: goimports vet staticcheck ## lint all

.PHONY: staticcheck
staticcheck: ## static analysis by staticcheck
	staticcheck ./...

.PHONY: vet
vet: ## static analysis by vet
	go vet ./...

.PHONY: goimports
goimports: ## update import lines
	goimports -w .

.PHONY: clean
clean: ## clean files
	rm -rf ./bin
	rm -rf ./VERSION

#
# Docs
#
.PHONY: format-markdown
format-markdown: ## format markdown by prettier
	$(DOCKER_RUN) prettier --write --parser=markdown **/*.md

#
# Release
#
.PHONY: release-major
release-major: ## release major version
	$(call release,major)

.PHONY: release-minor
release-minor: ## release minor version
	$(call release,minor)

.PHONY: release-patch
release-patch: ## release patch version
	$(call release,patch)

.PHONY: release-initial
release-initial: ## release initial version
	$(call release,init)

define release
	$(BUMP) $(1)
	$(MAKE) generate-changelog
	$(MAKE) git-push
endef

.PHONY: generate-changelog
generate-changelog:
	$(STANDARD_VERSION) --release-as "$(CURRENT_VERSION)"
	$(MAKE) format-markdown

.PHONY: git-push
git-push:
	$(GIT) switch -c "$(RELEASE_BRANCH)"
	$(GIT) add CHANGELOG.md VERSION
	$(GIT) commit -m "chore: release $(CURRENT_VERSION)"
	$(GIT) tag -a $(RELEASE_TAG) -m "release $(CURRENT_VERSION)"
	$(GIT) push origin "$(RELEASE_BRANCH)" "$(RELEASE_TAG)"

BUMP ?= $(shell which bump)
STANDARD_VERSION ?= $(DOCKER_RUN) standard-version --skip.commit --skip.tag
CURRENT_VERSION ?= $(shell cat VERSION)
RELEASE_BRANCH ?= release-$(CURRENT_VERSION)
RELEASE_TAG ?= v$(CURRENT_VERSION)

#
# CI
#
.PHONY: publish-release-notes
publish-release-notes: ## publish release notes
	grep -A 100000 -E '## .+[0-9]' CHANGELOG.md | sed '1d' | sed '1d' > $${TMPDIR}/tmp.md
	grep -m1 -B 100000 -E '## .+[0-9]' $${TMPDIR}/tmp.md | sed '$$d' | sed '$$d' > $${TMPDIR}/release.md \
	  || cp $${TMPDIR}/tmp.md $${TMPDIR}/release.md
	gh release create $(RELEASE_TAG) -F $${TMPDIR}/release.md

#
# Help
#
.PHONY: help
help: ## show help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
