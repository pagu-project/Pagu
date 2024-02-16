package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kehiy/RoboPac/log"
	"github.com/pactus-project/pactus/util/logger"
)

// Store is a thread-safe cache.
type Store struct {
	claimers             map[string]*Claimer
	twitterParties       map[string]*TwitterParty
	twitterWhitelisted   map[string]*WhitelistInfo
	claimersPath         string
	twitterPartiesPath   string
	twitterWhitelistPath string
	logger               *log.SubLogger
}

func loadMap[T any](path string, mapObj map[string]*T) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error loading data file: %w", err)
	}
	if len(data) == 0 {
		return fmt.Errorf("empty file: %s", path)
	}

	if err := json.Unmarshal(data, &mapObj); err != nil {
		return err
	}

	return nil
}

func saveMap[T any](path string, mapObj map[string]*T) error {
	logger.Debug("save map", "path", path)

	data, err := json.Marshal(mapObj)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func NewStore(storePath string, logger *log.SubLogger) (IStore, error) {
	claimers := make(map[string]*Claimer)
	twitterParties := make(map[string]*TwitterParty)
	twitterWhitelisted := make(map[string]*WhitelistInfo)

	claimersPath := path.Join(storePath, "claimers.json")
	twitterPartiesPath := path.Join(storePath, "twitter_campaign.json")
	twitterWhitelistPath := path.Join(storePath, "twitter_whitelisted.json")

	err := loadMap(claimersPath, claimers)
	if err != nil {
		return nil, err
	}

	err = loadMap(twitterPartiesPath, twitterParties)
	if err != nil {
		return nil, err
	}

	err = loadMap(twitterWhitelistPath, twitterWhitelisted)
	if err != nil {
		return nil, err
	}

	ss := &Store{
		claimers:             claimers,
		twitterParties:       twitterParties,
		twitterWhitelisted:   twitterWhitelisted,
		claimersPath:         claimersPath,
		twitterPartiesPath:   twitterPartiesPath,
		twitterWhitelistPath: twitterWhitelistPath,
		logger:               logger,
	}
	return ss, nil
}

func (s *Store) ClaimerInfo(testnetAddr string) *Claimer {
	entry, found := s.claimers[testnetAddr]
	if !found {
		return nil
	}

	return entry
}

func (s *Store) AddClaimTransaction(testnetAddr string, txID string) error {
	entry, found := s.claimers[testnetAddr]
	if !found {
		return fmt.Errorf("testnetAddr not found: %s", testnetAddr)
	}

	entry.ClaimedTxID = txID
	err := s.saveClaimers()
	if err != nil {
		return err
	}

	s.logger.Info("new claim transaction added",
		"discordID", entry.DiscordID,
		"amount", entry.TotalReward,
		"txID", txID)

	return nil
}

func (s *Store) ClaimStatus() (int64, int64, int64, int64) {
	var claimed int64
	var claimedAmount int64

	var notClaimed int64
	var notClaimedAmount int64

	for _, c := range s.claimers {
		if c.IsClaimed() {
			claimed++
			claimedAmount += c.TotalReward
		} else {
			notClaimed++
			notClaimedAmount += c.TotalReward
		}
	}
	return claimed, claimedAmount, notClaimed, notClaimedAmount
}

func (s *Store) saveClaimers() error {
	return saveMap(s.claimersPath, s.claimers)
}

func (s *Store) saveTwitterParties() error {
	return saveMap(s.twitterPartiesPath, s.twitterParties)
}

func (s *Store) saveTwitterWhitelist() error {
	return saveMap(s.twitterWhitelistPath, s.twitterWhitelisted)
}

func (s *Store) SaveTwitterParty(party *TwitterParty) error {
	s.twitterParties[party.TwitterID] = party

	return s.saveTwitterParties()
}

func (s *Store) FindTwitterParty(twitterName string) *TwitterParty {
	for _, party := range s.twitterParties {
		if strings.EqualFold(party.TwitterName, twitterName) {
			return party
		}
	}
	return nil
}

func (s *Store) WhitelistTwitterAccount(twitterID, twitterName, authorizedDiscordID string) error {
	_, exists := s.twitterWhitelisted[twitterID]
	if exists {
		return fmt.Errorf("the Twitter `%v` is already whitelisted", twitterName)
	}

	s.twitterWhitelisted[twitterID] = &WhitelistInfo{
		TwitterID:     twitterID,
		TwitterName:   twitterName,
		WhitelistedBy: authorizedDiscordID,
	}

	return s.saveTwitterWhitelist()
}

func (s *Store) IsWhitelisted(twitterID string) bool {
	_, exists := s.twitterWhitelisted[twitterID]

	return exists
}
