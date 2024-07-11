package entity

type AppID int

const (
	AppIDCLI      AppID = 1
	AppIDDiscord  AppID = 2
	AppIDgRPC     AppID = 3
	AppIDHTTP     AppID = 4
	AppIDTelegram AppID = 5
)

func (appID AppID) String() string {
	switch appID {
	case AppIDCLI:
		return "CLI"
	case AppIDDiscord:
		return "Discord"
	case AppIDgRPC:
		return "gRPC"
	case AppIDHTTP:
		return "HTTP"
	case AppIDTelegram:
		return "Telegram"
	}

	return ""
}

func AllAppIDs() []AppID {
	return []AppID{
		AppIDCLI,
		AppIDDiscord,
		AppIDgRPC,
		AppIDHTTP,
		AppIDTelegram,
	}
}
