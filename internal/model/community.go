package model

type CommunityData struct {
	Id      string       `json:"id"`
	Members []MemberData `json:"members"`
}

type MemberData struct {
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
	Battletag string `json:"player"`
	Plot      int    `json:"plot"`
	Score     int    `json:"score"`
}
