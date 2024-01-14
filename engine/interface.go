package engine

type Engine interface {
	NetworkHealth([]string) (*NetHealthResponse, error)
	NetworkStatus([]string) (*NetStatus, error)
	NodeInfo([]string) (*NodeInfo, error)
	MyInfo([]string) (string, error)
	Withdraw([]string) (string, error)

	Stop()
}
