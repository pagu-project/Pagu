package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/kehiy/RoboPac/config"
	"github.com/kehiy/RoboPac/log"
)

// Store is a thread-safe cache.
type Store struct {
	syncMap *sync.Map
	cfg     *config.Config
	logger  *log.SubLogger
}

func LoadStore(cfg *config.Config, logger *log.SubLogger) (IStore, error) {
	file, err := os.ReadFile(cfg.StorePath)
	if err != nil {
		return nil, fmt.Errorf("error loading data file: %w", err)
	}
	if len(file) == 0 {
		ss := &Store{
			syncMap: &sync.Map{},
			cfg:     cfg,
		}
		return ss, nil
	}

	data, err := unmarshalJSON(file)
	if err != nil {
		return nil, fmt.Errorf("error un-marshalling validator data: %w", err)
	}

	ss := &Store{
		syncMap: data,
		cfg:     cfg,
		logger:  logger,
	}
	return ss, nil
}

func (s *Store) ClaimerInfo(testNetValAddr string) *Claimer {
	entry, found := s.syncMap.Load(testNetValAddr)
	if !found {
		return nil
	}

	claimerInfo := entry.(*Claimer)
	return claimerInfo
}

func (s *Store) AddClaimTransaction(amount float64, time int64, txID, discordID, testNetValAddr string) error {
	s.syncMap.Store(testNetValAddr, &Claimer{
		DiscordID: discordID,
		ClaimTransaction: &ClaimTransaction{
			TxID:   txID,
			Amount: amount,
			Time:   time,
		},
	})

	s.logger.Info("new claim transaction added", "discordID", discordID, "amount",
		amount, "time", time, "txID", txID)

	// save record.
	data, err := marshalJSON(s.syncMap)
	if err != nil {
		s.logger.Panic("can't marshal json new claim transaction", "discordID", discordID, "amount",
			amount, "time", time, "txID", txID, "err", err)

		return fmt.Errorf("error marshalling validator data file: %w", err)
	}

	if err := os.WriteFile(s.cfg.StorePath, data, 0o600); err != nil {
		s.logger.Panic("can't write new claim transaction", "discordID", discordID, "amount",
			amount, "time", time, "txID", txID, "err", err)

		return fmt.Errorf("failed to write to %s: %w", s.cfg.StorePath, err)
	}

	return nil
}

func marshalJSON(m *sync.Map) ([]byte, error) {
	tmpMap := make(map[string]*Claimer)

	m.Range(func(k, v interface{}) bool {
		tmpMap[k.(string)] = v.(*Claimer)
		return true
	})
	return json.MarshalIndent(tmpMap, "  ", "  ")
}

func unmarshalJSON(data []byte) (*sync.Map, error) {
	var tmpMap map[string]*Claimer
	m := &sync.Map{}

	if err := json.Unmarshal(data, &tmpMap); err != nil {
		return m, err
	}

	for key, value := range tmpMap {
		m.Store(key, value)
	}
	return m, nil
}
