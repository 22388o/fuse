#!/usr/bin/env bash

set -euo pipefail

BTC_CONTAINER=${1:-bitcoind}

# standard options for lncli and bitcoin-cli
LNCLI="lncli --lnddir=/lnd -n regtest"
BITCOINCLI="bitcoin-cli -chain=regtest -rpcuser=regtest -rpcpassword=regtest -rpcwait"

# run a command in a docker compose managed container
CONTAINER_EXEC="docker exec"
BITCOIND_CMD="${CONTAINER_EXEC} ${BTC_CONTAINER} ${BITCOINCLI}"

# mine $1 blocks to $2 address
mine () {
  ${BITCOIND_CMD} generatetoaddress $1 $2 > /dev/null
}

# send $1 bitcoin from bitcoind (pre-mined) to $2 address
send () {
  ${BITCOIND_CMD} sendtoaddress $2 $1 > /dev/null
  # just mining some blocks to make official
  mine 6 $LND_ADDRESS
}

# fund the $1 LND node
fund_lnd () {
  # make sure lnd has started
  echo -n Fund $1... 1>&2
  until ${CONTAINER_EXEC} $1 ${LNCLI} getinfo > /dev/null 2>&1
  do
    echo -n "." 1>&2
    sleep 1
  done

  # short circuit if funded
  if [ $(${CONTAINER_EXEC} $1 ${LNCLI} walletbalance 2> /dev/null | jq '.confirmed_balance|tonumber') -gt 0 ] > /dev/null 2>&1
  then 
    echo "already funded" 1>&2
    return 0
  fi 

  # generate new address and 10 btc for funds
  LND_ADDRESS=$(${CONTAINER_EXEC} $1 ${LNCLI} newaddress np2wkh | jq -r .address)
  send 10 $LND_ADDRESS

  echo "done" 1>&2 
}

# generate a bunch of blocks to get some btc to premine 
# need at least 100 to trigger payout in regtest and maybe 400ish for segwit
${BITCOIND_CMD} createwallet regtest > /dev/null 2>&1 || ${BITCOIND_CMD} loadwallet regtest > /dev/null 2>&1 || true
mine 401 $(${BITCOIND_CMD} getnewaddress)

fund_lnd "fuse_lnd"

exit 0
