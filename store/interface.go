package store

type IStore interface {
	ClaimerInfo(testNetValAddr string) *Claimer
	AddClaimTransaction(testNetValAddr string, txID string) error
	ClaimStatus() *ClaimStatus

	SaveTwitterParty(party *TwitterParty) error
	FindTwitterParty(twitterName string) *TwitterParty

	WhitelistTwitterAccount(twitterID, twitterName, authorizedDiscordID string) error
	IsWhitelisted(twitterID string) bool
	BoosterStatus() *BoosterStatus
}
