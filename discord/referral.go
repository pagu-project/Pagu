package discord

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kehiy/RoboPac/config"
)

type Referral struct {
	ReferralCode   string `json:"referral_code"`
	AccountAddress string `json:"account_address"`
	ReferralCounts int    `json:"referral_count"`
	DiscordName    string `json:"discord_name"`
	DiscordID      string `json:"discord_id"`
}

// SafeStore is a thread-safe cache.
type ReferralStore struct {
	syncMap *sync.Map
	cfg     *config.Config
}

func LoadReferralData(cfg *config.Config) (*ReferralStore, error) {
	file, err := os.ReadFile(cfg.ReferralDataPath)
	if err != nil {
		log.Printf("error loading validator data: %v", err)
		return nil, fmt.Errorf("error loading data file: %w", err)
	}
	if len(file) == 0 {
		rs := &ReferralStore{
			syncMap: &sync.Map{},
			cfg:     cfg,
		}
		return rs, nil
	}

	data, err := unmarshalJSON(file)
	if err != nil {
		log.Printf("error unmarshalling validator data: %v", err)
		return nil, fmt.Errorf("error unmarshalling validator data: %w", err)
	}
	rs := &ReferralStore{
		syncMap: data,
		cfg:     cfg,
	}
	return rs, nil
}

// SetData Set a given value to the data storage.
func (rs *ReferralStore) SetData(address string, count int) error {
	rs.syncMap.Store(address, &Referral{
		ReferralCounts: count,
	})
	// save record
	data, err := marshalJSON(rs.syncMap)
	if err != nil {
		log.Printf("error marshalling validator data file: %v", err)
		return fmt.Errorf("error marshalling validator data file: %w", err)
	}
	if err := os.WriteFile(rs.cfg.ReferralDataPath, data, 0o600); err != nil {
		log.Printf("failed to write to %s: %v", rs.cfg.ReferralDataPath, err)
		return fmt.Errorf("failed to write to %s: %w", rs.cfg.ReferralDataPath, err)
	}
	return nil
}

// GetData retrieves the given key from the storage.
func (rs *ReferralStore) GetData(address string) (*Referral, bool) {
	entry, found := rs.syncMap.Load(address)
	if !found {
		return nil, false
	}
	referral := entry.(*Referral)
	return referral, true
}
