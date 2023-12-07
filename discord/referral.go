package discord

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kehiy/RoboPac/config"
)

type Referral struct {
	ReferralCode string `json:"referral_code"`
	Points       int    `json:"points"`
	DiscordName  string `json:"discord_name"`
	DiscordID    string `json:"discord_id"`
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
func (rs *ReferralStore) NewReferral(discordId, discordName, referralCode string) error {
	rs.syncMap.Store(referralCode, &Referral{
		Points:       0,
		DiscordName:  discordName,
		ReferralCode: referralCode,
		DiscordID:    discordId,
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
func (rs *ReferralStore) GetData(code string) (*Referral, bool) {
	entry, found := rs.syncMap.Load(code)
	if !found {
		return nil, false
	}
	referral := entry.(*Referral)
	return referral, true
}

// GetAllReferrals retrieves all referrals in store.
func (rs *ReferralStore) GetAllReferrals() []*Referral {
	result := []*Referral{}

	rs.syncMap.Range(func(key, value any) bool {
		referral, _ := value.(*Referral)
		result = append(result, referral)
		if value == nil {
			return true
		}
		return false
	})

	return result
}

// AddPoint add one point for a referral.
func (rs *ReferralStore) AddPoint(code string) bool {
	entry, found := rs.syncMap.Load(code)
	if !found {
		return false
	}

	if found {
		referral := entry.(*Referral)
		referral.Points++
		rs.syncMap.Store(referral.ReferralCode, referral)
		return true
	}

	// save record
	data, err := marshalJSON(rs.syncMap)
	if err != nil {
		log.Printf("error marshalling validator data file: %v", err)
		return false
	}

	if err := os.WriteFile(rs.cfg.ReferralDataPath, data, 0o600); err != nil {
		log.Printf("failed to write to %s: %v", rs.cfg.ReferralDataPath, err)
		return false
	}

	return false
}
