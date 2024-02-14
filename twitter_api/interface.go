package twitter_api

import (
	"context"
	"time"
)

type UserInfo struct {
	TwitterID   string
	TwitterName string
	CreatedAt   time.Time
	Followers   int
	IsVerified  bool
}

type TweetInfo struct {
	ID        string
	Link      string
	CreatedAt time.Time
}

type IClient interface {
	UserInfo(ctx context.Context, twitterName string) (*UserInfo, error)
	RetweetSearch(ctx context.Context, discordID, twitterName string) (*TweetInfo, error)
}
