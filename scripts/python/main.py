import csv
import re
import base58
import grpc
import sys
import json
import os

# Why can't??
# from pactus import network_pb2
# from pactus import network_pb2_grpc

import public_key
import network_pb2
import network_pb2_grpc
import blockchain_pb2
import blockchain_pb2_grpc
from google.protobuf.json_format import MessageToJson


def load_json_file(file_path):
    try:
        with open(file_path, "r") as file:
            data = json.load(file)
        return data
    except FileNotFoundError:
        print(f"File not found: {file_path}")
        return None
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")
        return None


def get_validator_info(grpc_validator_stub, addr):
    req = blockchain_pb2.GetValidatorRequest(address=addr)
    try:
        res = grpc_validator_stub.GetValidator(req)

        return res
    except:
        return None


def get_network_info(url):
    channel = grpc.insecure_channel(url)
    stub = network_pb2_grpc.NetworkStub(channel)
    req = network_pb2.GetNetworkInfoRequest()
    info = stub.GetNetworkInfo(req)

    return info


def extract_agent_info(input_string):
    # Regular expression pattern to extract desired information
    pattern = r"node=([^/]+)/node-version=([^/]+)/.*?/os=([^/]+)/arch=([^/]+)(/reachability=([^/]+))?"
    match = re.match(pattern, input_string)

    if match:
        node = match.group(1)
        node_version = match.group(2)
        os = match.group(3)
        arch = match.group(4)
        reachability = match.group(5)
        return node, node_version, os, arch, reachability
    else:
        return "Unknown", "Unknown", "Unknown", "Unknown", "Unknown"


def extract_network_vals(grpc_validator_stub):
    network_res = []
    network_res.append(get_network_info("95.217.181.93:50052")) # 0.21.1
    network_res.append(get_network_info("104.194.156.57:50052")) # 0.21.1
    network_res.append(get_network_info("5.75.169.99:50052")) # 0.21.1
    network_res.append(get_network_info("157.90.153.208:50052")) # 0.21.1
    network_res.append(get_network_info("81.0.218.193:50052")) # 0.21.1
    network_res.append(get_network_info("139.99.198.111:50052")) # 0.21.1
    network_res.append(get_network_info("49.13.126.38:50052")) # 0.21.1

    network_res.append(get_network_info("172.104.46.145:50052"))
    network_res.append(get_network_info("172.233.152.129:50052"))
    network_res.append(get_network_info("172.232.108.191:50052"))
    network_res.append(get_network_info("94.101.184.118:50052"))
    network_res.append(get_network_info("13.115.190.71:50052"))
    network_res.append(get_network_info("51.158.118.181:50052"))
    network_res.append(get_network_info("20.55.77.66:50052"))
    network_res.append(get_network_info("20.205.173.231:50052"))

    network_vals = {}
    for res in network_res:
        for peer in res.connected_peers:
            if len(peer.consensus_keys) == 0:
                continue

            pub = public_key.PublicKey.from_string(peer.consensus_keys[0])
            val_addr = pub.validator_address().string()

            val_info = network_vals.get(val_addr, None)
            if val_info is None:
                val_node_info = get_validator_info(grpc_validator_stub, val_addr)
                if val_node_info is None:
                    continue

                val_info = network_vals[val_addr] = {
                    "val_addr": val_addr,
                    "last_received": 0,
                    "stake": val_node_info.validator.stake / 10**9,
                    "last_sortition_height": val_node_info.validator.last_sortition_height,
                    "last_bonding_height": val_node_info.validator.last_bonding_height,
                    "unbonding_height": val_node_info.validator.unbonding_height,
                    "availability_score": val_node_info.validator.availability_score,
                    "ip": "",
                    "port": "",
                    "node": "",
                    "node_version": "",
                    "os": "",
                    "reachability": "",
                    "moniker": "",
                    "address": "",
                    "peerId": "",
                    "agent": "",
                }

            if val_info["last_received"] < peer.last_received:
                ip, port = extract_ip_and_port(peer.address)
                node, node_version, os, arch, reachability = extract_agent_info(
                    peer.agent
                )
                b58 = base58.b58encode(peer.peer_id)

                val_info["last_received"] = peer.last_received
                val_info["address"] = peer.address
                val_info["ip"] = ip
                val_info["port"] = port
                val_info["node"] = node
                val_info["node_version"] = node_version
                val_info["os"] = os
                val_info["arch"] = arch
                val_info["reachability"] = reachability
                val_info["moniker"] = peer.moniker
                val_info["peerId"] = b58[2 : len(b58) - 2]
                val_info["agent"] = peer.agent

    return network_vals


