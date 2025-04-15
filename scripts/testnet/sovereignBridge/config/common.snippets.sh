checkWalletBalanceOnMainChain() {
    local BALANCE=$(drtpy account get --address ${WALLET_ADDRESS} --proxy ${PROXY} --balance)
    if [ "$BALANCE" == "0" ]; then
        echo -e "Your wallet balance is zero on main chain"
        return 1
    fi
    return 0
}

fund() {
    if [ "$#" -ne 1 ]; then
        echo "Usage: fund <address>"
        return 1
    fi

    echo "Funding wallet address $1 on sovereign chain..."

    local OUTFILE="${OUTFILE_PATH}/get-funds-sovereign.interaction.json"
    drtpy tx new \
        --pem=${WALLET_SOVEREIGN} \
        --proxy=${PROXY_SOVEREIGN} \
        --chain=${CHAIN_ID_SOVEREIGN} \
        --receiver=$1 \
        --value=100000000000000000000000 \
        --gas-limit=50000 \
        --outfile=${OUTFILE} \
        --recall-nonce \
        --wait-result \
        --send
}

gitPullAllChanges()
{
    pushd .

    # Traverse up to the parent directory of "drt-go-chain-sovereign"
    while [[ ! -d "drt-go-chain-sovereign" && $(pwd) != "/" ]]; do
      cd ..
    done

    # Check if we found the directory
    if [[ ! -d "drt-go-chain-sovereign" ]]; then
      echo "drt-go-chain-sovereign directory not found"
      popd
      return 1
    fi

    echo -e "Pulling changes for drt-go-chain-sovereign..."
    cd drt-go-chain-sovereign
    git pull
    cd ..

    echo -e "Pulling changes for drt-go-chain-deploy-sovereign..."
    cd drt-go-chain-deploy-sovereign
    git pull
    cd ..

    echo -e "Pulling changes for drt-go-chain-proxy-sovereign..."
    cd drt-go-chain-proxy-sovereign
    git pull
    cd ..

    echo -e "Pulling changes for drt-go-chain-sovereign-bridge..."
    pushd .
    cd drt-go-chain-sovereign-bridge
    git pull
    cd cert/cmd/cert
    go build
    ./cert
    popd

    echo -e "Pulling changes for drt-go-chain-tools..."
    pushd .
    cd drt-go-chain-tools
    git pull
    cd elasticreindexer/cmd/indices-creator/
    go build
    popd

    popd
}
