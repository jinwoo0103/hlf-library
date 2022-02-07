#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1
starttime=$(date +%s)

# clean out any old identites in the wallets
rm -rf application-javascript/wallet/*

pushd ./application-javascript
npm install
popd

pushd ./chaincode-go
go mod init main
go mod tidy
go mod vendor
popd

# launch network; create channel and join peer to channel
pushd ../test-network
./network.sh down
./network.sh up createChannel -ca
./network.sh deployCC -ccn library -ccv 1 -cci initLedger -ccl go -ccp "../library/chaincode-go/"
popd

pushd ./application-javascript
node enrollAdmin.js
node registerAndEnrollUser.js
popd
