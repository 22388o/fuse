## start: starts wallet api
.PHONY: start
start:
	go run cmd/fuse/fuse.go

# network_up: starts local lnd + bitcoin node
.PHONY: network_up
network_up:
	@docker-compose up -d --build && make fund && make get_lnd_creds

# network_down: stops local lnd + bitcoin node
.PHONY: network_down
network_down:
	@docker-compose down --remove-orphans

# start_lnd: starts lnd node
.PHONY: start_lnd
start_lnd:
	docker-compose -f ./docker-compose.lnd.yml --env-file ./.env up -d  && make get_lnd_creds

# stop_lnd: stops lnd node
.PHONY: stop_lnd
stop_lnd:
	docker-compose -f ./docker-compose.lnd.yml down

# fund: funds local lnd node
.PHONY: fund
fund:
	./scripts/fund.sh

# get_lnd_creds: copies tls.cert and admin.macaroon from lnd node into .fuse
.PHONY: get_lnd_creds
get_lnd_creds:
	rm -rf ./.fuse
	mkdir ./.fuse
	docker cp fuse_lnd:/lnd/tls.cert ./.fuse/tls.cert
	docker cp fuse_lnd:/lnd/data/chain/bitcoin/regtest/admin.macaroon ./.fuse/admin.macaroon

# fusecli: installs the fuse cli client at bin/mapi
.PHONY: fusecli
fusecli:
	GOBIN=$(shell pwd)/bin go install cmd/fusecli/fusecli.go

# test: runs tests
.PHONY: test
test:
	@GOBIN=$(FUSE_DIR)/bin go install github.com/mfridman/tparse@v0.9.0
	go test -json -coverprofile cover.out ./... | $(FUSE_DIR)/bin/tparse -nocolor

## help: print help message
.DEFAULT_GOAL := help
.PHONY: help
help: Makefile
	@echo "MASH"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
