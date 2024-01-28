package engine

import (
	"errors"
	"testing"
	"time"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/log"
	rpstore "github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/wallet"
	"github.com/libp2p/go-libp2p/core/peer"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (*BotEngine, *client.MockIClient, *rpstore.MockIStore, *wallet.MockIWallet) {
	t.Helper()
	ctrl := gomock.NewController(t)

	// mocking client manager.
	sl := log.NewSubLogger("test")
	mockClient := client.NewMockIClient(ctrl)

	cm := client.NewClientMgr()
	cm.AddClient("addr-1", mockClient)

	// mocking mockWallet.
	mockWallet := wallet.NewMockIWallet(ctrl)

	// mocking mockStore.
	mockStore := rpstore.NewMockIStore(ctrl)

	eng := newBotEngine(sl, cm, mockWallet, mockStore)
	return eng, mockClient, mockStore, mockWallet
}

func TestNetworkStatus(t *testing.T) {
	eng, client, _, _ := setup(t)

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

	status, err := eng.NetworkStatus()
	assert.NoError(t, err)

	assert.Equal(t, uint32(5), status.ConnectedPeersCount)
	assert.Equal(t, "test", status.NetworkName)
	assert.Equal(t, int64(1234), status.TotalNetworkPower)
}

func TestNetworkHealth(t *testing.T) {
	eng, client, _, _ := setup(t)

	t.Run("should be healthy", func(t *testing.T) {
		currentTime := time.Now().Unix()
		client.EXPECT().LastBlockTime().Return(uint32(currentTime), uint32(100), nil)

		time.Sleep(2 * time.Second)

		healthy, err := eng.NetworkHealth()
		assert.NoError(t, err)

		assert.Equal(t, true, healthy.HealthStatus)
		assert.Equal(t, uint32(100), healthy.LastBlockHeight)
		assert.Equal(t, currentTime, healthy.LastBlockTime.Unix())
		assert.Equal(t, currentTime+2, healthy.CurrentTime.Unix())
		assert.Equal(t, int64(2), healthy.TimeDifference)
	})

	t.Run("should be unhealthy", func(t *testing.T) {
		currentTime := time.Now().Unix() - 16 // time difference is more than 15 seconds.
		client.EXPECT().LastBlockTime().Return(uint32(currentTime), uint32(100), nil)

		healthy, err := eng.NetworkHealth()
		assert.NoError(t, err)

		assert.Equal(t, false, healthy.HealthStatus)
	})
}

func TestNodeInfo(t *testing.T) {
	eng, client, _, _ := setup(t)

	t.Run("should return error, invalid input", func(t *testing.T) {
		client.EXPECT().GetNetworkInfo().Return(
			nil, errors.New(""),
		)

		info, err := eng.NodeInfo("pc1Invalid")

		assert.Nil(t, info)
		assert.Error(t, err)
	})

	t.Run("should work, valid address", func(t *testing.T) {
		valAddress := "valid-address"
		pubKey := "pub-key"

		peerID, err := peer.Decode("12D3KooWNwudyHVEwtyRTkTx9JoWgHo65hkPUxU12pKviAreVJYg")
		assert.NoError(t, err)

		client.EXPECT().GetNetworkInfo().Return(
			&pactus.GetNetworkInfoResponse{
				ConnectedPeers: []*pactus.PeerInfo{
					{
						ConsensusKeys:    []string{pubKey},
						ConsensusAddress: []string{valAddress},
						Height:           100,
						PeerId:           []byte(peerID),
						Agent:            "node=pactus-gui.exe/node-version=v0.20.0/protocol-version=1/os=windows/arch=amd64",
						Address:          "/ip4/000.000.000.000/tcp/21777",
					},
					{
						ConsensusKeys:    []string{pubKey},
						ConsensusAddress: []string{valAddress},
					},
				},
			}, nil,
		)

		client.EXPECT().GetValidatorInfo(valAddress).Return(
			&pactus.GetValidatorResponse{
				Validator: &pactus.ValidatorInfo{
					PublicKey:         pubKey,
					Stake:             int64(1_000),
					Address:           valAddress,
					Number:            1,
					AvailabilityScore: 0.9,
				},
			}, nil,
		).AnyTimes()

		info, err := eng.NodeInfo(valAddress)
		assert.NoError(t, err)

		assert.Equal(t, int64(1_000), info.StakeAmount)
		assert.Equal(t, float64(0.9), info.AvailabilityScore)
	})
}

func TestClaim(t *testing.T) {
	eng, client, store, wallet := setup(t)

	t.Run("everything normal and good", func(t *testing.T) {
		mainnetAddr := "mainnet-addr"
		testnetAddr := "testnet-addr"
		discordID := "123456789"
		amount := int64(30)
		memo := "TestNet reward claim from RoboPac"
		txID := "tx-id"

		wallet.EXPECT().Balance().Return(
			int64(501 * 1e9),
		).AnyTimes()

		store.EXPECT().ClaimerInfo(testnetAddr).Return(
			&rpstore.Claimer{
				DiscordID:   discordID,
				TotalReward: amount,
				ClaimedTxID: "",
			},
		)

		client.EXPECT().GetNetworkInfo().Return(
			&pactus.GetNetworkInfoResponse{
				ConnectedPeers: []*pactus.PeerInfo{
					{
						ConsensusAddress: []string{mainnetAddr},
						ConsensusKeys:    []string{"public-key-1"},
					},
				},
			}, nil,
		)

		wallet.EXPECT().BondTransaction("public-key-1", mainnetAddr, memo, amount).Return(
			txID, nil,
		)

		store.EXPECT().AddClaimTransaction(testnetAddr, txID).Return(
			nil,
		)

		expectedTx, err := eng.Claim(discordID, testnetAddr, mainnetAddr)
		assert.NoError(t, err)
		assert.NotNil(t, expectedTx, txID)

		//! can't claim twice:
		store.EXPECT().ClaimerInfo(testnetAddr).Return(
			&rpstore.Claimer{
				DiscordID:   discordID,
				TotalReward: amount,
				ClaimedTxID: txID,
			},
		)

		expectedTx, err = eng.Claim(discordID, testnetAddr, mainnetAddr)
		assert.Error(t, err)
		assert.Empty(t, expectedTx)
	})
}
