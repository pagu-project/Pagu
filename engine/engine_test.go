package engine_test

import (
	"testing"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/wallet"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (engine.Engine, client.MockIClient, error) {
	t.Helper()
	ctrl := gomock.NewController(t)

	// mocking client manager
	sl := log.NewSubLogger("test")
	client1 := client.NewMockIClient(ctrl)

	cm := client.NewClientMgr()
	cm.AddClient("addr-1", client1)

	// mocking wallet
	wallet := wallet.NewMockIWallet(ctrl)

	// mocking wallet
	store := store.NewMockIStore(ctrl)

	eng, err := engine.NewBotEngine(sl, cm, wallet, store)
	return eng, *client1, err
}

func TestNetworkStatus(t *testing.T) {
	eng, client, err := setup(t)
	assert.NoError(t, err)

	client.EXPECT().GetNetworkInfo().Return(
		&pactus.GetNetworkInfoResponse{
			ConnectedPeersCount: 5,
			NetworkName:         "test",
		}, nil,
	)

	client.EXPECT().GetBlockchainInfo().Return(
		&pactus.GetBlockchainInfoResponse{
			TotalPower: 1234,
		}, nil,
	)

	status, err := eng.NetworkStatus([]string{})
	assert.NoError(t, err)

	assert.Equal(t, uint32(5), status.ConnectedPeersCount)
	assert.Equal(t, "test", status.NetworkName)
	assert.Equal(t, int64(1234), status.TotalNetworkPower)
}
