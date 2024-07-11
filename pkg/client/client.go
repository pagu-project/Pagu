package client

import (
	"context"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/pagu-project/Pagu/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	blockchainClient  pactus.BlockchainClient
	networkClient     pactus.NetworkClient
	transactionClient pactus.TransactionClient
	conn              *grpc.ClientConn
}

func NewClient(endpoint string) (*Client, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	log.Info("establishing new connection", "addr", endpoint)

	return &Client{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		networkClient:     pactus.NewNetworkClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
		conn:              conn,
	}, nil
}

func (c *Client) GetBlockchainInfo(ctx context.Context) (*pactus.GetBlockchainInfoResponse, error) {
	blockchainInfo, err := c.blockchainClient.GetBlockchainInfo(ctx, &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return nil, err
	}
	return blockchainInfo, nil
}

func (c *Client) GetBlockchainHeight(ctx context.Context) (uint32, error) {
	blockchainInfo, err := c.blockchainClient.GetBlockchainInfo(ctx, &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return 0, err
	}
	return blockchainInfo.LastBlockHeight, nil
}

func (c *Client) GetNetworkInfo(ctx context.Context) (*pactus.GetNetworkInfoResponse, error) {
	networkInfo, err := c.networkClient.GetNetworkInfo(ctx, &pactus.GetNetworkInfoRequest{})
	if err != nil {
		return nil, err
	}

	return networkInfo, nil
}

func (c *Client) GetPeerInfo(ctx context.Context, address string) (*pactus.PeerInfo, error) {
	networkInfo, _ := c.GetNetworkInfo(ctx)
	if networkInfo != nil {
		for _, p := range networkInfo.ConnectedPeers {
			for _, addr := range p.ConsensusAddresses {
				if addr != "" {
					if addr == address {
						return p, nil
					}
				}
			}
		}
	}
	return nil, NotFoundError{
		Search:  "peer",
		Address: address,
	}
}

func (c *Client) GetValidatorInfo(ctx context.Context, address string) (*pactus.GetValidatorResponse, error) {
	validator, err := c.blockchainClient.GetValidator(ctx,
		&pactus.GetValidatorRequest{Address: address})
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func (c *Client) GetValidatorInfoByNumber(ctx context.Context, num int32) (*pactus.GetValidatorResponse, error) {
	validator, err := c.blockchainClient.GetValidatorByNumber(ctx,
		&pactus.GetValidatorByNumberRequest{Number: num})
	if err != nil {
		return nil, err
	}

	return validator, nil
}

func (c *Client) GetTransactionData(ctx context.Context, hash string) (*pactus.GetTransactionResponse, error) {
	res, err := c.transactionClient.GetTransaction(ctx,
		&pactus.GetTransactionRequest{
			Id:        hash,
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_DATA,
		})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) GetLastBlockTime(ctx context.Context) (lastBlockTime, lastBlockHeight uint32) {
	info, err := c.blockchainClient.GetBlockchainInfo(ctx, &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		return 0, 0
	}

	lastBlock, err := c.blockchainClient.GetBlock(ctx, &pactus.GetBlockRequest{
		Height:    info.LastBlockHeight,
		Verbosity: pactus.BlockVerbosity_BLOCK_INFO,
	})
	if err != nil {
		return 0, 0
	}

	lastBlockHeight = lastBlock.Height
	lastBlockTime = lastBlock.BlockTime

	return lastBlockTime, lastBlockHeight
}

func (c *Client) GetNodeInfo(ctx context.Context) (*pactus.GetNodeInfoResponse, error) {
	info, err := c.networkClient.GetNodeInfo(ctx, &pactus.GetNodeInfoRequest{})
	if err != nil {
		return &pactus.GetNodeInfoResponse{}, err
	}

	return info, err
}

func (c *Client) GetBalance(ctx context.Context, address string) (int64, error) {
	account, err := c.blockchainClient.GetAccount(ctx, &pactus.GetAccountRequest{
		Address: address,
	})
	if err != nil {
		return 0, err
	}

	return account.Account.Balance, nil
}

func (c *Client) GetFee(ctx context.Context, amt int64) (int64, error) {
	res, err := c.transactionClient.CalculateFee(ctx, &pactus.CalculateFeeRequest{
		Amount:      amt,
		PayloadType: pactus.PayloadType_TRANSFER_PAYLOAD,
	})
	if err != nil {
		return 0, err
	}

	return res.Fee, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
