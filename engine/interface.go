package engine

type IEngine interface {
	Run(appID AppID, callerID string, inputs []string) (*CommandResult, error)
	Commands() []Command

	Stop()
	Start()
}
