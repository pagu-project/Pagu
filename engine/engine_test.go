package engine_test

import (
	"testing"
	"time"

	"github.com/kehiy/RoboPac/client"
	"github.com/kehiy/RoboPac/engine"
	"github.com/kehiy/RoboPac/log"
	"github.com/kehiy/RoboPac/store"
	"github.com/kehiy/RoboPac/wallet"
	"github.com/libp2p/go-libp2p/core/peer"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setup(t *testing.T) (engine.Engine, client.MockIClient, error) {
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
	store := store.NewMockIStore(ctrl)

	eng, err := engine.NewBotEngine(sl, cm, wallet, store)
	return eng, *mockClient, err
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

func TestNetworkHealth(t *testing.T) {
	eng, client, err := setup(t)
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
	eng, client, err := setup(t)
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
						ConsensusKeys: []string{pubKey},
						Height:        100,
						PeerId:        []byte(peerID),
						Agent:         "node=pactus-gui.exe/node-version=v0.20.0/protocol-version=1/os=windows/arch=amd64",
						Address:       "/ip4/000.000.000.000/tcp/21777",
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
