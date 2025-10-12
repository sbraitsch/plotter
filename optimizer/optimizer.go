package optimizer

func Optimize(config Config) {

	// First pass
	assignments := RunHungarian(config)

	PrintAssignmentTable(assignments, config.Plots, 5 /* config.MaxWeights() */)
	score := TotalScore(assignments, config.Plots)
	PrintHeader(score)
	PrintCheaters(score.cheaters)
}
