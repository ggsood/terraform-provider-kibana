SHELL := /bin/bash
export GO111MODULE ?= on
export VERSION := "1.0.4"
export BINARY := terraform-provider-kibana
export GOBIN = $(shell pwd)/bin

include scripts/Makefile.help
.DEFAULT_GOAL := help

include build/Makefile.build
include build/Makefile.test
include build/Makefile.dev
include build/Makefile.deps
include build/Makefile.tools
include build/Makefile.lint
include build/Makefile.format