def write_network_vals(network_vals):
    # Specify the CSV file path
    csv_file_path = "output/network_vals.csv"

    # Create a list to store the rows of data for the CSV
    csv_data_users = []

    for _, value in network_vals.items():
        csv_row = [
            value["val_addr"],
            value["stake"],
            value["last_received"],
            value["last_sortition_height"],
            value["last_bonding_height"],
            value["unbonding_height"],
            value["availability_score"],
            value["ip"],
            value["port"],
            value["node"],
            value["node_version"],
            value["os"],
            value["arch"],
            value["reachability"],
            value["moniker"],
            value["address"],
            value["peerId"],
            value["agent"],
        ]

        csv_data_users.append(csv_row)

    try:
        with open(csv_file_path, mode="w", newline="") as csv_file:
            csv_writer = csv.writer(csv_file)
            # Write the header row
            header = [
                "Address",
                "Stake",
                "Last Time Online",
                "Last Sortition Height",
                "Last Bonding Height",
                "Unbonding Height",
                "Availability Score",
                "Ip",
                "Port",
                "Node",
                "Node Version",
                "Os",
                "Arch",
                "Reachability",
                "Moniker",
                "Address",
                "PeerId",
                "Agent",
            ]
            csv_writer.writerow(header)
            # Write the data rows
            csv_writer.writerows(csv_data_users)
        print(f"Data saved to {csv_file_path}")
    except Exception as e:
        print(f"Error saving data to CSV: {e}")


def extract_ip_and_port(input_string):
    # Define a regular expression to match IP address and port number
    pattern = r"/ip[46]/([^/]+)/[a-zA-Z]+/(\d+)"

    match = re.match(pattern, input_string)

    if match:
        ip_address = match.group(1)
        port = int(match.group(2))
        return ip_address, port
    else:
        return None, None


def referrals_by_discord_id(referral_data):
    # example of referral data
    #
    #  "039370": {
    #   "referral_code": "039370",
    #   "points": 0,
    #   "discord_name": "xulee99",
    #   "discord_id": "923711288759177277"
    # },
    referrals = {}
    for key, value in referral_data.items():
        referral_code = value["referral_code"]
        points = value["points"]
        discord_name = value["discord_name"]
        discord_id = value["discord_id"]

        referral_data = referrals.get(discord_id, None)
        if referral_data is None:
            if referral_code != key:
                print("Duplicated referral code")
                sys.exit(1)

            referrals[discord_id] = {
                "referral_code": referral_code,
                "points": points,
                "discord_name": discord_name,
                "discord_id": discord_id,
            }
        else:
            print("Duplicated referral")
            sys.exit(1)

    return referrals


