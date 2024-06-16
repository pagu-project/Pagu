package entity

type AppID int

const (
	AppIdCLI      AppID = 1
	AppIdDiscord  AppID = 2
	AppIdgRPC     AppID = 3
	AppIdHTTP     AppID = 4
	AppIdTelegram AppID = 5
)

func (appID AppID) String() string {
	switch appID {
	case AppIdCLI:
		return "CLI"
	case AppIdDiscord:
		return "Discord"
	case AppIdgRPC:
		return "gRPC"
	case AppIdHTTP:
		return "HTTP"
	case AppIdTelegram:
		return "Telegram"
	}

	return ""
}

func AllAppIDs() []AppID {
	return []AppID{
		AppIdCLI,
		AppIdDiscord,
		AppIdgRPC,
		AppIdHTTP,
		AppIdTelegram,
	}
}
