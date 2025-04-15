computeFirstSovereignContractAddress() {
    echo $(python3 $SCRIPT_PATH/pyScripts/compute_contract_address.py $WALLET_ADDRESS 0)
}

computeSecondSovereignContractAddress() {
    echo $(python3 $SCRIPT_PATH/pyScripts/compute_contract_address.py $WALLET_ADDRESS 1)
}

getShardOfAddress() {
  echo $(python3 $SCRIPT_PATH/pyScripts/address_shard.py $WALLET_ADDRESS)
}

bech32ToHex() {
  echo $(python3 $SCRIPT_PATH/pyScripts/address_convert.py $1)
}

displayContracts() {
    echo "DCDT-SAFE: $DCDT_SAFE_ADDRESS"
    echo "DCDT-SAFE SOVEREIGN: $DCDT_SAFE_ADDRESS_SOVEREIGN"
    echo "FEE-MARKET: $FEE_MARKET_ADDRESS"
    echo "FEE-MARKET SOVEREIGN: $FEE_MARKET_ADDRESS_SOVEREIGN"
    echo "HEADER-VERIFIER: $HEADER_VERIFIER_ADDRESS"
}

updateAndStartBridgeService() {
    python3 $SCRIPT_PATH/pyScripts/bridge_service.py $WALLET $PROXY $DCDT_SAFE_ADDRESS $HEADER_VERIFIER_ADDRESS
}

updateNotifierNotarizationRound() {
    python3 $SCRIPT_PATH/pyScripts/notifier_round.py $PROXY $(getShardOfAddress)
}

setGenesisContract() {
    local DCDT_SAFE_INIT_PARAMS="$(bech32ToHex $FEE_MARKET_ADDRESS_SOVEREIGN)"
    local FEE_MARKET_INIT_PARAMS="$(bech32ToHex $DCDT_SAFE_ADDRESS_SOVEREIGN)@00"

    python3 $SCRIPT_PATH/pyScripts/genesis_contract.py $WALLET_ADDRESS $SOV_DCDT_SAFE_WASM $DCDT_SAFE_INIT_PARAMS $FEE_MARKET_WASM $FEE_MARKET_INIT_PARAMS
}

updateSovereignConfig() {
    if [ -z "$1" ]; then
        DCDT_PREFIX=$(generateRandomDcdtPrefix)
    else
        DCDT_PREFIX=$1
    fi

    if [ -z "$NATIVE_DCDT" ]; then
        echo "Error: NATIVE_DCDT was not registered"
        return 1
    fi

    python3 $SCRIPT_PATH/pyScripts/update_toml.py $DCDT_SAFE_ADDRESS $DCDT_SAFE_ADDRESS_SOVEREIGN $DCDT_PREFIX $USE_ELASTICSEARCH $MAIN_CHAIN_ELASTIC $NATIVE_DCDT
}

generateRandomDcdtPrefix() {
  LEN=$(shuf -i 1-4 -n 1)
  RANDOM_PREFIX=$(cat /dev/urandom | tr -dc 'a-z0-9' | head -c $LEN)

  echo $RANDOM_PREFIX
}
