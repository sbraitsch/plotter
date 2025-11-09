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

type SingleAssignmentRequest struct {
	Battletag string `json:"btag"`
	Char      string `json:"char"`
	PlotId    int    `json:"plot"`
}
