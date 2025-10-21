package model

type PlotMappingRequest struct {
	PlotData map[int]int `json:"plotData"`
}

type CommunityRankRequest struct {
	MinRank int `json:"minRank"`
}
