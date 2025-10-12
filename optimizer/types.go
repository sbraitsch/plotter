package optimizer

type Config struct {
	Plots     int            `mapstructure:"plots"`
	Adjacency map[int][]int  `mapstructure:"adjacency"`
	Players   []PlayerConfig `mapstructure:"players"`
}

type PlayerConfig struct {
	Name    string      `mapstructure:"name"`
	Weights map[int]int `mapstructure:"weights"`
	Cheater bool
}

type Assignment struct {
	Player  string
	Plot    int
	Score   int
	Cheater bool
}

type Score struct {
	total       int
	sadness     int
	happyness   int
	contentness int
	cheaters    []string
}

func (cfg *Config) MaxWeights() int {
	maxLen := 0
	for _, p := range cfg.Players {
		if l := len(p.Weights); l > maxLen {
			maxLen = l
		}
	}
	return maxLen
}

func TotalScore(assignments []Assignment, plotCount int) Score {
	sum := 0
	sadness := 0
	happyness := 0
	contentness := 0
	cheaters := []string{}

	for _, a := range assignments {
		switch a.Score {
		case plotCount:
			sadness++
		case 1:
			happyness++
		default:
			contentness++
		}
		if a.Cheater {
			cheaters = append(cheaters, a.Player)
		} else {
			sum += (plotCount - a.Score)
		}
	}
	return Score{sum, sadness, happyness, contentness, cheaters}
}
