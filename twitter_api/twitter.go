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
	res, err := client.UserLookup(context.Background(), []string{twitterID}, opts)
	if err != nil {
		return nil, err
	}

	if len(res.Raw.Users) == 0 || res.Raw.Users[0] == nil {
		return nil, fmt.Errorf("unable to find %v", twitterID)
	}

	logger.Info("found twitter", "name", res.Raw.Users[0].Name, "id", twitterID)

	return &Client{
		client:    client,
		twitterID: twitterID,
	}, nil
}

func (c *Client) UserInfo(ctx context.Context, twitterName string) (*UserInfo, error) {
	opts := twitter.UserLookupOpts{
		UserFields: []twitter.UserField{twitter.UserFieldCreatedAt, twitter.UserFieldVerified, twitter.UserFieldPublicMetrics},
	}
	res, err := c.client.UserNameLookup(ctx, []string{twitterName}, opts)
	if err != nil {
		logger.Error("user lookup error", "error", err)
		return nil, err
	}
	if len(res.Raw.Errors) > 0 {
		return nil, fmt.Errorf("the Twitter API returned error: %v", res.Raw.Errors[0].Detail)
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
			TwitterName: twitterName,
			CreatedAt:   createdAt,
			Followers:   userDic.User.PublicMetrics.Followers,
			IsVerified:  userDic.User.Verified,
		}
		return userInfo, nil
	}

	return nil, fmt.Errorf("no user found with %v", twitterName)
}

func (c *Client) RetweetSearch(ctx context.Context, discordID string, twitterName string) (*TweetInfo, error) {
	opts := twitter.TweetRecentSearchOpts{
		UserFields:  []twitter.UserField{twitter.UserFieldName},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt},
	}
	query := fmt.Sprintf("%v (#Pactus pr #PactusBoosterProgram) from:%v is:quote", discordID, twitterName)
	logger.Debug("search query", "query", query)

	res, err := c.client.TweetRecentSearch(ctx, query, opts)
	if err != nil {
		logger.Error("retweet lookup error", "error", err)
		return nil, err
	}
	if len(res.Raw.Errors) > 0 {
		return nil, fmt.Errorf("the Twitter API returned error: %v", res.Raw.Errors[0].Detail)
	}

	if len(res.Raw.Tweets) == 0 {
		return nil, fmt.Errorf("no eligible quote tweet found. "+
			"Please select a tweet from https://x.com/PactusChain and retweet it. "+
			"Don't forget to add '#PactusBoosterProgram' and your Discord ID '%v'. "+
			"For example you can tweet: \n\n"+
			"```"+
			"#Pactus Blockchain just started the Validator Booster Program on Twitter Campaign. Don't miss out! https://discord.com/invite/H5vZkNnXCu\n"+
			"#PactusBoosterProgram %v"+
			"```", discordID, discordID)
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
		Link:      fmt.Sprintf("https://x.com/%v/status/%v", twitterName, quoteTweet.ID),
		CreatedAt: createdAt,
	}, nil
}
