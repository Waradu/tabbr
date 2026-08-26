package matcher

import (
	"sort"
	"time"

	"github.com/waradu/tabbr/internal/store"
)

type RankedCommand struct {
	Command store.Command
	Score   float64
}

func Rank(commands []store.Command) []RankedCommand {
	now := time.Now()
	ranked := make([]RankedCommand, 0, len(commands))
	for _, command := range commands {
		ranked = append(ranked, RankedCommand{
			Command: command,
			Score:   float64(command.UseCount) * AgeMultiplier(command.LastUsed, now),
		})
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	return ranked
}
