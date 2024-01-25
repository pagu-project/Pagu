import csv
import re
import base58
import grpc
import sys
import json
import os
import socket

# Why can't??
# from pactus import network_pb2
# from pactus import network_pb2_grpc

import public_key
import network_pb2
import network_pb2_grpc
import blockchain_pb2
import blockchain_pb2_grpc


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


def is_port_open(ip, port):
    # Create a socket object
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.settimeout(2)  # Set a timeout for the connection attempt

    try:
        # Attempt to connect to the IP address and port
        sock.connect((ip, port))
        return True  # Port is open
    except socket.error:
        return False  # Port is closed
    finally:
        sock.close()


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

        if reachability is not None:
            if "/reachability=" in reachability:
                reachability = reachability.replace("/reachability=", "")

            if reachability == "pactus-testnet-v2":
                reachability = ""

        return node, node_version, os, arch, reachability
    else:
        return "Unknown", "Unknown", "Unknown", "Unknown", "Unknown"


def check_grpc_port(network_vals):
    for _, net_val in network_vals.items():
        ip = net_val["ip"]
        if is_port_open(ip, 50052):
            net_val["grpc_is_open"] = "Yes"
            channel = grpc.insecure_channel(ip + ":50052")
            stub = network_pb2_grpc.NetworkStub(channel)
            req = network_pb2.GetNodeInfoRequest()
            info = stub.GetNodeInfo(req)
            net_val["node_info"] += "agent: " + info.agent
            net_val["node_info"] += ", reachability: " + info.reachability
            net_val["node_info"] += ", addrs: " + str(list(info.addrs))

        if is_port_open(ip, 21777):
            net_val["p2p_is_open"] = "Yes"


def extract_network_vals(grpc_validator_stub):
    network_res = []
    network_res.append(get_network_info("95.217.181.93:50052"))  # 0.21.1
    network_res.append(get_network_info("104.194.156.57:50052"))  # 0.21.1
    network_res.append(get_network_info("5.75.169.99:50052"))  # 0.21.1
    network_res.append(get_network_info("157.90.153.208:50052"))  # 0.21.1
    network_res.append(get_network_info("81.0.218.193:50052"))  # 0.21.1
    network_res.append(get_network_info("139.99.198.111:50052"))  # 0.21.1
    network_res.append(get_network_info("49.13.126.38:50052"))  # 0.21.1

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
                    "grpc_is_open": "No",
                    "p2p_is_open": "No",
                    "moniker": "",
                    "address": "",
                    "peerId": "",
                    "agent": "",
                    "node_info": "",
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
            value["grpc_is_open"],
            value["p2p_is_open"],
            value["moniker"],
            value["address"],
            value["peerId"],
            value["agent"],
            value["node_info"],
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
                "gRPC Open",
                "P2P Open",
                "Moniker",
                "IP_Address",
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


grpc_validator_channel = grpc.insecure_channel("172.104.46.145:50052")
grpc_validator_stub = blockchain_pb2_grpc.BlockchainStub(grpc_validator_channel)

network_vals = extract_network_vals(grpc_validator_stub)

check_grpc_port(network_vals)
write_network_vals(network_vals)
