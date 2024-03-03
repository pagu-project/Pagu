package engine

type AppID int

type Command struct {
	Name    string
	Desc    string
	Help    string
	Args    []Args
	AppIDs  []AppID
	Handler func(source AppID, callerID string, args ...string) (*CommandResult, error)
}

type Args struct {
	Name     string
	Desc     string
	Optional bool
}

const (
	AppIdCLI     AppID = 1
	AppIdDiscord AppID = 2
)

func (be *BotEngine) RegisterCommands() {
	HelpCmd := Command{
		Name:   "Help",
		Desc:   "Help command.",
		Help:   "This is the help command.",
		Args:   []Args{},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	NetworkHealthCmd := Command{
		Name:   "NetworkHealth",
		Desc:   "Command to check network health.",
		Help:   "Network health.",
		Args:   []Args{},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	NetworkStatusCmd := Command{
		Name:   "NetworkStatus",
		Desc:   "Command to check the status of pactus chain.",
		Help:   "check the status of the Pactus network.",
		Args:   []Args{},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	NodeInfoCmd := Command{
		Name: "NodeInfo",
		Desc: "Command to see information on your node.",
		Help: "This command will help you check info on your node",
		Args: []Args{
			{
				Name:     "nodeID",
				Desc:     "ID of the node to get information about.",
				Optional: false,
			},
		},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	RewardCalculateCmd := Command{
		Name: "RewardCalculate",
		Desc: "Command to calculate your potential staking rewards.",
		Help: "This command will help you calculate your potential staking rewards.",
		Args: []Args{
			{
				Name:     "stake",
				Desc:     "Your validator stake amount",
				Optional: false,
			},
			{
				Name:     "time",
				Desc:     "In a day/month/year",
				Optional: false,
			},
		},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	ClaimerInfoCmd := Command{
		Name: "ClaimerInfo",
		Desc: "Get claimer info.",
		Help: "Command to fetch claimer info.",
		Args: []Args{
			{
				Name:     "claimer-info",
				Desc:     "Get claimer info",
				Optional: false,
			},
		},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	ClaimCmd := Command{
		Name: "Claim",
		Desc: "Claim your Pactus coins.",
		Help: "claim your Pactus coins.",
		Args: []Args{
			{
				Name:     "testnet-addr",
				Desc:     "Enter your testnet address.",
				Optional: false,
			},
			{
				Name:     "mainnet-addr",
				Desc:     "Enter your mainnet address",
				Optional: false,
			},
		},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	ClaimStatusCmd := Command{
		Name:   "ClaimStatus",
		Desc:   "Testnet reward claim status",
		Help:   "Claim status",
		Args:   []Args{},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	BotWalletCmd := Command{
		Name:   "BotWallet",
		Desc:   "Bot wallet balance.",
		Help:   "Bot wallet balance",
		Args:   []Args{},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	BoosterWhitelistCmd := Command{
		Name: "BoosterWhitelist",
		Desc: "Whitelist a non-active Twitter account in Validator Booster Program",
		Help: "Booster whitelist",
		Args: []Args{
			{
				Name:     "twitter-username",
				Desc:     "Twitter username",
				Optional: false,
			},
		},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	BoosterClaimCmd := Command{
		Name: "BoosterClaim",
		Desc: "Claim the stake PAC coin in Validator Booster Program",
		Help: "your Twitter username",
		Args: []Args{
			{
				Name:     "twitter-username",
				Desc:     "your Twitter username",
				Optional: false,
			},
		},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	BoosterPaymentCmd := Command{
		Name: "BoosterPayment",
		Desc: "Create payment link in Validator Booster Program",
		Help: "Create payment link in Validator Booster Program",
		Args: []Args{
			{
				Name:     "twitter-username",
				Desc:     "your Twitter username",
				Optional: false,
			},
			{
				Name:     "validator-address",
				Desc:     "your validator address",
				Optional: false,
			},
		},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	BoosterStatusCmd := Command{
		Name:   "BoosterStatus",
		Desc:   "Validator Booster Program Status",
		Help:   "Booster status",
		Args:   []Args{},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	DepositAddressCmd := Command{
		Name:   "deposit-address",
		Desc:   "create a deposit address or get your deposit address",
		Help:   "",
		Args:   []Args{},
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	CreateOfferCmd := Command{
		Name: "create",
		Desc: "create an offer",
		Help: "",
		Args: []Args{
			{
				Name:     "total-amount",
				Desc:     "total amount of PAC",
				Optional: false,
			},
			{
				Name:     "total-price",
				Desc:     "total price which includes gas fee",
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
		AppIDs: []AppID{AppIdCLI, AppIdDiscord},
	}

	be.Cmds = append(be.Cmds,
		HelpCmd,
		NetworkHealthCmd,
		NetworkStatusCmd,
		NodeInfoCmd,
		RewardCalculateCmd,
		ClaimerInfoCmd,
		ClaimCmd,
		ClaimStatusCmd,
		BotWalletCmd,
		BoosterWhitelistCmd,
		BoosterClaimCmd,
		BoosterPaymentCmd,
		BoosterStatusCmd,
		DepositAddressCmd,
		CreateOfferCmd,
	)
}

func (be *BotEngine) Commands() []Command {
	return be.Cmds
}
