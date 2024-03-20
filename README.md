<p align="center">
    <img alt="RoboPac" src="./assets/robopac.png" width="150" height="150" />
</p>

<h3 align="center">
RoboPac is a Robot that provides support and information about the Pactus Blockchain.
</h3>

### Pactus Bot Engine (RoboPac)


# Run

The RoboPac is require golang installed to be run. make sure you installed The Pactus daemon CLI from here:
https://pactus.org/download/#cli

You need to run 2 local pactus node (local-net) and add them to a file called local (it's on .gitignore) and run them. 

commands:

```pactus-daemon init -w=./local/net1 --localnet```

```pactus-daemon init -w=./local/net2 --localnet```

> Note: make sure you enable the gRPC for them.

Then you can make a .env file by following .env.example and fill it with you own node data and test wallet.

> Note2: you can make test-net wallets like: `pactus-wallet create --testnet`


Last step is to run `make build` and use the robopac-cli binary to start testing your new feature or command.

## Contributing

Contributions to the RoboPac are appreciated.

## License

RoboPac it under [MIT](./LICENSE).
