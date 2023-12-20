package discord

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kehiy/RoboPac/config"
)

// Validator is a value stored in the cache.
type Validator struct {
	DiscordName       string  `json:"discord_name"`
	DiscordID         string  `json:"discord_id"`
	ValidatorAddress  string  `json:"validator_address"`
	ReferrerDiscordID string  `json:"referrer_discord_id"`
	FaucetAmount      float64 `json:"faucet_amount"`
}

type ValidatorRoundOne struct {
	ID              int    `json:"id"`
	Address         string `json:"address"`
	DiscordUsername string `json:"discordUsername"`
	DiscordID       string `json:"discord_id"`
	Status          string `json:"status"`
	Twitter         string `json:"Twitter"`
	Total           int    `json:"Total"`
}

// SafeStore is a thread-safe cache.
type SafeStore struct {
	syncMap *sync.Map
	cfg     *config.Config
}

func LoadData(cfg *config.Config) (*SafeStore, error) {
	file, err := os.ReadFile(cfg.ValidatorDataPath)
	if err != nil {
		log.Printf("error loading validator data: %v", err)
		return nil, fmt.Errorf("error loading data file: %w", err)
	}
	if len(file) == 0 {
		ss := &SafeStore{
			syncMap: &sync.Map{},
			cfg:     cfg,
		}
		return ss, nil
	}

	data, err := unmarshalJSON(file)
	if err != nil {
		log.Printf("error unmarshalling validator data: %v", err)
		return nil, fmt.Errorf("error unmarshalling validator data: %w", err)
	}
	ss := &SafeStore{
		syncMap: data,
		cfg:     cfg,
	}
	return ss, nil
}

// SetData Set a given value to the data storage.
func (ss *SafeStore) SetData(peerID, address, discordName, discordID, referrerDiscordID string, amount float64) error {
	ss.syncMap.Store(peerID, &Validator{
		DiscordName: discordName, DiscordID: discordID,
		ValidatorAddress: address, FaucetAmount: amount,
		ReferrerDiscordID: referrerDiscordID,
	})
	// save record
	data, err := marshalJSON(ss.syncMap)
	if err != nil {
		log.Printf("error marshalling validator data file: %v", err)
		return fmt.Errorf("error marshalling validator data file: %w", err)
	}
	if err := os.WriteFile(ss.cfg.ValidatorDataPath, data, 0o600); err != nil {
		log.Printf("failed to write to %s: %v", ss.cfg.ValidatorDataPath, err)
		return fmt.Errorf("failed to write to %s: %w", ss.cfg.ValidatorDataPath, err)
	}
	return nil
}

// GetData retrieves the given key from the storage.
func (ss *SafeStore) GetData(peerID string) (*Validator, bool) {
	entry, found := ss.syncMap.Load(peerID)
	if !found {
		return nil, false
	}
	validator := entry.(*Validator)
	return validator, true
}

func (ss *SafeStore) FindDiscordID(discordID string) (*Validator, bool) {
	validator := &Validator{}
	exists := false

	ss.syncMap.Range(func(key, value any) bool {
		v := value.(*Validator)
		if validator.DiscordID == discordID {
			validator = v
			exists = true
		}
		return true
	})
	return validator, exists
}

func (ss *SafeStore) GetDistribution() (uint, float64) {
	totalDistribution := float64(0)
	totalValidators := uint(0)

	ss.syncMap.Range(func(key, value any) bool {
		v := value.(*Validator)
		if v != nil {
			totalDistribution += v.FaucetAmount
			totalValidators++
		}
		return true
	})
	return totalValidators, totalDistribution
}

func marshalJSON(m *sync.Map) ([]byte, error) {
	tmpMap := make(map[string]*Validator)

	m.Range(func(k, v interface{}) bool {
		tmpMap[k.(string)] = v.(*Validator)
		return true
	})
	return json.MarshalIndent(tmpMap, "  ", "  ")
}

func unmarshalJSON(data []byte) (*sync.Map, error) {
	var tmpMap map[string]*Validator
	m := &sync.Map{}

	if err := json.Unmarshal(data, &tmpMap); err != nil {
		return m, err
	}

	for key, value := range tmpMap {
		m.Store(key, value)
	}
	return m, nil
}
