package engine

import "github.com/kehiy/RoboPac/store"

type IEngine interface {
	NetworkHealth() (*NetHealthResponse, error)
	NetworkStatus() (*NetStatus, error)
	NodeInfo(addr string) (*NodeInfo, error)
	RewardCalculate(int64, string) (int64, string, int64, error)

	ClaimerInfo(discordID string) (*store.Claimer, error)
	Claim(discordID string, testnetAddr string, mainnetAddr string) (string, error)
	ClaimStatus() *store.ClaimStatus

	BotWallet() (string, int64)

	BoosterWhitelist(string, string) error
	BoosterClaim(string) (*store.TwitterParty, error)
	BoosterPayment(string, string, string) (*store.TwitterParty, error)
	BoosterStatus() *store.BoosterStatus

	Run(input string) (string, error)

	Stop()
	Start()
}

type CommandHandler interface {
	Execute(source AppID, callerID string, args ...string) (*CommandResult, error)
}
