package worm

import (
	"math/rand"
)

type MockWorm struct{}

func (m MockWorm) NoseTouch(val int) int {
	if val < 0 {
		val = 0
	}
	return rand.Intn(val + 1)
}

func (m MockWorm) Chemotaxis(val int) int {
	if val < 0 {
		val = 0
	}
	return rand.Intn(val + 1)
}
