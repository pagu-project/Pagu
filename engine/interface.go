package engine

type Engine interface {
	NetworkHealth([]string) (*NetHealthResponse, error)
	NetworkStatus([]string) (*NetStatus, error)
	NodeInfo([]string) (*NodeInfo, error)
	ClaimerInfo([]string) (string, error)
	Claim([]string) (string, error)

	Stop()
}
