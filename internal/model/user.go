package model

import "time"

type User struct {
	Battletag     string
	Char          string
	Note          string
	Community     UserCommunity
	CommunityRank int
	AccessToken   string
	Expiry        time.Time
}

type UserCommunity struct {
	Id          string
	Name        string
	OfficerRank int
	Locked      bool
	Realm       string
}
type ValidatedUser struct {
	Battletag string             `json:"battletag"`
	Char      string             `json:"char"`
	Note      string             `json:"note"`
	IsAdmin   bool               `json:"isAdmin"`
	Community ValidatedCommunity `json:"community"`
}
type ValidatedCommunity struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Realm  string `json:"realm"`
	Locked bool   `json:"locked"`
}
