package client

import (
	"errors"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func init() {
	crypto.AddressHRP = "tpc"
	crypto.PublicKeyHRP = "tpublic"
}

type Mgr struct {
	clients map[string]*Client
}

func NewClientMgr() *Mgr {
	return &Mgr{
		clients: make(map[string]*Client),
	}
}

func (cm *Mgr) AddClient(addr string, c *Client) {
	cm.clients[addr] = c
}

func (cm *Mgr) GetRandomClient() *Client {
	for _, c := range cm.clients {
		return c
	}

	return nil
}

func (cm *Mgr) GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	for _, c := range cm.clients {
		info, err := c.GetBlockchainInfo()
		if err != nil {
			continue
		}
		return info, nil
	}

	return nil, errors.New("unable to get blockchain info")
}

func (cm *Mgr) GetBlockchainHeight() (uint32, error) {
	for _, c := range cm.clients {
		height, err := c.GetBlockchainHeight()
		if err != nil {
			continue
		}
		return height, nil
	}

	return 0, errors.New("unable to get blockchain height")
}

func (cm *Mgr) GetNetworkInfo() (*pactus.GetNetworkInfoResponse, error) {
	for _, c := range cm.clients {
		info, err := c.GetNetworkInfo()
		if err != nil {
			continue
		}
		return info, nil
	}

	return nil, errors.New("unable to get network info")
}

func (cm *Mgr) GetPeerInfo(address string) (*pactus.PeerInfo, *bls.PublicKey, error) {
	for _, c := range cm.clients {
		networkInfo, err := c.GetNetworkInfo()
		if err != nil {
			continue
		}

		if networkInfo != nil {
			for _, p := range networkInfo.Peers {
				for _, key := range p.ConsensusKeys {
					pub, _ := bls.PublicKeyFromString(key)
					if pub != nil {
						if pub.ValidatorAddress().String() == address {
							return p, pub, nil
						}
					}
				}
			}
		}
	}

	return nil, nil, errors.New("peer does not exist")
}

func (cm *Mgr) IsValidator(address string) (bool, error) {
	for _, c := range cm.clients {
		exists, err := c.IsValidator(address)
		if err != nil {
			continue
		}
		return exists, nil
	}

	return false, errors.New("unable to get validator info")
}

func (cm *Mgr) Close() {
	for addr, c := range cm.clients {
		if err := c.Close(); err != nil {
			fmt.Printf("error on closing client %s\n", addr)
		}
	}
}
