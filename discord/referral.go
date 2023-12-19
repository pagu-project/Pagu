package discord

import (
	"encoding/json"
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
		log.Printf("error loading referral data: %v", err)
		return nil, fmt.Errorf("error loading data file: %w", err)
	}
	if len(file) == 0 {
		rs := &ReferralStore{
			syncMap: &sync.Map{},
			cfg:     cfg,
		}
		return rs, nil
	}

	data, err := unmarshalReferralJSON(file)
	if err != nil {
		log.Printf("error unmarshalling referral data: %v", err)
		return nil, fmt.Errorf("error unmarshalling referral data: %w", err)
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
	data, err := marshaReferralJSON(rs.syncMap)
	if err != nil {
		log.Printf("error marshalling referral data file: %v", err)
		return fmt.Errorf("error marshalling referral data file: %w", err)
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
		return true
	})

	return result
}

// AddPoint add one point for a referral.
func (rs *ReferralStore) AddPoint(code string) bool {
	entry, found := rs.syncMap.Load(code)
	if !found {
		return false
	}

	// updating record
	referral := entry.(*Referral)
	referral.Points++
	rs.syncMap.Store(referral.ReferralCode, referral)

	// saving record
	data, err := marshaReferralJSON(rs.syncMap)
	if err != nil {
		log.Printf("error marshalling referral data file: %v", err)
		return false
	}

	if err := os.WriteFile(rs.cfg.ReferralDataPath, data, 0o600); err != nil {
		log.Printf("failed to write to %s: %v", rs.cfg.ReferralDataPath, err)
		return false
	}

	return true
}

func marshaReferralJSON(m *sync.Map) ([]byte, error) {
	tmpMap := make(map[string]*Referral)

	m.Range(func(k, v interface{}) bool {
		tmpMap[k.(string)] = v.(*Referral)
		return true
	})
	return json.MarshalIndent(tmpMap, "  ", "  ")
}

func unmarshalReferralJSON(data []byte) (*sync.Map, error) {
	var tmpMap map[string]*Referral
	m := &sync.Map{}

	if err := json.Unmarshal(data, &tmpMap); err != nil {
		return m, err
	}

	for key, value := range tmpMap {
		m.Store(key, value)
	}
	return m, nil
}
