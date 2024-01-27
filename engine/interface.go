package engine

type IEngine interface {
	// NetworkHealth() (*NetHealthResponse, error)
	// NetworkStatus([]string) (*NetStatus, error)
	// NodeInfo(addr string) (*NodeInfo, error)
	// ClaimerInfo([]string) (*store.Claimer, error)
	// Claim(discordID string, testnetAddr string, mainnetAddr string) (*store.ClaimTransaction, error)

	Run(input string) (string, error)

	Stop()
	Start()
}
