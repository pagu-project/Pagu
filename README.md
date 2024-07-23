<p align="center">
    <img alt="Pagu" src="./assets/PAGU.png" width="150" height="150" />
</p>

<h3 align="center">
The PAGU is a Robot that provides support and information about the Pactus Blockchain.
</h3>

# Run

The Pagu is required golang installed to be run. make sure you installed The Pactus daemon CLI (and wallet CLI) from here:
https://pactus.org/download/#cli

You need to run 2 local pactus node (local-net) and add them to a file called local (it's on .gitignore) and run them. 

commands:

```pactus-daemon init -w=./local/net1 --localnet```

```pactus-daemon init -w=./local/net2 --localnet```

> Note: make sure you enable the gRPC for them.

Make sure you run a postgres instance using docker on your local machine and make a proper config based on [this](./config/), also you can find deployment info and guidelines [here](./deployment/) to run your local instance on Pagu using docker compose.

> Note2: you can make test-net wallets like: `pactus-wallet create --testnet`


Last step is to run `make build` and use the pagu-cli binary to start testing your new feature or command.

## Assets

The Pagu logo and other assets are available on [here](./assets/) on this repo for usage.

## Contributing

Contributions to the Pagu are appreciated.

## License

The Pagu it under [MIT](./LICENSE).
