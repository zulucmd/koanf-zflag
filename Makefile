.PHONY: lint test unittest clean

ifeq (, $(shell which golangci-lint 2> /dev/null))
	$(error Unable to locate golangci-lint! Ensure it is installed and available in PATH before re-running.)
endif

gotest =
ifeq (, $(shell which richgo 2> /dev/null))
	gotest = go test
else
	gotest = richgo test
endif

default: all

all: lint test clean

lint:
	@echo '********** LINT TEST **********'
	golangci-lint run

unittest:
	@echo '********** UNIT TEST **********'
	@$(gotest) -failfast -v -race -cover

test: unittest lint
