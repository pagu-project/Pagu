package client

import (
	"context"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type IClient interface {
	GetBlockchainInfo(context.Context) (*pactus.GetBlockchainInfoResponse, error)
	GetBlockchainHeight(context.Context) (uint32, error)
	GetLastBlockTime(context.Context) (uint32, uint32)
	GetNetworkInfo(context.Context) (*pactus.GetNetworkInfoResponse, error)
	GetValidatorInfo(context.Context, string) (*pactus.GetValidatorResponse, error)
	GetValidatorInfoByNumber(context.Context, int32) (*pactus.GetValidatorResponse, error)
	GetTransactionData(context.Context, string) (*pactus.GetTransactionResponse, error)
	GetBalance(context.Context, string) (int64, error)
	GetFee(context.Context, int64) (int64, error)
	Close() error
}

type Manager interface {
	Start()
	Stop()
	updateValMap()
	AddClient(c IClient)
	GetLocalClient() IClient
	GetRandomClient() IClient
	GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error)
	GetBlockchainHeight() (uint32, error)
	GetLastBlockTime() (uint32, uint32)
	GetNetworkInfo() (*pactus.GetNetworkInfoResponse, error)
	GetPeerInfo(address string) (*pactus.PeerInfo, error)
	GetValidatorInfo(address string) (*pactus.GetValidatorResponse, error)
	GetValidatorInfoByNumber(num int32) (*pactus.GetValidatorResponse, error)
	GetTransactionData(txID string) (*pactus.GetTransactionResponse, error)
	GetBalance(addr string) (int64, error)
	GetFee(amt int64) (int64, error)
	GetCirculatingSupply() (int64, error)
	FindPublicKey(address string, firstVal bool) (string, error)
}
