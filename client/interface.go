package client

import (
	"context"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

type IClient interface {
	GetBlockchainInfo(context.Context) (*pactus.GetBlockchainInfoResponse, error)
	GetBlockchainHeight(context.Context) (uint32, error)
	LastBlockTime(context.Context) (uint32, uint32, error)
	GetNetworkInfo(context.Context) (*pactus.GetNetworkInfoResponse, error)
	GetValidatorInfo(context.Context, string) (*pactus.GetValidatorResponse, error)
	GetValidatorInfoByNumber(context.Context, int32) (*pactus.GetValidatorResponse, error)
	GetTransactionData(context.Context, string) (*pactus.GetTransactionResponse, error)
	GetBalance(context.Context, string) (int64, error)
	Close() error
}
