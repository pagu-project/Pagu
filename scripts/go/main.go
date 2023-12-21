package main

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/kehiy/RoboPac/discord"
)

func main() {
	referrals, err := LoadReferralData("./referral.json")
	if err != nil {
		panic(err)
	}

	sort.Slice(referrals, func(i, j int) bool { // sort referrals based on points in descending order
		return referrals[i].Points > referrals[j].Points
	})

	top10 := [10]discord.Referral{}
	for i, r := range referrals {
		if i >= 10 {
			break
		}
		top10[i] = r
	}

	data, err := json.Marshal(top10)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("top10.json", data, 0o600)
	if err != nil {
		panic(err)
	}
}

func LoadReferralData(path string) ([]discord.Referral, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	referrals := map[string]discord.Referral{}
	err = json.Unmarshal(file, &referrals)
	if err != nil {
		return nil, err
	}

	result := []discord.Referral{}
	for _, r := range referrals {
		result = append(result, r)
	}

	return result, nil
}
