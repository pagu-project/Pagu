package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/log"
)

// Store is a thread-safe cache.
type Store struct {
	claimers       map[string]*Claimer
	twitterParties map[string]*TwitterParty
	cfg            *config.Config
	logger         *log.SubLogger
}

func NewStore(cfg *config.Config, logger *log.SubLogger) (IStore, error) {
	loadClaimers := func() (map[string]*Claimer, error) {
		data, err := os.ReadFile(cfg.StorePath)
		if err != nil {
			return nil, fmt.Errorf("error loading data file: %w", err)
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("empty file: %s", cfg.StorePath)
		}

		claimers := make(map[string]*Claimer)
		if err := json.Unmarshal(data, &claimers); err != nil {
			return nil, err
		}

		return claimers, nil
	}

	loadTwitterParties := func() (map[string]*TwitterParty, error) {
		data, err := os.ReadFile(cfg.TwitterStorePath)
		if err != nil {
			return nil, fmt.Errorf("error loading data file: %w", err)
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("empty file: %s", cfg.TwitterStorePath)
		}

		parties := make(map[string]*TwitterParty)
		if err := json.Unmarshal(data, &parties); err != nil {
			return nil, err
		}

		return parties, nil
	}

	claimers, err := loadClaimers()
	if err != nil {
		return nil, err
	}
	twitterParties, err := loadTwitterParties()
	if err != nil {
		return nil, err
	}

	ss := &Store{
		claimers:       claimers,
		twitterParties: twitterParties,
		cfg:            cfg,
		logger:         logger,
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

func (s *Store) Status() (int64, int64, int64, int64) {
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
	data, err := json.Marshal(s.claimers)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.cfg.StorePath, data, 0o600)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) saveTwitterParties() error {
	data, err := json.Marshal(s.twitterParties)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.cfg.TwitterStorePath, data, 0o600)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) AddTwitterParty(party *TwitterParty) error {
	found, exists := s.twitterParties[party.TwitterName]
	if exists {
		return fmt.Errorf("the Twitter `%v` already registered for the campagna. Discount code is %v",
			found.TwitterName, party.DiscountCode)
	}

	s.twitterParties[party.TwitterName] = party

	err := s.saveTwitterParties()
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) GetTwitterParty(twitterName string) *TwitterParty {
	return s.twitterParties[twitterName]
}
