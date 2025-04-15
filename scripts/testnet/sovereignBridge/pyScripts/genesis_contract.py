import json
import os
import shutil
import sys


def copy_wasm_in_project(file_path: str, wasm_path: str):
    try:
        if os.path.exists(file_path):
            os.remove(file_path)

        shutil.copy2(os.path.expanduser(wasm_path), file_path)
        print(f"{wasm_path} copied successfully.")
    except FileNotFoundError:
        print(f"File {wasm_path} not found.")
    except PermissionError:
        print("Permission denied.")
    except Exception as e:
        print(f"An error occurred: {e}")


def push_genesis_contract(file_path: str, genesis_contract):
    try:
        with open(file_path, 'r') as file:
            data = json.load(file)

        found = False
        for item in data:
            if item["filename"] == genesis_contract["filename"]:
                item.update(genesis_contract)
                found = True
                break
        if not found:
            data.append(genesis_contract)

        with open(file_path, 'w') as file:
            json.dump(data, file, indent=2)

        with open(file_path, 'a') as file:
            file.write('\n')

        print(f"genesis contract pushed successfully")
    except Exception as e:
        print(f"An error occurred: {e}")


def main():
    # input arguments
    owner_address = sys.argv[1]
    dcdt_safe_path = sys.argv[2]
    dcdt_safe_init_params = sys.argv[3]
    fee_market_path = sys.argv[4]
    fee_market_init_params = sys.argv[5]

    current_path = os.getcwd()
    project = 'drt-go-chain-sovereign'
    index = current_path.find(project)
    project_path = current_path[:index + len(project)]
    json_path = project_path + "/cmd/node/config/genesisSmartContracts.json"

    # dcdt-safe ---------------------
    dcdt_safe_wasm_path = project_path + "/cmd/node/config/genesisContracts/sov-dcdt-safe.wasm"
    copy_wasm_in_project(dcdt_safe_wasm_path, dcdt_safe_path)

    dcdt_safe_genesis_contract = {
        "owner": owner_address,
        "filename": "./config/genesisContracts/sov-dcdt-safe.wasm",
        "init-parameters": dcdt_safe_init_params,
        "vm-type": "0500",
        "type": "dcdt"
    }
    push_genesis_contract(json_path, dcdt_safe_genesis_contract)

    # fee-market -----------------
    fee_market_wasm_path = project_path + "/cmd/node/config/genesisContracts/fee-market.wasm"
    copy_wasm_in_project(fee_market_wasm_path, fee_market_path)

    fee_market_genesis_contract = {
        "owner": owner_address,
        "filename": "./config/genesisContracts/fee-market.wasm",
        "init-parameters": fee_market_init_params,
        "vm-type": "0500",
        "type": "fee"
    }
    push_genesis_contract(json_path, fee_market_genesis_contract)


if __name__ == "__main__":
    main()
