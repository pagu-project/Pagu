package engine

import (
	"errors"
	"testing"
	"time"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/log"
	rpstore "github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/utils"
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
		pubKey := "public-key"
		discordID := "123456789"
		amount := int64(30)
		memo := "TestNet reward claim from RoboPac"
		txID := "tx-id"

		wallet.EXPECT().Balance().Return(
			utils.CoinToAtomic(501),
		).MaxTimes(2)

		client.EXPECT().GetValidatorInfo(mainnetAddr).Return(
			&pactus.GetValidatorResponse{
				Validator: &pactus.ValidatorInfo{
					Stake: 0,
				},
			}, nil,
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
						ConsensusKeys:    []string{pubKey},
					},
				},
			}, nil,
		)

		wallet.EXPECT().BondTransaction(pubKey, mainnetAddr, memo, amount).Return(
			txID, nil,
		).MaxTimes(1)

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

	t.Run("should fail, low balance", func(t *testing.T) {
		mainnetAddr := "mainnet-addr-fail-balance"
		testnetAddr := "testnet-addr-fail-balance"
		discordID := "123456789-fail-balance"

		wallet.EXPECT().Balance().Return(
			utils.CoinToAtomic(499),
		)

		client.EXPECT().GetValidatorInfo(mainnetAddr).Return(
			&pactus.GetValidatorResponse{
				Validator: &pactus.ValidatorInfo{
					Stake: 0,
				},
			}, nil,
		)
		expectedTx, err := eng.Claim(discordID, testnetAddr, mainnetAddr)
		assert.EqualError(t, err, "insufficient wallet balance")
		assert.Empty(t, expectedTx)
	})

	t.Run("should fail, claimer not found", func(t *testing.T) {
		mainnetAddr := "mainnet-addr-fail-notfound"
		testnetAddr := "testnet-addr-fail-notfound"
		discordID := "123456789-fail-notfound"

		wallet.EXPECT().Balance().Return(
			utils.CoinToAtomic(501),
		)

		client.EXPECT().GetValidatorInfo(mainnetAddr).Return(
			&pactus.GetValidatorResponse{
				Validator: &pactus.ValidatorInfo{
					Stake: 0,
				},
			}, nil,
		)

		store.EXPECT().ClaimerInfo(testnetAddr).Return(
			nil,
		)

		expectedTx, err := eng.Claim(discordID, testnetAddr, mainnetAddr)
		assert.EqualError(t, err, "claimer not found")
		assert.Empty(t, expectedTx)
	})

	t.Run("should fail, different Discord ID", func(t *testing.T) {
		mainnetAddr := "mainnet-addr-fail-different-id"
		testnetAddr := "testnet-addr-fail-different-id"
		discordID := "123456789-fail-different-id"

		wallet.EXPECT().Balance().Return(
			utils.CoinToAtomic(501),
		)

		client.EXPECT().GetValidatorInfo(mainnetAddr).Return(
			&pactus.GetValidatorResponse{
				Validator: &pactus.ValidatorInfo{
					Stake: 0,
				},
			}, nil,
		)

		store.EXPECT().ClaimerInfo(testnetAddr).Return(
			&rpstore.Claimer{
				DiscordID: "invalid-discord-id",
			},
		)

		expectedTx, err := eng.Claim(discordID, testnetAddr, mainnetAddr)
		assert.EqualError(t, err, "invalid claimer")
		assert.Empty(t, expectedTx)
	})

	t.Run("should fail, not first validator address", func(t *testing.T) {
		mainnetAddr := "mainnet-addr-fail-not-first-validator"
		testnetAddr := "testnet-addr-fail-not-first-validator"
		discordID := "123456789-fail-not-first-validator"

		wallet.EXPECT().Balance().Return(
			utils.CoinToAtomic(501),
		)

		client.EXPECT().GetValidatorInfo(mainnetAddr).Return(
			&pactus.GetValidatorResponse{
				Validator: &pactus.ValidatorInfo{
					Stake: 0,
				},
			}, nil,
		)

		store.EXPECT().ClaimerInfo(testnetAddr).Return(
			&rpstore.Claimer{
				DiscordID:   discordID,
				ClaimedTxID: "",
			},
		)

		client.EXPECT().GetNetworkInfo().Return(
			&pactus.GetNetworkInfoResponse{
				ConnectedPeers: []*pactus.PeerInfo{
					{
						ConsensusAddress: []string{"invalid-address", mainnetAddr},
					},
				},
			}, nil,
		)

		expectedTx, err := eng.Claim(discordID, testnetAddr, mainnetAddr)
		assert.EqualError(t, err, "please enter the first validator address")
		assert.Empty(t, expectedTx)
	})

	t.Run("should fail, validator not found", func(t *testing.T) {
		mainnetAddr := "mainnet-addr-fail-validator-not-found"
		testnetAddr := "testnet-addr-fail-validator-not-found"
		discordID := "123456789-fail-validator-not-found"

		wallet.EXPECT().Balance().Return(
			utils.CoinToAtomic(501),
		)

		client.EXPECT().GetValidatorInfo(mainnetAddr).Return(
			&pactus.GetValidatorResponse{
				Validator: &pactus.ValidatorInfo{
					Stake: 0,
				},
			}, nil,
		)

		store.EXPECT().ClaimerInfo(testnetAddr).Return(
			&rpstore.Claimer{
				DiscordID:   discordID,
				ClaimedTxID: "",
			},
		)

		client.EXPECT().GetNetworkInfo().Return(
			&pactus.GetNetworkInfoResponse{
				ConnectedPeers: []*pactus.PeerInfo{
					{
						ConsensusAddress: []string{"invalid-address", "invalid-address-2"},
					},
				},
			}, nil,
		)

		expectedTx, err := eng.Claim(discordID, testnetAddr, mainnetAddr)
		assert.EqualError(t, err, "peer does not exist")
		assert.Empty(t, expectedTx)
	})

	t.Run("should fail, empty transaction hash", func(t *testing.T) {
		mainnetAddr := "mainnet-addr-fail-empty-tx-hash"
		testnetAddr := "testnet-addr-fail-empty-tx-hash"
		pubKey := "public-key-fail-empty-tx-hash"
		discordID := "123456789-fail-empty-tx-hash"
		amount := int64(30)
		memo := "TestNet reward claim from RoboPac"

		wallet.EXPECT().Balance().Return(
			utils.CoinToAtomic(501),
		)

		client.EXPECT().GetValidatorInfo(mainnetAddr).Return(
			&pactus.GetValidatorResponse{
				Validator: &pactus.ValidatorInfo{
					Stake: 0,
				},
			}, nil,
		)

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
						ConsensusKeys:    []string{pubKey},
					},
				},
			}, nil,
		)

		wallet.EXPECT().BondTransaction(pubKey, mainnetAddr, memo, amount).Return(
			"", nil,
		)

		expectedTx, err := eng.Claim(discordID, testnetAddr, mainnetAddr)
		assert.EqualError(t, err, "can't send bond transaction")
		assert.Empty(t, expectedTx)
	})

	t.Run("should panic, add claimer failed", func(t *testing.T) {
		mainnetAddr := "mainnet-addr-panic-add-claimer-failed"
		testnetAddr := "testnet-addr-panic-add-claimer-failed"
		pubKey := "public-key-panic-add-claimer-failed"
		discordID := "123456789-panic-add-claimer-failed"
		amount := int64(30)
		memo := "TestNet reward claim from RoboPac"
		txID := "tx-id-panic-add-claimer-failed"

		wallet.EXPECT().Balance().Return(
			utils.CoinToAtomic(501),
		)

		client.EXPECT().GetValidatorInfo(mainnetAddr).Return(
			&pactus.GetValidatorResponse{
				Validator: &pactus.ValidatorInfo{
					Stake: 0,
				},
			}, nil,
		)
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
						ConsensusKeys:    []string{pubKey},
					},
				},
			}, nil,
		)

		wallet.EXPECT().BondTransaction(pubKey, mainnetAddr, memo, amount).Return(
			txID, nil,
		)

		store.EXPECT().AddClaimTransaction(testnetAddr, txID).Return(
			errors.New(""),
		)

		assert.Panics(t, func() {
			_, _ = eng.Claim(discordID, testnetAddr, mainnetAddr)
		})
	})
}
