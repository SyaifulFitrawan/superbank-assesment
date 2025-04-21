package utils

import (
	"math/rand"
	"time"
)

func RandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomRoundedAmount() float64 {
	if RandomBool() {
		return float64(RandomInt(1, 9)) * 100_000
	} else {
		return float64(RandomInt(1, 10)) * 1_000_000
	}
}

func RandomBool() bool {
	return rand.Intn(2) == 1
}

func RandomDate(layout string) time.Time {
	start := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()
	delta := end.Unix() - start.Unix()
	sec := rand.Int63n(delta) + start.Unix()
	return time.Unix(sec, 0)
}