def extract_users_map(validator_data, referrals):
    # Example of validator data
    #
    # "12D3KooWSbDJWMhYgrFqymf78q4vhZEAV8n1LUUZ7Y9VenR6PNdN": {
    #   "discord_name": "vikanren",
    #   "discord_id": "840519270004686878",
    #   "validator_address": "tpc1pjrvumvpsutpklgg2hwhuaaujc3pju6z9kkffzt",
    #   "referrer_discord_id": "",
    #   "faucet_amount": 100
    # },
    users_map = {}
    for _, value in validator_data.items():
        discord_name = value["discord_name"]
        discord_id = value["discord_id"]
        validator_address = value["validator_address"]
        referrer_discord_id = value["referrer_discord_id"]
        faucet_amount = value["faucet_amount"]

        user_data = users_map.get(discord_id, None)
        if user_data is None:
            user_data = users_map[discord_id] = {
                "discord_id": discord_id,
                "discord_name": set(),
                "faucet_amount": [],
                "referrer_discord_id": [],
                "referrer_discord_name": [],
                "referral_points": 0,
                "referral_code": 0,
                "stakes": "",
                "total_stakes": 0,
                "total_reward": 0,
                "total_stakes_online": 0,
                "total_reward_online": 0,
                "num_validators": 0,
                "validators": [],
            }

        user_data["discord_name"].add(discord_name)
        user_data["faucet_amount"].append(faucet_amount)
        user_data["validators"].append(validator_address)

        if discord_id in referrals:
            ref_data = referrals[discord_id]
            referral_code = ref_data["referral_code"]

            user_data["referral_points"] = ref_data["points"]
            user_data["referral_code"] = referral_code

        if referrer_discord_id != "":
            ref_data = referrals[referrer_discord_id]

            user_data["referrer_discord_id"].append(referrer_discord_id)
            user_data["referrer_discord_name"].append(ref_data["discord_name"])

    return users_map


def write_users_map(users_map):
    # Specify the CSV file path
    csv_file_path = "output/users_map.csv"

    # Create a list to store the rows of data for the CSV
    csv_data_users = []

    for _, value in users_map.items():
        csv_row = [
            value["discord_id"],
            value["discord_name"],
            value["total_stakes_online"],
            value["total_reward_online"],
            value["total_stakes"],
            value["total_reward"],
            value["faucet_amount"],
            value["referrer_discord_id"],
            value["referrer_discord_name"],
            value["referral_points"],
            value["referral_code"],
            value["num_validators"],
            value["validators"],
            value["stakes"],
        ]

        csv_data_users.append(csv_row)

    try:
        with open(csv_file_path, mode="w", newline="") as csv_file:
            csv_writer = csv.writer(csv_file)
            # Write the header row
            header = [
                "Discord ID",
                "Discord Name",
                "Total Stakes (Online)",
                "Total Rewards (Online)",
                "Total Stakes",
                "Total Rewards",
                "Faucet amount",
                "Referrer Discord Id",
                "Referrer Discord Name",
                "Referral Points",
                "Referral Code",
                "Num of Validators",
                "Validators",
                "Stakes",
            ]
            csv_writer.writerow(header)
            # Write the data rows
            csv_writer.writerows(csv_data_users)
        print(f"Data saved to {csv_file_path}")
    except Exception as e:
        print(f"Error saving data to CSV: {e}")


def extract_vals_map(validator_data, referrals, network_vals):
    # Example of validator data
    #
    # "12D3KooWSbDJWMhYgrFqymf78q4vhZEAV8n1LUUZ7Y9VenR6PNdN": {
    #   "discord_name": "vikanren",
    #   "discord_id": "840519270004686878",
    #   "validator_address": "tpc1pjrvumvpsutpklgg2hwhuaaujc3pju6z9kkffzt",
    #   "referrer_discord_id": "",
    #   "faucet_amount": 100
    # },
    vals_map = {}
    for _, value in validator_data.items():
        discord_name = value["discord_name"]
        discord_id = value["discord_id"]
        validator_address = value["validator_address"]
        referrer_discord_id = value["referrer_discord_id"]
        faucet_amount = value["faucet_amount"]

        val_data = vals_map.get(validator_address, None)
        if val_data is None:
            stake = "0"
            last_received = 0
            last_sortition_height = "0"
            last_bonding_height = "0"
            unbonding_height = "0"
            availability_score = "0"

            if validator_address in network_vals:
                network_val_info = network_vals[validator_address]

                last_received = network_val_info["last_received"]

            val_node_info = get_validator_info(grpc_validator_stub, validator_address)

            stake = val_node_info.validator.stake / 10**9
            last_sortition_height = val_node_info.validator.last_sortition_height
            last_bonding_height = val_node_info.validator.last_bonding_height
            unbonding_height = val_node_info.validator.unbonding_height
            availability_score = val_node_info.validator.availability_score

            val_data = vals_map[validator_address] = {
                "validator_address": validator_address,
                "discord_id": discord_id,
                "discord_name": discord_name,
                "faucet_amount": faucet_amount,
                "stake": stake,
                "referrer_discord_id": referrer_discord_id,
                "referrer_discord_name": "",
                "last_received": last_received,
                "last_sortition_height": last_sortition_height,
                "last_bonding_height": last_bonding_height,
                "unbonding_height": unbonding_height,
                "availability_score": availability_score,
            }

        if referrer_discord_id != "":
            ref_data = referrals[referrer_discord_id]

            val_data["referrer_discord_name"] = ref_data["discord_name"]

    return vals_map


