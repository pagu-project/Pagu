# Installing and Running Pagu

This document provides detailed instructions on how to install and run Pagu for development purposes.

## Prerequisites

Ensure you have the following installed on your system before proceeding:

- **Go**: Pagu is developed using the Go programming language. You can find installation instructions [here](https://go.dev/doc/install).
- **MySQL**: Pagu uses MySQL as its primary database. Download and install MySQL from [here](https://dev.mysql.com/downloads/workbench/).

## Installation Steps

Follow the steps below to install and configure Pagu on your local machine.

### 1. Clone the Repository

First, clone the Pagu repository to your local machine:

```bash
git clone https://github.com/pagu-project/Pagu.git
cd Pagu
```

### 2. Install Development Tools

Install the necessary development tools by running the following commands:

```bash
make devtools
cp ./config/config-sample.yml ./config/config.yml
```

Modify `config.yml` with the appropriate values for your setup.

### 3. Running Local Pactus Nodes

To run local nodes and configure them in your `config.yml`, refer to the [Pactus Daemon documentation](https://docs.pactus.org/get-started/pactus-daemon/).

### 4. Wallet Requirements

Pagu requires a Pactus wallet to manage transactions. If you don’t have a wallet, follow the instructions to create one [here](https://docs.pactus.org/tutorials/pactus-wallet/#create-a-wallet).

### 5. Discord Setup

If you plan to run Pagu in a Discord server, you will need a Guild ID and a Discord application token. You can obtain these by following the [Discord Developer Guide](https://discord.com/developers/docs/quick-start/getting-started).

## Running Pagu

Once the installation and configuration are complete, you can run Pagu using the following commands.

### 1. Discord Engine

Run the Discord engine using:

```bash
make check
go run ./cmd/discord -c ./config/config.main.yml run
```

## Contributing

We welcome contributions to Pagu! Please follow the steps below to get started:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add some feature'`).
5. Push to the branch (`git push origin feature/YourFeature`).
6. Open a Pull Request.
