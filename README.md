# Pagu
Pagu is a Bot engine that provides support and information about the [Pactus](https://pactus.org) Blockchain.

<p align="center">
    <img alt="Pagu" src="./assets/PAGU.png" width="150" height="150" />
</p>

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Features

- Pactus network health, status and statistics
- Pactus fee and reward calculation
- PAC coin market prices
- Phoenix (Pactus testnet) health and status 
- Phoenix faucet
  
## Installation

To get started with Pagu, follow these steps:

1. **Clone the repository**:
    ```bash
    git clone https://github.com/pagu-project/Pagu.git
    cd Pagu
    ```

2. **Install dependencies**:

   Install [Go](https://go.dev/doc/install) if you have not installed before and also install [Mysql](https://dev.mysql.com/downloads/workbench/) as main database of Pagu
then run below commands
   
    ```bash 
       make devtools
       cp ./config/config-sample.yml ./config/config.yml # fill the file with correct values
    ```
3. **Run Local Nodes**
   
   To run local node and set thier address in config file please follow below instruction

   https://docs.pactus.org/get-started/pactus-daemon/


4. **Wallet requirements**:
   
   Pagu needs a Pactus wallet to call transaction methods. If you have no wallet follow bellow instruction to make one.
 
   https://docs.pactus.org/tutorials/pactus-wallet/#create-a-wallet


5. **Discord Server**:
   
   To run Pagu in Discord server you need a GuildID of server and discord application token. To make them please follow below link

   https://discord.com/developers/docs/quick-start/getting-started
   
## Run
1. **Discord Engine**

    ```bash
    make check
    go run ./cmd/discord -c ./config/config.main.yml run
    ```

## Contributing

We welcome contributions! If you'd like to contribute to Pagu, please follow these guidelines:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add some feature'`).
5. Push to the branch (`git push origin feature/YourFeature`).
6. Open a Pull Request.

---