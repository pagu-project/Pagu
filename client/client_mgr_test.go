package client

import (
	"context"
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func setup(t *testing.T) (*Mgr, *MockIClient) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockClient := NewMockIClient(ctrl)

	clientMgr := NewClientMgr(context.Background())
	clientMgr.AddClient(mockClient)

	mockClient.EXPECT().GetNetworkInfo(clientMgr.ctx).Return(
		&pactus.GetNetworkInfoResponse{
			ConnectedPeers: []*pactus.PeerInfo{
				{
					ConsensusKeys:    []string{"pubKey-1", "pubKey-2"},
					ConsensusAddress: []string{"addr-1", "addr-2"},
				},
				{
					ConsensusKeys:    []string{"pubKey-3", "pubKey-4"},
					ConsensusAddress: []string{"addr-3", "addr-4"},
				},
			},
		}, nil,
	).AnyTimes()

	clientMgr.Start()

	return clientMgr, mockClient
}

func TestFindPublicKey(t *testing.T) {
	clientMgr, _ := setup(t)

	t.Run("not found", func(t *testing.T) {
		pubKey, err := clientMgr.FindPublicKey("not-exists", false)
		assert.Error(t, err)
		assert.Empty(t, pubKey)
	})

	t.Run("not first", func(t *testing.T) {
		pubKey, err := clientMgr.FindPublicKey("addr-4", true)
		assert.Error(t, err)
		assert.Empty(t, pubKey)
	})

	t.Run("first-ok", func(t *testing.T) {
		pubKey, err := clientMgr.FindPublicKey("addr-3", true)
		assert.NoError(t, err)
		assert.Equal(t, pubKey, "pubKey-3")
	})

	t.Run("any-ok", func(t *testing.T) {
		pubKey, err := clientMgr.FindPublicKey("addr-4", false)
		assert.NoError(t, err)
		assert.Equal(t, pubKey, "pubKey-4")
	})
}
