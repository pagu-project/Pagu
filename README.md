# Pactus Universal Robot (pactus-bot)

## Compile
To compile the bot, on your terminal navigate to the project cmd folder and run this command: `GOOS=linux GOARCH=amd64 go build -o=./bin/pbot`

## Installation
1. Copy the compiled binary `pbot` from project bin folder e.g `./cmd/bin/pbot` as compiled above. You can copy it your installation folder, e.g `/opt/pbot`
2. In the installation folder, create a nother folder that will hold the bot data e.g `mkdir data`
3. In that data directory, add the three files namely `config.json`, `wallet.json`, and `validator.json`
4. Make sure that `pbot` binary is executable e.g by running this command: `chmod +x pbot`
5. Your folder structure should be like this:
        /opt/pbot
        ├── pbot
        ├── data
            ├── config.json
            ├── validator.json
            └── wallet.json
6. You can either create a bpot service and start the service by running `service pbot start` or open a new screen and run your `pbot` in terminal e.g `./pbot` and then exit the screen and leave the pbot service running

