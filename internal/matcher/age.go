package matcher

import "time"

func AgeMultiplier(lastUsed int64, now time.Time) float64 {
	age := now.Sub(time.Unix(lastUsed, 0))

	if age < time.Hour {
		return 4.0
	} else if age < 24*time.Hour {
		return 2.0
	} else if age < 7*24*time.Hour {
		return 0.5
	}

	return 0.25
}
