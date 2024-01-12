package store

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kehiy/RoboPac/config"
)

// TO:DO:::

// Store is a thread-safe cache.
type Store struct {
	syncMap *sync.Map
	cfg     *config.Config
}

func LoadStore(cfg *config.Config) (IStore, error) {
	file, err := os.ReadFile(cfg.StorePath)
	if err != nil {
		log.Printf("error loading validator data: %v", err)
		return nil, fmt.Errorf("error loading data file: %w", err)
	}
	if len(file) == 0 {
		ss := &Store{
			syncMap: &sync.Map{},
			cfg:     cfg,
		}
		return ss, nil
	}

	// data, err := unmarshalJSON(file)
	// if err != nil {
	// 	log.Printf("error unmarshalling validator data: %v", err)
	// 	return nil, fmt.Errorf("error unmarshalling validator data: %w", err)
	// }

	ss := &Store{
		// syncMap: data,
		cfg: cfg,
	}
	return ss, nil
}

func (s *Store) Set() bool {
	return true
}

func (s *Store) Get() {}

// func marshalJSON(m *sync.Map) ([]byte, error) {
// 	tmpMap := make(map[string]*Validator)

// 	m.Range(func(k, v interface{}) bool {
// 		tmpMap[k.(string)] = v.(*Validator)
// 		return true
// 	})
// 	return json.MarshalIndent(tmpMap, "  ", "  ")
// }

// func unmarshalJSON(data []byte) (*sync.Map, error) {
// 	var tmpMap map[string]*Validator
// 	m := &sync.Map{}

// 	if err := json.Unmarshal(data, &tmpMap); err != nil {
// 		return m, err
// 	}

// 	for key, value := range tmpMap {
// 		m.Store(key, value)
// 	}
// 	return m, nil
// }
