package client

import (
	"context"
	"errors"
	"log"

	"github.com/k0kubun/pp"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
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
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	pp.Println("connection established...")

	return &Client{
		blockchainClient:  pactus.NewBlockchainClient(conn),
		networkClient:     pactus.NewNetworkClient(conn),
		transactionClient: pactus.NewTransactionClient(conn),
		conn:              conn,
	}, nil

}

func (c *Client) GetBlockchainInfo() (*pactus.GetBlockchainInfoResponse, error) {
	blockchainInfo, err := c.blockchainClient.GetBlockchainInfo(context.Background(), &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		log.Printf("error obtaining block height: %v", err)
		return nil, err
	}
	return blockchainInfo, nil
}
func (c *Client) GetBlockchainHeight() (uint32, error) {
	blockchainInfo, err := c.blockchainClient.GetBlockchainInfo(context.Background(), &pactus.GetBlockchainInfoRequest{})
	if err != nil {
		log.Printf("error obtaining block height: %v", err)
		return 0, err
	}
	return blockchainInfo.LastBlockHeight, nil
}

func (c *Client) GetNetworkInfo() (*pactus.GetNetworkInfoResponse, error) {
	networkInfo, err := c.networkClient.GetNetworkInfo(context.Background(), &pactus.GetNetworkInfoRequest{})
	if err != nil {
		log.Printf("error obtaining network information: %v", err)

		return nil, err
	}

	return networkInfo, nil
}

func (c *Client) GetPeerInfo(address string) (*pactus.PeerInfo, *bls.PublicKey, error) {
	networkInfo, _ := c.GetNetworkInfo()
	crypto.PublicKeyHRP = "tpublic"
	if networkInfo != nil {
		for _, p := range networkInfo.Peers {
			for _, key := range p.ConsensusKeys {
				pub, _ := bls.PublicKeyFromString(key)
				if pub != nil {
					if pub.Address().String() == address {
						return p, pub, nil
					}
				}
			}
		}
	}
	return nil, nil, errors.New("peer does not exist")
}

func (c *Client) Close() error {
	return c.conn.Close()
}
