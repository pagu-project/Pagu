package client

import (
	"errors"

	"github.com/kehiy/RoboPac/log"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func init() {
	crypto.AddressHRP = "pc"
	crypto.PublicKeyHRP = "public"
	crypto.PrivateKeyHRP = "secret"
	crypto.XPublicKeyHRP = "xpublic"
	crypto.XPrivateKeyHRP = "xsecret"
}

type Mgr struct {
	clients map[string]IClient
}

func NewClientMgr() *Mgr {
	return &Mgr{
		clients: make(map[string]IClient),
	}
}

func (cm *Mgr) AddClient(addr string, c IClient) {
	cm.clients[addr] = c
}

func (cm *Mgr) GetRandomClient() IClient {
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

func (cm *Mgr) GetLastBlockTime() (uint32, uint32) {
	var lastBlockTime uint32 = 0
	var lastBlockHeight uint32 = 0
	for _, c := range cm.clients {
		t, h, err := c.LastBlockTime()
		if err != nil {
			continue
		}
		if t > lastBlockTime {
			lastBlockTime = t
			lastBlockHeight = h
		}
	}

	return lastBlockTime, lastBlockHeight
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

func (cm *Mgr) GetPeerInfoFirstVal(address string) (*pactus.PeerInfo, *bls.PublicKey, error) {
	for _, c := range cm.clients {
		networkInfo, err := c.GetNetworkInfo()
		if err != nil {
			continue
		}

		if networkInfo != nil {
			for _, p := range networkInfo.ConnectedPeers {
				for i, key := range p.ConsensusKeys {
					pub, _ := bls.PublicKeyFromString(key)
					if pub != nil {
						if pub.ValidatorAddress().String() == address {
							if i != 0 {
								return nil, nil, errors.New("please enter the first validator address")
							}
							return p, pub, nil
						}
					}
				}
			}
		}
	}

	return nil, nil, errors.New("peer does not exist")
}

func (cm *Mgr) GetPeerInfo(address string) (*pactus.PeerInfo, *bls.PublicKey, error) {
	for _, c := range cm.clients {
		networkInfo, err := c.GetNetworkInfo()
		if err != nil {
			continue
		}

		if networkInfo != nil {
			for _, p := range networkInfo.ConnectedPeers {
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

func (cm *Mgr) GetValidatorInfo(address string) (*pactus.GetValidatorResponse, error) {
	for _, c := range cm.clients {
		val, err := c.GetValidatorInfo(address)
		if err != nil {
			continue
		}
		return val, nil
	}

	return nil, errors.New("unable to get validator info")
}

func (cm *Mgr) GetValidatorInfoByNumber(num int32) (*pactus.GetValidatorResponse, error) {
	for _, c := range cm.clients {
		val, err := c.GetValidatorInfoByNumber(num)
		if err != nil {
			continue
		}
		return val, nil
	}

	return nil, errors.New("unable to get validator info")
}

func (cm *Mgr) Close() {
	for addr, c := range cm.clients {
		if err := c.Close(); err != nil {
			log.Error("could not close connection to RPC node", "err", err, "RPCAddr", addr)
		}
	}
}
