BINNAME ?= clipper
BINDIR ?=
GO ?= go
GOBIN ?= $(shell $(GO) env GOBIN)
GOLANGCI_LINT ?= golangci-lint
GOFUMPT ?= gofumpt
TEST_FLAGS ?= -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt

ifeq ($(OS),Windows_NT)
BINEXT := .exe
ifeq ($(strip $(BINDIR)),)
BINDIR := $(LOCALAPPDATA)/Programs/clipper
endif
else
BINEXT :=
ifeq ($(strip $(BINDIR)),)
BINDIR := ${HOME}/.local/bin
endif
endif

ifeq ($(strip $(GOBIN)),)
GOBIN := $(shell $(GO) env GOPATH)/bin
endif
# Go on Windows may report backslashes; normalize for recipe portability.
GOBIN := $(subst \,/,$(GOBIN))
BINDIR := $(subst \,/,$(BINDIR))

INSTSRC := $(GOBIN)/$(BINNAME)$(BINEXT)
INSTDEST := $(BINDIR)/$(BINNAME)$(BINEXT)

.PHONY: all build test lint fmt check install

all: build check

build:
	$(GO) install ./...

test:
	$(GO) test $(TEST_FLAGS) ./...

lint:
	$(GOLANGCI_LINT) run

fmt:
	$(GOFUMPT) -w .

check: test lint

install: build
	@echo "Installing $(BINNAME)..."
ifeq ($(OS),Windows_NT)
	@powershell -NoProfile -Command "$$ErrorActionPreference='Stop'; New-Item -ItemType Directory -Force -Path '$(BINDIR)' | Out-Null; Move-Item -LiteralPath '$(INSTSRC)' -Destination '$(INSTDEST)' -Force"
	@powershell -NoProfile -Command "$$ErrorActionPreference='Stop'; $$installDir=[System.IO.Path]::GetFullPath('$(BINDIR)'); $$userPath=[Environment]::GetEnvironmentVariable('Path','User'); if ($$userPath -notlike ('*'+$$installDir+'*')) { $$newPath = if ([string]::IsNullOrEmpty($$userPath)) { $$installDir } else { $$installDir+';'+$$userPath }; [Environment]::SetEnvironmentVariable('Path', $$newPath, 'User'); Write-Host ('Added to your user PATH: '+$$installDir) }"
else
	@mkdir -p "$(BINDIR)"
	@mv "$(INSTSRC)" "$(INSTDEST)"
endif
	@echo "Done."
