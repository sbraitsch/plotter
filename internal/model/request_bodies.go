package model

type PlotMappingRequest struct {
	PlotData map[int]int `json:"plotData"`
}

type CommunityRankRequest struct {
	AdminRank  int `json:"adminRank"`
	MemberRank int `json:"memberRank"`
}
