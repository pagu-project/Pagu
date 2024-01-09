package discord

type Result struct {
	Agent            string  `json:"agent"`
	ValidatorAddress string  `json:"validator_address"`
	PIP19Score       float64 `json:"pip19_score"`
	IsActive         bool    `json:"is_active"`
	RemoteAddress    string  `json:"remote_address"`
}
