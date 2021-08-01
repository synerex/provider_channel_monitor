# Makefile for Synerex.

GOCMD=go
GOBUILD=$(GOCMD) build -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=`date +%Y-%m-%d_%T` -X main.gitver=`git describe --tag`"
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
RM=rm


# Main target

.PHONY: build 
build: channel_monitor

channel_monitor: channel_monitor.go 
	$(GOBUILD)

.PHONY: clean
clean: 
	$(RM) channel_monitor



