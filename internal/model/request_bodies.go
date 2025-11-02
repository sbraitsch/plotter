package model

type PlayerUpdateRequest struct {
	Note     string      `json:"note"`
	PlotData map[int]int `json:"plotData"`
}

type CommunityRankRequest struct {
	AdminRank  int `json:"adminRank"`
	MemberRank int `json:"memberRank"`
}

type AssignmentUpload struct {
	Members []struct {
		Assignment Assignment `json:"assignment"`
	} `json:"members"`
}
