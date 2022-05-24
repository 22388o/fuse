#!/usr/bin/env bash

set -euo pipefail

BTC_HOST=$1 BTC_RPC_PORT=$2 BTC_RPC_USER=$3 BTC_RPC_PASS=$4 BTC_ZQMPUBRAWBLOCK_PORT=$5 BTC_ZQMPUBRAWTX_PORT=$6 docker-compose -f ./docker-compose.lnd.yml up -d
