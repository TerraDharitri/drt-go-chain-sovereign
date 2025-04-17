downloadCrossChainContracts() {
    echo "Downloading cross-chain contracts..."

    mkdir -p $(eval echo "${CONTRACTS_DIRECTORY}")
    version=$(basename `curl -s https://github.com/TerraDharitri/drt-sovereign-sc/releases/latest -I | grep location | awk -F"https:/" '{print $2}' | tr -d "\r"`)
    wget -O $(eval echo ${DRT_DCDT_SAFE_WASM}) https://github.com/TerraDharitri/drt-sovereign-sc/releases/download/${version}/drt-dcdt-safe.wasm
    wget -O $(eval echo ${SOV_DCDT_SAFE_WASM}) https://github.com/TerraDharitri/drt-sovereign-sc/releases/download/${version}/sov-dcdt-safe.wasm
    wget -O $(eval echo ${FEE_MARKET_WASM}) https://github.com/TerraDharitri/drt-sovereign-sc/releases/download/${version}/fee-market.wasm
    wget -O $(eval echo ${HEADER_VERIFIER_WASM}) https://github.com/TerraDharitri/drt-sovereign-sc/releases/download/${version}/header-verifier.wasm
}
