import os
import re
import subprocess
import sys
from datetime import datetime


def update_env(lines, identifier, value) -> []:
    updated_lines = []

    for line in lines:
        if line.startswith(identifier):
            line = re.sub(r'"(.*?)"', f'"{value}"', line)

        updated_lines.append(line)

    return updated_lines


def build_and_run_bridge_server(server_path):
    os.chdir(server_path)

    build_command = "go build"
    kill_screen_service = "screen -ls | grep 'sovereignBridgeService' | awk -F. '{print $1}' | xargs -I{} screen -X -S {} quit"
    current_time = datetime.now().strftime("%Y%m%d_%H%M%S")
    run_service = f"screen -dmS sovereignBridgeService -L -Logfile sovereignBridgeService_{current_time}.log ./server"

    build_process = subprocess.run(build_command, shell=True, capture_output=True, text=True)
    if build_process.returncode == 0:
        print("Go build successful")
    else:
        print("Error during Go build")
        return

    subprocess.run(kill_screen_service, shell=True, capture_output=True, text=True)

    # TODO start terminal with server app
    run_process = subprocess.run(run_service, shell=True, capture_output=True, text=True)
    if run_process.returncode == 0:
        print("Bridge service started successfully")
    else:
        print("Error starting bridge service")
        return


def main():
    # input arguments
    wallet = sys.argv[1]
    proxy = sys.argv[2]
    dcdt_safe_address = sys.argv[3]
    header_verifier_address = sys.argv[4]

    current_path = os.getcwd()
    project = 'drt-go-chain-sovereign'
    index = current_path.find(project)
    project_path = current_path[:index]
    bridge_service_path = os.path.join(project_path, 'drt-go-chain-sovereign-bridge')
    server_path = bridge_service_path + "/server/cmd/server"
    env_path = server_path + "/.env"

    with open(env_path, 'r') as file:
        lines = file.readlines()

    updated_lines = update_env(lines, "WALLET_PATH", os.path.expanduser(wallet))
    updated_lines = update_env(updated_lines, "DHARITRI_PROXY", os.path.expanduser(proxy))
    updated_lines = update_env(updated_lines, "HEADER_VERIFIER_SC_ADDRESS", header_verifier_address)
    updated_lines = update_env(updated_lines, "DCDT_SAFE_SC_ADDRESS", dcdt_safe_address)
    updated_lines = update_env(updated_lines, "CERT_FILE", os.path.expanduser("~/Dharitri/testnet/node/config/certificate.crt"))
    updated_lines = update_env(updated_lines, "CERT_PK_FILE", os.path.expanduser("~/Dharitri/testnet/node/config/private_key.pem"))

    with open(env_path, 'w') as file:
        file.writelines(updated_lines)

    build_and_run_bridge_server(server_path)


if __name__ == "__main__":
    main()
