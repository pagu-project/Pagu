import grpc
import sys
import json
# Why can't??
# from pactus import network_pb2
# from pactus import network_pb2_grpc

import public_key
import network_pb2
import network_pb2_grpc


def load_json_file(file_path):
    try:
        with open(file_path, 'r') as file:
            data = json.load(file)
        return data
    except FileNotFoundError:
        print(f"File not found: {file_path}")
        return None
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")
        return None


def update_validator_map(url, valMap):
    channel = grpc.insecure_channel(url)
    stub = network_pb2_grpc.NetworkStub(channel)
    req = network_pb2.GetNetworkInfoRequest()
    info = stub.GetNetworkInfo(req)

    for peer in info.peers:
        for i, key in enumerate(peer.consensus_keys):
            pub = public_key.PublicKey.from_string(key)
            valAddr = pub.validator_address().string()
            valMap[valAddr] = i

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: python main.py <file_path>")
    else:
        file_path = sys.argv[1]
        json_data = load_json_file(file_path)

    valMap = {}
    update_validator_map('172.104.46.145:50052', valMap)
    update_validator_map('94.101.184.118:50052', valMap)
    update_validator_map('51.158.118.181:50052', valMap)
    update_validator_map('172.232.108.191:50052', valMap)

    for key, value in json_data.items():
        userValAddr = value['validator_address']
        index = valMap.get(userValAddr, None)
        if index is not None:
            if index != 0:
                print("user {} staked in wrong validator. index: {}".format(
                    value['discord_name'], index))
        else:
            print("unable to find validator {} information for user {}".format(
                value['validator_address'],
                value['discord_name']))