def write_vals_map(vals_map):
    # Specify the CSV file path
    csv_file_path = "output/vals_map.csv"

    # Create a list to store the rows of data for the CSV
    csv_data_vals = []

    for _, value in vals_map.items():
        csv_row = [
            value["validator_address"],
            value["discord_id"],
            value["discord_name"],
            value["faucet_amount"],
            value["stake"],
            value["referrer_discord_id"],
            value["referrer_discord_name"],
            value["last_received"],
            value["last_sortition_height"],
            value["last_bonding_height"],
            value["unbonding_height"],
            value["availability_score"],
        ]

        csv_data_vals.append(csv_row)

    try:
        with open(csv_file_path, mode="w", newline="") as csv_file:
            csv_writer = csv.writer(csv_file)
            # Write the header row
            header = [
                "Validator Address",
                "Discord Id",
                "Discord Name",
                "Faucet Amount",
                "Stake",
                "Referrer Discord Id",
                "Referrer Discord Name",
                "Last Time Online",
                "Last Sortition Height",
                "Last Bonding Height",
                "Unbonding Height",
                "Availability Score",
            ]
            csv_writer.writerow(header)
            # Write the data rows
            csv_writer.writerows(csv_data_vals)
        print(f"Data saved to {csv_file_path}")
    except Exception as e:
        print(f"Error saving data to CSV: {e}")


def calculate_rewards(users_map, vals_map):
    for _, user in users_map.items():
        total_stakes = 0
        total_reward = 0
        total_stakes_online = 0
        total_reward_online = 0

        total_reward = user["referral_points"]
        total_reward_online = user["referral_points"]

        for val_addr in user["validators"]:
            val_stake = vals_map[val_addr]["stake"]
            if val_stake == "0":
                print("Something is wrong: " + val_addr)
                sys.exit(1)

            total_stakes += val_stake
            if vals_map[val_addr]["last_received"] > 0:
                total_stakes_online += val_stake

        total_reward += total_stakes / 10
        total_reward_online += total_stakes_online / 10

        user["total_stakes"] = total_stakes
        user["total_reward"] = total_reward
        user["total_stakes_online"] = total_stakes_online
        user["total_reward_online"] = total_reward_online


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python main.py <validators_path> <referral_path>")
        sys.exit(1)

    validator_path = sys.argv[1]
    referral_path = sys.argv[2]

    grpc_validator_channel = grpc.insecure_channel("172.104.46.145:50052")
    grpc_validator_stub = blockchain_pb2_grpc.BlockchainStub(grpc_validator_channel)

    validator_data = load_json_file(validator_path)
    referral_data = load_json_file(referral_path)

    network_vals = extract_network_vals(grpc_validator_stub)
    referrals = referrals_by_discord_id(referral_data)
    vals_map = extract_vals_map(validator_data, referrals, network_vals)
    users_map = extract_users_map(validator_data, referrals)

    calculate_rewards(users_map, vals_map)

    write_network_vals(network_vals)
    write_users_map(users_map)
    write_vals_map(vals_map)
