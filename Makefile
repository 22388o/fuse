## start: starts wallet api
.PHONY: start
start:
	go run cmd/fuse/fuse.go

.PHONY: network_up
network_up:
	@docker-compose up -d --build && make fund && make get_lnd_creds

.PHONY: network_down
network_down:
	@docker-compose down --remove-orphans

.PHONY: fund
fund:
	./scripts/fund.sh

.PHONY: get_lnd_creds
get_lnd_creds:
	rm -rf ./.fuse
	mkdir ./.fuse
	docker cp fuse_lnd:/lnd/tls.cert ./.fuse/tls.cert
	docker cp fuse_lnd:/lnd/data/chain/bitcoin/regtest/admin.macaroon ./.fuse/admin.macaroon

## help: print help message
.DEFAULT_GOAL := help
.PHONY: help
help: Makefile
	@echo "MASH"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
