package p2pmarket

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/database"
	"github.com/kehiy/RoboPac/engine/command"
	"github.com/kehiy/RoboPac/wallet"
	"github.com/pactus-project/pactus/util"
)

const (
	P2PCommandName            = "p2p"
	DepositAddressCommandName = "deposit-address"
	CreateOfferCommandName    = "create-offer"
	P2PHelpCommandName        = "help"
)

type P2PMarket struct {
	ctx       context.Context
	AdminIDs  []string
	db        *database.DB
	wallet    wallet.IWallet
	clientMgr *client.Mgr
}

func NewP2PMarket(ctx context.Context,
	adminIDs []string,
	db database.DB,
	wallet wallet.IWallet,
	clientMgr *client.Mgr,
) *P2PMarket {
	return &P2PMarket{
		ctx:       ctx,
		AdminIDs:  adminIDs,
		db:        &db,
		wallet:    wallet,
		clientMgr: clientMgr,
	}
}

func (be *P2PMarket) GetCommand() *command.Command {
	subCmdDepositAddress := command.Command{
		Name:        DepositAddressCommandName,
		Desc:        "Create a deposit address for P2P offer",
		Help:        "It will show your address if you already have an deposit address",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		Handler:     be.depositAddressHandler,
	}

	subCmdCreateOffer := command.Command{
		Name: CreateOfferCommandName,
		Desc: "Create an offer for P2P market",
		Help: "",
		Args: []command.Args{
			{
				Name:     "total-amount",
				Desc:     "Total amount of PAC",
				Optional: false,
			},
			{
				Name:     "total-price",
				Desc:     "Total price which includes gas fee",
				Optional: false,
			},
			{
				Name:     "chain-type",
				Desc:     "e.g. BTCUSDT",
				Optional: false,
			},
			{
				Name:     "address",
				Desc:     "",
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		Handler:     be.createOfferHandler,
	}

	cmdP2PMarket := command.Command{
		Name:        P2PCommandName,
		Desc:        "Peer to Peer market for Pactus trading",
		Help:        "",
		Args:        nil,
		AppIDs:      []command.AppID{command.AppIdCLI, command.AppIdDiscord},
		SubCommands: []*command.Command{&subCmdCreateOffer, &subCmdDepositAddress},
		Handler:     nil,
	}

	cmdP2PMarket.AddSubCommand(&subCmdDepositAddress)
	cmdP2PMarket.AddSubCommand(&subCmdCreateOffer)

	cmdP2PMarket.AddHelpSubCommand()

	return &cmdP2PMarket
}

func (bpm *P2PMarket) depositAddressHandler(cmd *command.Command, _ command.AppID, callerID string, _ ...string) *command.CommandResult {
	u, err := bpm.db.GetUser(callerID)
	if err == nil {
		return &command.CommandResult{
			Successful: true,
			Message:    fmt.Sprintf("You already have a deposit address: %s", u.DepositAddress),
		}
	}

	addr, err := bpm.wallet.NewAddress(fmt.Sprintf("deposit address for %s", callerID))
	if err != nil {
		return &command.CommandResult{
			Successful: false,
			Message:    fmt.Sprintf("Can't make a new address: %v", err),
		}
	}

	err = bpm.db.AddUser(
		&database.DiscordUser{
			DiscordID:      callerID,
			DepositAddress: addr,
		},
	)
	if err != nil {
		return &command.CommandResult{
			Successful: false,
			Message:    fmt.Sprintf("Can't add discord user to database: %v", err),
		}
	}

	return &command.CommandResult{
		Successful: true,
		Message:    fmt.Sprintf("Deposit address created for you successfully: %s", addr),
	}
}

func (pm *P2PMarket) createOfferHandler(cmd *command.Command, source command.AppID, callerID string, args ...string) *command.CommandResult {
	u, err := pm.db.GetUser(callerID)
	if err != nil {
		return &command.CommandResult{
			Successful: false,
			Error:      err.Error(),
		}
	}

	totalAmount, err := strconv.Atoi(args[0])
	if err != nil {
		return &command.CommandResult{
			Successful: false,
			Error:      err.Error(),
		}
	}

	totalPrice, err := strconv.Atoi(args[1])
	if err != nil {
		return &command.CommandResult{
			Successful: false,
			Error:      err.Error(),
		}
	}

	chainType := args[2]
	address := args[3]

	uBalance, err := pm.clientMgr.GetBalance(u.DepositAddress)
	if err != nil {
		return &command.CommandResult{
			Successful: false,
			Error:      err.Error(),
		}
	}

	if float64(totalAmount) != util.ChangeToCoin(uBalance) {
		return &command.CommandResult{
			Successful: false,
			Error: fmt.Sprintf("the deposit balance: %d is not equal to offered amount: %d",
				uBalance, totalAmount),
		}
	}

	unitPrice := float64(totalPrice / totalAmount)

	offer := &database.Offer{
		TotalAmount: int64(totalAmount),
		TotalPrice:  int64(totalPrice),
		UnitPrice:   unitPrice,
		ChainType:   chainType,
		Address:     address,
		DiscordUser: *u,
	}

	if err = pm.db.CreateOffer(offer); err != nil {
		return &command.CommandResult{
			Successful: false,
			Error:      err.Error(),
		}
	}

	return &command.CommandResult{
		Successful: true,
		Message:    fmt.Sprintf("Offer successfully created, your offer ID: %s", "TODO!!!!!!!"),
	}
}
