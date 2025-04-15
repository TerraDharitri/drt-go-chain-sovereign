DCDT_SAFE_ADDRESS=$(drtpy data load --use-global --partition=${CHAIN_ID} --key=address-dcdt-safe-contract)
DCDT_SAFE_ADDRESS_SOVEREIGN=$(drtpy data load --use-global --partition=sovereign --key=address-dcdt-safe-contract)

deployDcdtSafeContract() {
    echo "Deploying DCDT Safe contract on main chain..."

    local OUTFILE="${OUTFILE_PATH}/deploy-dcdt-safe.interaction.json"
    drtpy contract deploy \
        --bytecode=$(eval echo ${DRT_DCDT_SAFE_WASM}) \
        --pem=${WALLET} \
        --proxy=${PROXY} \
        --chain=${CHAIN_ID} \
        --gas-limit=200000000 \
        --arguments \
            ${HEADER_VERIFIER_ADDRESS} \
        --outfile=${OUTFILE} \
        --recall-nonce \
        --wait-result \
        --send || return

    printTxStatus ${OUTFILE} || return

    local ADDRESS=$(drtpy data parse --file=${OUTFILE} --expression="data['contractAddress']")
    drtpy data store --use-global --partition=${CHAIN_ID} --key=address-dcdt-safe-contract --value=${ADDRESS}
    DCDT_SAFE_ADDRESS=$(drtpy data load --use-global --partition=${CHAIN_ID} --key=address-dcdt-safe-contract)
    echo -e "DCDT Safe contract: ${ADDRESS}"

    local SOVEREIGN_CONTRACT_ADDRESS=$(computeFirstSovereignContractAddress)
    drtpy data store --use-global --partition=sovereign --key=address-dcdt-safe-contract --value=${SOVEREIGN_CONTRACT_ADDRESS}
    DCDT_SAFE_ADDRESS_SOVEREIGN=$(drtpy data load --use-global --partition=sovereign --key=address-dcdt-safe-contract)
    echo -e "DCDT Safe sovereign contract: ${SOVEREIGN_CONTRACT_ADDRESS}\n"
}

upgradeDcdtSafeContract() {
    echo "Upgrading DCDT Safe contract on main chain..."
    checkVariables DCDT_SAFE_ADDRESS || return

    local OUTFILE="${OUTFILE_PATH}/upgrade-dcdt-safe.interaction.json"
    upgradeDcdtSafeContractCall ${DCDT_SAFE_ADDRESS} ${PROXY} ${CHAIN_ID} ${OUTFILE}
}

upgradeDcdtSafeContractSovereign() {
    echo "Upgrading DCDT Safe contract on sovereign chain..."
    checkVariables DCDT_SAFE_ADDRESS_SOVEREIGN || return

    local OUTFILE="${OUTFILE_PATH}/upgrade-dcdt-safe-sovereign.interaction.json"
    upgradeDcdtSafeContractCall ${DCDT_SAFE_ADDRESS_SOVEREIGN} ${PROXY_SOVEREIGN} ${CHAIN_ID_SOVEREIGN} ${OUTFILE}
}

