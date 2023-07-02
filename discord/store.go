package discord

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"pactus-faucet/config"
	"sync"
)

// CacheEntry is a value stored in the cache.
type Validator struct {
	DiscordName      string  `json:"discord_name"`
	DiscordID        string  `json:"discord_id"`
	ValidatorAddress string  `json:"validator_address"`
	FaucetAmount     float64 `json:"faucet_amount"`
}

// SafeCache is a thread-safe cache.
type SafeStore struct {
	syncMap *sync.Map
	cfg     *config.Config
}

func LoadData(cfg *config.Config) (*SafeStore, error) {
	file, err := os.ReadFile(cfg.ValidatorDataPath)
	if err != nil {
		log.Printf("error loading validator data: %v", err)
		return nil, fmt.Errorf("error loading data file: %v", err)
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
		return nil, fmt.Errorf("error unmarshalling validator data: %v", err)
	}
	ss := &SafeStore{
		syncMap: data,
		cfg:     cfg,
	}
	return ss, nil
}

// Set a given value to the data storage
func (ss *SafeStore) SetData(address, discordName, discordID string, amount float64) error {
	ss.syncMap.Store(discordID, &Validator{DiscordName: discordName, DiscordID: discordID, ValidatorAddress: address, FaucetAmount: amount})
	//save record
	data, err := marshalJSON(ss.syncMap)
	if err != nil {
		log.Printf("error marshalling validator data file: %v", err)
		return fmt.Errorf("error marshalling validator data file: %v", err)
	}
	if err := os.WriteFile(ss.cfg.ValidatorDataPath, data, 0600); err != nil {
		log.Printf("failed to write to %s: %v", ss.cfg.ValidatorDataPath, err)
		return fmt.Errorf("failed to write to %s: %v", ss.cfg.ValidatorDataPath, err)
	}
	return nil
}

// Get retrives the given key from the storage
func (ss *SafeStore) GetData(discordID string) (*Validator, bool) {
	entry, found := ss.syncMap.Load(discordID)
	if !found {
		return nil, false
	}
	validator := entry.(*Validator)
	return validator, true
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
