# Fuse LN Wallet API

Simple Wallet API that supports LNURL. Meant to run against local docker environment to easily test different functionality when developing on the lightning network.

Currently supports:

- lnurlp: [LUD-06](https://github.com/fiatjaf/lnurl-rfc/blob/luds/06.md)

## Getting Started

### No Existing Environment

This method is for when you have no existing bitcoind or lnd nodes running on your local environment.

1. To get the network up and running: `make network_up`
2. Start API: `make start`

You will now have a LND Node, Bitcoin Node and the API running.

### Existing bitcoind + lightning nodes

This method is for when you want to hook into an existing lightning network running on docker locally.

> This guide currently assumes bitcoind is running under container name `bitcoind`. It also has some preconfigured rpc creds it assumes. These can be be found in the makefile

1. Start your LND node: `NETWORK={docker_network_name} make start_lnd`
1. Pull your LND creds: `make get_lnd_creds`
1. Fund your node: `make fund`

### fusecli

There is a cli to help interact with the API and lightning network

1. Build CLI: `make fusecli`
1. Add bin directory to your path or just use `./bin/fusecli`

## Interacting with API and Network

The simplest method to interact with the network locally is to open a channel through the API. 

1. Open a channel: `fusecli channels new -node <pubkey@host> <localSat> <pushSat>`
1. Mine some blocks to confirm the channel:
  a. `fusecli ln newaddress` 
  a. `fusecli btcd mine -block <num_blocks> <address_from_prev_step>`
