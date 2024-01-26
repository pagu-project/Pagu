package engine_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
	rpstore "github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/wallet"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (engine.Engine, client.MockIClient, rpstore.MockIStore, wallet.MockIWallet, error) {
	t.Helper()
	ctrl := gomock.NewController(t)

	// mocking client manager.
	sl := log.NewSubLogger("test")
	mockClient := client.NewMockIClient(ctrl)

	cm := client.NewClientMgr()
	cm.AddClient("addr-1", mockClient)

	// mocking wallet.
	wallet := wallet.NewMockIWallet(ctrl)

	// mocking store.
	store := rpstore.NewMockIStore(ctrl)

	eng, err := engine.NewBotEngine(sl, cm, wallet, store)
	return eng, *mockClient, *store, *wallet, err
}

func TestNetworkStatus(t *testing.T) {
	eng, client, _, _, err := setup(t)
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

func TestNetworkHealth(t *testing.T) {
	eng, client, _, _, err := setup(t)
	assert.NoError(t, err)

	t.Run("should be healthy", func(t *testing.T) {
		currentTime := time.Now().Unix()
		client.EXPECT().LastBlockTime().Return(uint32(currentTime), uint32(100), nil)

		time.Sleep(2 * time.Second)

		healthy, err := eng.NetworkHealth([]string{})
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

		healthy, err := eng.NetworkHealth([]string{})
		assert.NoError(t, err)

		assert.Equal(t, false, healthy.HealthStatus)
	})
}

func TestNodeInfo(t *testing.T) {
	eng, client, _, _, err := setup(t)
	assert.NoError(t, err)

	t.Run("should return error, invalid input", func(t *testing.T) {
		info, err := eng.NodeInfo([]string{})

		assert.Nil(t, info)
		assert.Error(t, err)
	})

	t.Run("should work, valid address", func(t *testing.T) {
		valAddress := "pc1p74scge5dyzjktv9q70xtr0pjmyqcqk7nuh8nzp"
		pubKey := "public1pk85lz4hkymm7ke3539p7fssqz8hqkvlwlcx0jvxn3s6has834l2skdr3649fznt6xdvkvz6rum0gxq4nunr9ta0vapz0wt92kdr2dj6qxt5qnm92j2mv8u8e8rj3nylyr7q9pn88myp49kht85eqxkqdsu5t39gx"

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
						ConsensusKeys: []string{"publicInvalid"},
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

		info, err := eng.NodeInfo([]string{valAddress})
		assert.NoError(t, err)

		assert.Equal(t, int64(1_000), info.StakeAmount)
		assert.Equal(t, float64(0.9), info.AvailabilityScore)
	})
}

func TestClaim(t *testing.T) {
	eng, client, store, wallet, err := setup(t)
	assert.NoError(t, err)

	t.Run("everything normal and good", func(t *testing.T) {
		valAddress := "pc1p74scge5dyzjktv9q70xtr0pjmyqcqk7nuh8nzp"
		discordID := "123456789"
		txID := "0x123456789"
		amount := 74.68
		time := time.Now().Unix()

		client.EXPECT().IsValidator(valAddress).Return(
			true, nil,
		)

		store.EXPECT().ClaimerInfo(discordID).Return(
			&rpstore.Claimer{
				DiscordID:        discordID,
				TotalReward:      amount,
				ClaimTransaction: nil,
			},
		)

		memo := fmt.Sprintf("RP to: %v", discordID)
		wallet.EXPECT().BondTransaction("", valAddress, memo, amount).Return(
			txID, nil,
		)

		client.EXPECT().GetTransactionData(txID).Return(
			&pactus.GetTransactionResponse{
				BlockTime: uint32(time),
				Transaction: &pactus.TransactionInfo{
					Id:    []byte(txID),
					Value: util.CoinToChange(amount),
					Memo:  memo,
				},
			}, nil,
		)

		store.EXPECT().AddClaimTransaction(txID, amount, time, discordID).Return(
			nil,
		)

		store.EXPECT().ClaimerInfo(discordID).Return(
			&rpstore.Claimer{
				DiscordID:   discordID,
				TotalReward: amount,
				ClaimTransaction: &rpstore.ClaimTransaction{
					TxID:   txID,
					Amount: amount,
					Time:   time,
				},
			},
		).AnyTimes()

		claimTx, err := eng.Claim([]string{valAddress, discordID})
		assert.NoError(t, err)
		assert.NotNil(t, claimTx)

		assert.Equal(t, amount, claimTx.Amount)
		assert.Equal(t, txID, claimTx.TxID)
		assert.Equal(t, time, claimTx.Time)

		//! can't claim twice.
		claimTx, err = eng.Claim([]string{valAddress, discordID})
		assert.EqualError(t, err, "this claimer have already claimed rewards")
		assert.Nil(t, claimTx)
	})

	t.Run("missing arguments", func(t *testing.T) {
		claimTx, err := eng.Claim([]string{})
		assert.EqualError(t, err, "missing argument: validator address")
		assert.Nil(t, claimTx)
	})

	t.Run("claimer not found", func(t *testing.T) {
		valAddress := "pc1p74scge5dyzjktv9q70xtr0pjmyqcqk7nuh8nzp"
		discordID := "987654321"

		store.EXPECT().ClaimerInfo(discordID).Return(
			nil,
		)

		claimTx, err := eng.Claim([]string{valAddress, discordID})
		assert.EqualError(t, err, "claimer not found")
		assert.Nil(t, claimTx)
	})

	t.Run("not validator address", func(t *testing.T) {
		valAddress := "pc1p74scge5dyzjktv9q70xtr0pjmyqcqk7nuh8nzp"
		discordID := "1234567890"
		amount := 74.68

		store.EXPECT().ClaimerInfo(discordID).Return(
			&rpstore.Claimer{
				DiscordID:   discordID,
				TotalReward: amount,
			},
		)

		client.EXPECT().IsValidator(valAddress).Return(
			false, nil,
		)

		claimTx, err := eng.Claim([]string{valAddress, discordID})
		assert.EqualError(t, err, "invalid argument: validator address")
		assert.Nil(t, claimTx)
	})

	t.Run("empty transaction ID", func(t *testing.T) {
		valAddress := "pc1p74scge5dyzjktv9q70xtr0pjmyqcqk7nuh8nzp"
		discordID := "1234567890"
		amount := 74.68

		client.EXPECT().IsValidator(valAddress).Return(
			true, nil,
		)

		store.EXPECT().ClaimerInfo(discordID).Return(
			&rpstore.Claimer{
				DiscordID:        discordID,
				TotalReward:      amount,
				ClaimTransaction: nil,
			},
		)

		memo := fmt.Sprintf("RP to: %v", discordID)
		wallet.EXPECT().BondTransaction("", valAddress, memo, amount).Return(
			"", nil,
		)

		claimTx, err := eng.Claim([]string{valAddress, discordID})
		assert.EqualError(t, err, "can't send bond transaction")
		assert.Nil(t, claimTx)
	})
}
