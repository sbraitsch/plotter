package optimizer

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/mattn/go-runewidth"
)

func PrettyPrint(v any) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal model for printing: %v", err)
	}

	fmt.Println(string(out))
}

func PrintOptimizedAssignments(assignments []Assignment, plotCount int, happyRange int) {
	reset := "\033[0m"
	for _, ass := range assignments {
		color := colorForScore(ass.Score, happyRange)
		if ass.Score == plotCount {
			fmt.Printf("%s#%d - %s ← This player did not get a prioritized spot and might be unhappy%s\n", color, ass.Plot, ass.Player, reset)
		} else {
			fmt.Printf("%s#%d - %s[%d]%s\n", color, ass.Plot, ass.Player, ass.Score, reset)
		}
	}
}

func colorForScore(score, happyRange int) string {
	if score <= 0 {
		score = 1
	}
	if happyRange < 1 {
		happyRange = 1
	}

	if score > happyRange {
		// Beyond happy range → bright red
		return "\033[38;2;255;0;0m"
	}

	t := float64(score-1) / float64(happyRange-1)
	if happyRange == 1 {
		t = 0
	}

	// We'll blend in two stages:
	// 0.0 → 0.5: green → yellow
	// 0.5 → 1.0: yellow → orange-red
	var r, g, b int

	if t < 0.5 {
		// Green (0,255,0) → Yellow (255,255,0)
		f := t / 0.5
		r = int(255 * f)
		g = 255
		b = 0
	} else {
		// Yellow (255,255,0) → Orange (255,128,0)
		f := (t - 0.5) / 0.5
		r = 255
		g = int(255 - 127*f)
		b = 0
	}

	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

func PrintHeader(score Score) {
	lines := []string{
		" Optimized for: Plot Priority ",
		fmt.Sprintf(" Total Score: %d ", score.total),
		fmt.Sprintf(" Happy Players: %d ", score.happyness),
		fmt.Sprintf(" Pleased Players: %d ", score.contentness),
		fmt.Sprintf(" Sad Players: %d ", score.sadness),
	}

	// find widest line for box width
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}

	// pad all lines to same width
	for i, l := range lines {
		if len(l) < width {
			lines[i] = l + strings.Repeat(" ", width-len(l))
		}
	}

	// box drawing
	h := "─"
	top := "┌" + strings.Repeat(h, width) + "┐"
	bot := "└" + strings.Repeat(h, width) + "┘"

	fmt.Println(top)
	for _, l := range lines {
		fmt.Println("│" + l + "│")
	}
	fmt.Println(bot)
	fmt.Println()
}

func PrintCheaters(cheaters []string) {
	if len(cheaters) == 0 {
		fmt.Println("No cheaters detected. ✅")
		return
	}

	lines := []string{
		" Cheaters ",
	}

	// add each cheater name on its own line
	for _, name := range cheaters {
		lines = append(lines, fmt.Sprintf(" - %s", name))
	}

	// find widest line for box width
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}

	// pad all lines to same width
	for i, l := range lines {
		if len(l) < width {
			lines[i] = l + strings.Repeat(" ", width-len(l))
		}
	}

	// box drawing
	h := "─"
	top := "┌" + strings.Repeat(h, width) + "┐"
	bot := "└" + strings.Repeat(h, width) + "┘"

	fmt.Println(top)
	for _, l := range lines {
		fmt.Println("│" + l + "│")
	}
	fmt.Println(bot)
	fmt.Println()
}

func repeat(s string, n int) string {
	res := ""
	for range n {
		res += s
	}
	return res
}

func PrintAssignmentTable(assignments []Assignment, plotCount, happyRange int) {
	// Sort assignments by plot
	sort.Slice(assignments, func(i, j int) bool {
		return assignments[i].Plot < assignments[j].Plot
	})

	// Determine column widths
	playerWidth := len("Player")
	plotWidth := len("Plot")
	scoreWidth := len("Priority")
	for _, a := range assignments {
		if len(a.Player) > playerWidth {
			playerWidth = len(a.Player)
		}
	}

	// Box drawing characters
	h := "─"
	topLeft, topRight := "┌", "┐"
	botLeft, botRight := "└", "┘"
	vert, midVert := "│", "┼"

	// Total table width (sum of columns + separators + spaces)
	totalWidth := playerWidth + plotWidth + scoreWidth + 8

	// Top border
	fmt.Println(topLeft + repeat(h, totalWidth) + topRight)

	// Column headers
	fmt.Printf("%s %-*s %s %-*s %s %-*s %s\n",
		vert, playerWidth, "Player",
		vert, plotWidth, "Plot",
		vert, scoreWidth, "Priority",
		vert)

	// Header separator
	fmt.Println(midVert + repeat(h, playerWidth+2) + midVert + repeat(h, plotWidth+2) + midVert + repeat(h, scoreWidth+2) + midVert)

	reset := "\033[0m"

	// Print rows
	for _, a := range assignments {
		color := colorForScore(a.Score, happyRange)
		// Use padRight on plain number to align correctly
		var scoreStr string
		if a.Score == plotCount {
			scoreStr = padRight("💀", scoreWidth)
		} else {
			scoreStr = padRight(fmt.Sprintf("%d", a.Score), scoreWidth)
		}
		fmt.Printf("%s %-*s %s %-*d %s %s%s%s %s\n",
			vert, playerWidth, a.Player,
			vert, plotWidth, a.Plot,
			vert, color, scoreStr, reset,
			vert)
	}

	// Bottom border
	fmt.Println(botLeft + repeat(h, totalWidth) + botRight)
}

func padRight(s string, width int) string {
	for runewidth.StringWidth(s) < width {
		s += " "
	}
	return s
}
