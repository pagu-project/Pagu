package twitter_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/pactus-project/pactus/util/logger"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

type Client struct {
	client    *twitter.Client
	twitterID string
}

func NewClient(bearerToken string, twitterID string) (*Client, error) {
	client := &twitter.Client{
		Authorizer: authorize{
			Token: bearerToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	// test the connection and everything is OK
	opts := twitter.UserLookupOpts{}
	username, err := client.UserLookup(context.Background(), []string{twitterID}, opts)
	if err != nil {
		return nil, err
	}

	if len(username.Raw.Users) == 0 || username.Raw.Users[0] == nil {
		return nil, fmt.Errorf("unable to find %v", twitterID)
	}

	logger.Info("found twitter", "name", username.Raw.Users[0].Name, "id", twitterID)

	return &Client{
		client:    client,
		twitterID: twitterID,
	}, nil
}

func (c *Client) UserInfo(ctx context.Context, username string) (*UserInfo, error) {
	opts := twitter.UserLookupOpts{
		UserFields: []twitter.UserField{twitter.UserFieldCreatedAt, twitter.UserFieldVerified, twitter.UserFieldPublicMetrics},
	}
	res, err := c.client.UserNameLookup(ctx, []string{username}, opts)
	if err != nil {
		// TODO: I used the global logger instance.
		// Change it to local loggers, or just use a global instance for all modules.
		// Also, try to set log-level inside the .env file
		// Here, one logger is simpler and better to manage
		logger.Error("user lookup error", "error", err)
		return nil, err
	}

	dictionaries := res.Raw.UserDictionaries()

	for _, userDic := range dictionaries {
		createdAt, err := time.Parse(time.RFC3339, userDic.User.CreatedAt)
		if err != nil {
			return nil, err
		}

		enc, _ := json.Marshal(userDic)
		logger.Debug("found user", "tweet", string(enc))

		userInfo := &UserInfo{
			TwitterID:   userDic.User.ID,
			TwitterName: username,
			CreatedAt:   createdAt,
			Followers:   userDic.User.PublicMetrics.Followers,
			IsVerified:  userDic.User.Verified,
		}
		return userInfo, nil
	}

	return nil, fmt.Errorf("no user found with %v", username)
}

func (c *Client) RetweetSearch(ctx context.Context, discordName string, username string) (*TweetInfo, error) {
	opts := twitter.TweetRecentSearchOpts{
		UserFields:  []twitter.UserField{twitter.UserFieldName},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt},
	}
	query := fmt.Sprintf("#Pactus AND %v from:%v is:quote", discordName, username)
	logger.Debug("search query", "query", query)

	res, err := c.client.TweetRecentSearch(ctx, query, opts)
	if err != nil {
		logger.Error("retweet lookup error", "error", err)
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("UserRetweetLookup result is nil")
	}
	if len(res.Raw.Tweets) == 0 {
		return nil, fmt.Errorf("no quote tweet with the hashtag '#Pactus' found. "+
			"Please select a tweet from https://x.com/PactusChain and retweet it. "+
			"Don't forget to add '#Pactus' and your Discord name '%v'. "+
			"For example you can tweet: \n"+
			"`#Pactus Blockchain is running a Twitter campaign to add more validators. Don't miss out!\n%v`", discordName, discordName)
	}

	quoteTweet := res.Raw.Tweets[0]

	enc, _ := json.Marshal(quoteTweet)
	logger.Debug("found quote tweet", "tweet", string(enc))

	createdAt, err := time.Parse(time.RFC3339, quoteTweet.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &TweetInfo{
		ID:        quoteTweet.ID,
		Link:      fmt.Sprintf("https://x.com/%v/status/%v", username, quoteTweet.ID),
		CreatedAt: createdAt,
	}, nil
}
