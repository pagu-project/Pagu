package client

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kehiy/RoboPac/log"
	"github.com/pactus-project/pactus/util/logger"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type Mgr struct {
	valMapLock sync.RWMutex
	valMap     map[string]*pactus.PeerInfo

	ctx     context.Context
	clients []IClient
}

func NewClientMgr(ctx context.Context) *Mgr {
	return &Mgr{
		clients:    make([]IClient, 0),
		valMap:     make(map[string]*pactus.PeerInfo),
		valMapLock: sync.RWMutex{},
		ctx:        ctx,
	}
}

func (cm *Mgr) Start() {
	ticker := time.NewTicker(30 * time.Minute)

	go func() {
		for {
			select {
			case <-cm.ctx.Done():
				return

			case <-ticker.C:
				logger.Info("updating validator map started")
				cm.updateValMap()
			}
		}
	}()

	cm.updateValMap()
}

func (cm *Mgr) Stop() {
	for addr, c := range cm.clients {
		if err := c.Close(); err != nil {
			log.Error("could not close connection to RPC node", "err", err, "RPCAddr", addr)
		}
	}
}

func (cm *Mgr) updateValMap() {
	freshValMap := make(map[string]*pactus.PeerInfo)

	for _, c := range cm.clients {
		networkInfo, err := c.GetNetworkInfo(cm.ctx)
		if err != nil {
			continue
		}

		if networkInfo == nil {
			logger.Warn("network info is nil")
			continue
		}

		for _, p := range networkInfo.ConnectedPeers {
			for _, addr := range p.ConsensusAddress {
				current := freshValMap[addr]
				if current != nil {
					if current.LastSent < p.LastSent {
						freshValMap[addr] = p
					}
				} else {
					freshValMap[addr] = p
				}
			}
		}
	}

	cm.valMapLock.Lock()
	clear(cm.valMap)
	cm.valMap = freshValMap
	cm.valMapLock.Unlock()

	logger.Info("validator map updated successfully")
}

// AddClient should call before Start.
func (cm *Mgr) AddClient(c IClient) {
	cm.clients = append(cm.clients, c)
}

// NOTE: local client is always the first client.
func (cm *Mgr) getLocalClient() IClient {
	return cm.clients[0]
}

func (cm *Mgr) GetRandomClient() IClient {
	for _, c := range cm.clients {
		return c
	}

	return nil
}

func (cm *Mgr) GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	localClient := cm.getLocalClient()
	info, err := localClient.GetBlockchainInfo(cm.ctx)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (cm *Mgr) GetBlockchainHeight() (uint32, error) {
	localClient := cm.getLocalClient()
	height, err := localClient.GetBlockchainHeight(cm.ctx)
	if err != nil {
		return 0, err
	}
	return height, nil
}

func (cm *Mgr) GetLastBlockTime() (uint32, uint32) {
	localClient := cm.getLocalClient()
	lastBlockTime, lastBlockHeight, err := localClient.LastBlockTime(cm.ctx)
	if err != nil {
		return 0, 0
	}

	return lastBlockTime, lastBlockHeight
}

func (cm *Mgr) GetNetworkInfo() (*pactus.GetNetworkInfoResponse, error) {
	for _, c := range cm.clients {
		info, err := c.GetNetworkInfo(cm.ctx)
		if err != nil {
			continue
		}
		return info, nil
	}

	return nil, errors.New("unable to get network info")
}

func (cm *Mgr) FindPublicKey(address string, firstVal bool) (string, error) {
	peerInfo, err := cm.GetPeerInfo(address)
	if err != nil {
		return "", err
	}

	for i, addr := range peerInfo.ConsensusAddress {
		if addr == address {
			if firstVal && i != 0 {
				return "", errors.New("please enter the first validator address")
			}
			return peerInfo.ConsensusKeys[i], nil
		}
	}

	panic("unreachable")
}

func (cm *Mgr) GetPeerInfo(address string) (*pactus.PeerInfo, error) {
	cm.valMapLock.Lock()
	defer cm.valMapLock.Unlock()

	peerInfo, ok := cm.valMap[address]
	if !ok {
		return nil, fmt.Errorf("peer does not exist with this address: %v", address)
	}

	return peerInfo, nil
}

func (cm *Mgr) GetValidatorInfo(address string) (*pactus.GetValidatorResponse, error) {
	localClient := cm.getLocalClient()
	val, err := localClient.GetValidatorInfo(cm.ctx, address)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (cm *Mgr) GetValidatorInfoByNumber(num int32) (*pactus.GetValidatorResponse, error) {
	localClient := cm.getLocalClient()
	val, err := localClient.GetValidatorInfoByNumber(cm.ctx, num)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (cm *Mgr) GetTransactionData(txID string) (*pactus.GetTransactionResponse, error) {
	localClient := cm.getLocalClient()
	txData, err := localClient.GetTransactionData(cm.ctx, txID)
	if err != nil {
		return nil, err
	}
	return txData, nil
}

func (cm *Mgr) GetBalance(addr string) (int64, error) {
	return cm.getLocalClient().GetBalance(cm.ctx, addr)
}

func (cm *Mgr) GetCirculatingSupply() (int64, error) {
	localClient := cm.getLocalClient()

	height, err := localClient.GetBlockchainInfo(cm.ctx)
	if err != nil {
		return 0, err
	}
	minted := float64(height.LastBlockHeight) * 1e9
	staked := height.TotalPower
	warm := int64(630_000_000_000_000)

	var addr1Out int64 = 0
	var addr2Out int64 = 0
	var addr3Out int64 = 0
	var addr4Out int64 = 0
	var addr5Out int64 = 0 // warm wallet
	var addr6Out int64 = 0 // warm wallet

	balance1, err := localClient.GetBalance(cm.ctx, "pc1z2r0fmu8sg2ffa0tgrr08gnefcxl2kq7wvquf8z")
	if err == nil {
		addr1Out = 8_400_000_000_000_000 - balance1
	}

	balance2, err := localClient.GetBalance(cm.ctx, "pc1zprhnvcsy3pthekdcu28cw8muw4f432hkwgfasv")
	if err == nil {
		addr2Out = 6_300_000_000_000_000 - balance2
	}

	balance3, err := localClient.GetBalance(cm.ctx, "pc1znn2qxsugfrt7j4608zvtnxf8dnz8skrxguyf45")
	if err == nil {
		addr3Out = 4_200_000_000_000_000 - balance3
	}

	balance4, err := localClient.GetBalance(cm.ctx, "pc1zs64vdggjcshumjwzaskhfn0j9gfpkvche3kxd3")
	if err == nil {
		addr4Out = 2_100_000_000_000_000 - balance4
	}

	balance5, err := localClient.GetBalance(cm.ctx, "pc1zuavu4sjcxcx9zsl8rlwwx0amnl94sp0el3u37g")
	if err == nil {
		addr5Out = 420_000_000_000_000 - balance5
	}

	balance6, err := localClient.GetBalance(cm.ctx, "pc1zf0gyc4kxlfsvu64pheqzmk8r9eyzxqvxlk6s6t")
	if err == nil {
		addr6Out = 210_000_000_000_000 - balance6
	}

	circulating := (addr1Out + addr2Out + addr3Out + addr4Out + addr5Out + addr6Out + int64(minted)) - staked - warm
	return circulating, nil
}
