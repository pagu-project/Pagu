package engine

type Engine interface {
	NetworkHealth([]string) (string, error)
	NetworkData([]string) (string, error)
	MyInfo([]string) (string, error)
	Withdraw([]string) (string, error)
	NodeInfo([]string) (string, error)

	Stop()
}
