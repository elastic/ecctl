SHELL := /bin/bash
export VERSION ?= v1.14.2
export GO111MODULE ?= on
export GOBIN = $(shell pwd)/bin
BINARY := ecctl

include scripts/Makefile.help
.DEFAULT_GOAL := help

include build/Makefile.build
include build/Makefile.dev
include build/Makefile.deps
include build/Makefile.version
