import re

def extract_info(input_string):
    # Regular expression pattern to extract desired information
    pattern = r'node=([^/]+)/node-version=([^/]+)/.*?/os=([^/]+)/arch=([^/]+)(/reachability=([^/]+))?'
    match = re.match(pattern, input_string)

    if match:
        node = match.group(1)
        node_version = match.group(2)
        os = match.group(3)
        arch = match.group(4)
        reachability = "Unknown"
        if match.group(5) is not None:
            reachability = match.group(5)
        return {
            "Node": node,
            "Node Version": node_version,
            "OS": os,
            "Arch": arch,
            "Reachability": reachability
        }
    else:
        print("Invalid input string:", input_string)
        return None

# Sample input strings
strings = [
    "node=pactus-daemon/node-version=v0.20.0/protocol-version=1/os=linux/arch=amd64/reachability=pactus-testnet-v2",
    "node=pactus-gui.exe/node-version=v0.18.4/protocol-version=1/os=windows/arch=amd64",
    "node=pactus-daemon/node-version=v0.20.1/protocol-version=1/os=linux/arch=amd64/reachability=Private",
]

for s in strings:
    result = extract_info(s)
    if result:
        print(result)
