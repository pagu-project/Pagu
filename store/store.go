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
	claimers map[string]*Claimer
	cfg      *config.Config
	logger   *log.SubLogger
}

func NewStore(cfg *config.Config, logger *log.SubLogger) (IStore, error) {
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

	ss := &Store{
		claimers: claimers,
		cfg:      cfg,
		logger:   logger,
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
	err := s.save()
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

func (s *Store) save() error {
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
