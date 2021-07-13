SHELL := /bin/bash -euo pipefail

VERSION=0.1.0

VERSION_PARTS := $(subst ., ,$(VERSION))

MAJOR := $(word 1,$(VERSION_PARTS))
MINOR := $(word 2,$(VERSION_PARTS))
PATCH := $(word 3,$(VERSION_PARTS))

test:
	@go test ./...

lint:
	@golangci-lint run --fast

build:
	@go build -mod=readonly ./...

release-patch: guard
	@make release PATCH=$$(( $(PATCH) + 1 ))

release-minor: guard
	@make release MINOR=$$(( $(MINOR) + 1 )) PATCH=0

release-major: guard
	@make release MAJOR=$$(( $(MAJOR) + 1 )) MINOR=0 PATCH=0

release: guard
	@sed -i'.bak' 's/^VERSION=.*$$/VERSION=$(MAJOR).$(MINOR).$(PATCH)/' Makefile
	@rm Makefile.bak
	@git add Makefile
	@git commit -m 'bump version to $(MAJOR).$(MINOR).$(PATCH)'
	@git tag -a v$(MAJOR).$(MINOR).$(PATCH) -m 'v$(MAJOR).$(MINOR).$(PATCH)'
	@git push --follow-tags

guard:
	@git diff-index --quiet HEAD || (echo "There are changes in the repo, won't release. Commit everything and run this from a clean repo"; exit 1)
ifneq ($(shell echo `git branch --show-current`),master)
	@echo "Releases can only be done from master" && exit 1
endif
