package service

type PlayerData struct {
	Name     string      `json:"name"`
	PlotData map[int]int `json:"plotdata"`
}

type Player struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}
