DUSER ?= demo
GETPASS ?= 1234
POSTPASS ?= 5678
DDATA ?= $(shell echo -n "$(DUSER):$(POSTPASS)" | base64)

DOMAIN_NAME ?= http://localhost
SERVER_ADDR ?= 8080

.DEFAULT_GOAL := help
.PHONY: test-web test-checker

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ \
	{ printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

web:	## Run web server - no bin
	SERVER_ADDR=$(SERVER_ADDR) go run \
	-ldflags "-X 'main.user=$(DUSER)' \
	-X 'main.getPass=$(GETPASS)' \
	-X 'main.postPass=$(POSTPASS)'" \
	.

checker:	## Run checker - no bin
	go run \
	-ldflags "-X 'main.address=$(DOMAIN_NAME):$(SERVER_ADDR)' \
	-X 'main.data=$(DDATA)'" \
	./bin/checker.go
