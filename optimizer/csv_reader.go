package optimizer

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func ReadConfigCSV(url string) (Config, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Config{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-OK HTTP status: %s", resp.Status)
	}

	r := csv.NewReader(resp.Body)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV data: %v", err)
	}

	// Check if data is empty or too short
	if len(records) < 2 {
		log.Fatal("CSV data is empty or missing headers.")
	}

	// --- PARSING THE DATA ---

	// Row 0 of the CSV (records[0]) corresponds to Sheet Row 1 (Plot Numbers).
	// Column 0 of the CSV (records[x][0]) corresponds to Sheet Column A (Player Names).

	// The player data starts at CSV row 1 (Sheet Row 2)
	playerDataStartRow := 1

	// The plot number headers start at CSV column 1 (Sheet Column B)
	plotHeaderStartCol := 1

	// The data values for weights start at CSV column 1 (Sheet Column B)
	weightsDataStartCol := 1

	// 3. Extract Plot Numbers (Header Row)
	// records[0] is the row containing plot numbers (Sheet Row 1)
	rawPlotHeaders := records[0]
	plotNumbers := make(map[int]int) // map[CSV_INDEX]Plot_Number

	for i := plotHeaderStartCol; i < len(rawPlotHeaders); i++ {
		// Attempt to parse the plot number from the header row (B1, C1, etc.)
		plotNum, err := strconv.Atoi(rawPlotHeaders[i])
		if err == nil && plotNum > 0 {
			// Store the plot number, keyed by its column index in the CSV array
			plotNumbers[i] = plotNum
		}
	}

	// 4. Extract Player Configs (Data Rows)
	var players []PlayerConfig
	// Iterate through each row starting from the first player row (Sheet Row 2)
	for rIdx := playerDataStartRow; rIdx < len(records); rIdx++ {
		row := records[rIdx]

		// Skip rows if the first column (Player Name) is empty
		playerName := row[0]
		if playerName == "" {
			continue
		}

		player := PlayerConfig{
			Name:    playerName,
			Weights: make(map[int]int),
		}

		var seenWeights map[int]bool

		// Iterate through the columns containing weights
		for cIdx := weightsDataStartCol; cIdx < len(row); cIdx++ {
			// Check if we have a valid plot number header for this column
			plotNum, ok := plotNumbers[cIdx]
			if !ok {
				continue // Skip if no valid plot number header exists for this column
			}

			// Get the weight value (which is a number or possibly a blank string)
			rawWeight := row[cIdx]

			// Parse the weight value. Atoi returns 0 for a blank string, which is fine
			// for the scenario where your formula blanks out cells.
			weight, err := strconv.Atoi(rawWeight)
			if err != nil {
				// If the value isn't a valid integer (e.g., if it's the blank "" from the formula),
				// treat it as max value
				weight = 50
			}

			if seenWeights == nil {
				seenWeights = make(map[int]bool)
			}

			if (weight != len(plotNumbers) && seenWeights[weight]) || weight > 50 || weight < 1 {
				// Duplicate detected — log, track, or flag this player
				fmt.Printf("⚠️ Invalid (duplication or range violation) weight %d found for player %s (plot %d)\n",
					weight, player.Name, plotNum)
				player.Cheater = true
				player.Weights = make(map[int]int)
				break
			} else {
				seenWeights[weight] = true
				player.Weights[plotNum] = weight
			}
		}

		players = append(players, player)
	}

	// 5. Construct Final Config
	finalConfig := Config{
		Plots:   len(plotNumbers),
		Players: players,
	}

	// PrintPlayerConfigs(finalConfig)

	return finalConfig, nil
}

func PrintPlayerConfigs(finalConfig Config) {
	fmt.Println("\n=======================================================")
	fmt.Printf("Parsed Configuration: %d Players, %d Total Plots\n", len(finalConfig.Players), finalConfig.Plots)
	fmt.Println("=======================================================")

	for _, player := range finalConfig.Players {
		// Prepare a list of plots that DO NOT have weight 50
		var weightedPlots []struct {
			Plot   int
			Weight int
		}

		// Sort the plots by number for readable output
		// We'll use a simple approach by iterating from 1 to 50
		for plotNum := 1; plotNum <= finalConfig.Plots; plotNum++ {
			weight, ok := player.Weights[plotNum]

			// Only include plots that were in the sheet (ok) AND whose weight is NOT 50
			if ok && weight != 50 {
				weightedPlots = append(weightedPlots, struct {
					Plot   int
					Weight int
				}{
					Plot:   plotNum,
					Weight: weight,
				})
			}
		}

		fmt.Printf("👤 Player: %s\n", player.Name)

		if len(weightedPlots) == 0 {
			fmt.Println("  - All plots have a weight of 50 or are unweighted.")
		} else {
			fmt.Println("  - Weighted Plots (excluding weight 50):")
			for _, item := range weightedPlots {
				// Use a formatted string to align the plot and weight
				fmt.Printf("    - Plot %02d: Weight %d\n", item.Plot, item.Weight)
			}
		}
		fmt.Println("-------------------------------------------------------")
	}
}