upgradeDcdtSafeContractCall() {
    if [ $# -lt 4 ]; then
        echo "Usage: ${FUNCNAME[0]} <arg1> <arg2> <arg3> <arg4>"
        return 1
    fi

    local ADDRESS=$1
    local URL=$2
    local CHAIN=$3
    local OUTFILE=$4

    drtpy contract upgrade ${ADDRESS} \
        --bytecode=$(eval echo ${DRT_DCDT_SAFE_WASM}) \
        --pem=${WALLET} \
        --proxy=${URL} \
        --chain=${CHAIN} \
        --gas-limit=200000000 \
        --outfile=${OUTFILE} \
        --recall-nonce \
        --wait-result \
        --send || return

    printTxStatus ${OUTFILE}
}

pauseDcdtSafeContract() {
    echo "Pausing DCDT Safe contract on main chain..."

    local OUTFILE="${OUTFILE_PATH}/pause-dcdt-safe.interaction.json"
    pauseDcdtSafeContractCall ${DCDT_SAFE_ADDRESS} ${PROXY} ${CHAIN_ID} ${OUTFILE}
}
pauseDcdtSafeContractSovereign() {
    echo "Pausing DCDT Safe contract on sovereign chain..."

    local OUTFILE="${OUTFILE_PATH}/pause-dcdt-safe-sovereign.interaction.json"
    pauseDcdtSafeContractCall ${DCDT_SAFE_ADDRESS_SOVEREIGN} ${PROXY_SOVEREIGN} ${CHAIN_ID_SOVEREIGN} ${OUTFILE}
}
pauseDcdtSafeContractCall() {
    if [ $# -lt 4 ]; then
        echo "Usage: ${FUNCNAME[0]} <arg1> <arg2> <arg3> <arg4>"
        return 1
    fi

    local ADDRESS=$1
    local URL=$2
    local CHAIN=$3
    local OUTFILE=$4

    drtpy contract call ${ADDRESS} \
        --pem=${WALLET} \
        --proxy=${URL} \
        --chain=${CHAIN} \
        --gas-limit=10000000 \
        --function="pause" \
        --outfile=${OUTFILE} \
        --recall-nonce \
        --wait-result \
        --send || return

    printTxStatus ${OUTFILE}
}

unpauseDcdtSafeContract() {
    echo "Unpausing DCDT Safe contract on main chain..."

    local OUTFILE="${OUTFILE_PATH}/unpause-dcdt-safe.interaction.json"
    unpauseDcdtSafeContractCall ${DCDT_SAFE_ADDRESS} ${PROXY} ${CHAIN_ID} ${OUTFILE}
}
unpauseDcdtSafeContractSovereign() {
    echo "Unpausing DCDT Safe contract on sovereign chain..."

    local OUTFILE="${OUTFILE_PATH}/unpause-dcdt-safe-sovereign.interaction.json"
    unpauseDcdtSafeContractCall ${DCDT_SAFE_ADDRESS_SOVEREIGN} ${PROXY_SOVEREIGN} ${CHAIN_ID_SOVEREIGN} ${OUTFILE}
}
unpauseDcdtSafeContractCall() {
    if [ $# -lt 4 ]; then
        echo "Usage: ${FUNCNAME[0]} <arg1> <arg2> <arg3> <arg4>"
        return 1
    fi

    local ADDRESS=$1
    local URL=$2
    local CHAIN=$3
    local OUTFILE=$4

    drtpy contract call ${ADDRESS} \
        --pem=${WALLET} \
        --proxy=${URL} \
        --chain=${CHAIN} \
        --gas-limit=10000000 \
        --function="unpause" \
        --outfile=${OUTFILE} \
        --recall-nonce \
        --wait-result \
        --send || return

    printTxStatus ${OUTFILE}
}

setFeeMarketAddress() {
    echo "Setting Fee Market address in DCDT Safe contract on main chain..."

    local OUTFILE="${OUTFILE_PATH}/set-feemarket-address.interaction.json"
    setFeeMarketAddressCall ${DCDT_SAFE_ADDRESS} ${FEE_MARKET_ADDRESS} ${PROXY} ${CHAIN_ID} ${OUTFILE}
}
setFeeMarketAddressSovereign() {
    echo "Setting Fee Market address in DCDT Safe contract on sovereign chain..."

    local OUTFILE="${OUTFILE_PATH}/set-feemarket-address-sovereign.interaction.json"
    setFeeMarketAddressCall ${DCDT_SAFE_ADDRESS_SOVEREIGN} ${FEE_MARKET_ADDRESS_SOVEREIGN} ${PROXY_SOVEREIGN} ${CHAIN_ID_SOVEREIGN} ${OUTFILE}
}
setFeeMarketAddressCall() {
    if [ $# -lt 5 ]; then
        echo "Usage: ${FUNCNAME[0]} <arg1> <arg2> <arg3> <arg4> <arg5>"
        return 1
    fi

    local DCDT_SAFE_CONTRACT_ADDRESS=$1
    local FEE_MARKET_CONTRACT_ADDRESS=$2
    local URL=$3
    local CHAIN=$4
    local OUTFILE=$5

    drtpy contract call ${DCDT_SAFE_CONTRACT_ADDRESS} \
        --pem=${WALLET} \
        --proxy=${URL} \
        --chain=${CHAIN} \
        --gas-limit=10000000 \
        --function="setFeeMarketAddress" \
        --arguments ${FEE_MARKET_CONTRACT_ADDRESS} \
        --outfile=${OUTFILE} \
        --recall-nonce \
        --wait-result \
        --send || return

    printTxStatus ${OUTFILE}
}

registerNativeToken() {
    echo "Registering new native token in DCDT Safe contract on main chain..."
    checkVariables DCDT_SAFE_ADDRESS || return

    local OUTFILE="${OUTFILE_PATH}/register-native-token.interaction.json"
    drtpy contract call ${DCDT_SAFE_ADDRESS} \
        --pem=${WALLET} \
        --proxy=${PROXY} \
        --chain=${CHAIN_ID} \
        --value=${DCDT_ISSUE_COST} \
        --gas-limit=100000000 \
        --function="registerNativeToken" \
        --arguments \
            str:${NATIVE_DCDT_TICKER} \
            str:${NATIVE_DCDT_NAME} \
        --outfile=${OUTFILE} \
        --recall-nonce \
        --wait-result \
        --send || return

    printTxStatus ${OUTFILE}

    local RESULT=$(readNativeToken)
    local TOKEN_HEX=$(echo "$RESULT" | jq -r '.[0]')
    NATIVE_DCDT=$(hex_to_string "$TOKEN_HEX")
    echo "Native Token identifier: ${NATIVE_DCDT}"
}

readNativeToken() {
    checkVariables DCDT_SAFE_ADDRESS || return

    drtpy contract query ${DCDT_SAFE_ADDRESS} \
        --proxy=${PROXY} \
        --function="getNativeToken"
}
