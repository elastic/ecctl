SHELL := /bin/bash
export VERSION ?= 1.0.0-rc1
# Required for Go 1.13 onwards. Remove after OSS.
export GOPRIVATE ?= github.com/elastic/cloud-sdk-go
export GO111MODULE ?= on
export GOBIN = $(shell pwd)/bin
BINARY := ecctl

include scripts/Makefile.help
.DEFAULT_GOAL := help

include build/Makefile.build
include build/Makefile.dev
include build/Makefile.deps
