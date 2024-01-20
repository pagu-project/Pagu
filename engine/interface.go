package engine

import "github.com/kehiy/RoboPac/store"

type Engine interface {
	NetworkHealth([]string) (*NetHealthResponse, error)
	NetworkStatus([]string) (*NetStatus, error)
	NodeInfo([]string) (*NodeInfo, error)
	ClaimerInfo([]string) (*store.Claimer, error)
	Claim([]string) (*store.ClaimTransaction, error)

	Stop()
	Start()
}
