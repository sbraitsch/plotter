package model

type CommunityData struct {
	Id      string       `json:"id"`
	Members []MemberData `json:"members"`
}

type MemberData struct {
	Character string      `json:"char"`
	BattleTag string      `json:"battletag"`
	PlotData  map[int]int `json:"plotData"`
}

type Community struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Realm  string `json:"realm"`
	Locked bool   `json:"locked"`
}

type Roster struct {
	Members []member `json:"members"`
}

type member struct {
	Character character `json:"character"`
	Rank      int       `json:"rank"`
}

type character struct {
	Name string `json:"name"`
}

type Assignment struct {
	Character string `json:"char"`
	Battletag string `json:"btag"`
	Plot      int    `json:"plot"`
	Score     int    `json:"score"`
}

type Settings struct {
	OfficerRank int `json:"officerRank"`
	MemberRank  int `json:"memberRank"`
}

type FullCommunityData struct {
	Id      string           `json:"id"`
	Members []FullMemberData `json:"members"`
}

type FullMemberData struct {
	Assignment Assignment  `json:"assignment"`
	Note       string      `json:"note"`
	PlotData   map[int]int `json:"plotSelection"`
}
