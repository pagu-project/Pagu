package engine

import "github.com/kehiy/RoboPac/store"

type IEngine interface {
	NetworkHealth() (*NetHealthResponse, error)
	NetworkStatus() (*NetStatus, error)
	NodeInfo(addr string) (*NodeInfo, error)
	ClaimerInfo(discordID string) (*store.Claimer, error)
	Claim(discordID string, testnetAddr string, mainnetAddr string) (string, error)
	BotWallet() (string, int64)
	RewardCalculate(int64, string) (int64, string, int64, error)

	Run(input string) (string, error)

	Stop()
	Start()
}
