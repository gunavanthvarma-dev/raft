package raft

import (
	"math/rand"
	"time"
)

type TimeoutGenerator interface {
	GenerateTimeout() uint64
}

type RandomTimeoutGenerator struct {
	MinimumTicks int64
	MaximumTicks int64
	randGen      *rand.Rand
}

func NewRandomTimeoutGenerator(MinimumTicks int64, MaximumTicks int64) *RandomTimeoutGenerator {
	return &RandomTimeoutGenerator{MinimumTicks: MinimumTicks, MaximumTicks: MaximumTicks, randGen: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (timeoutGen *RandomTimeoutGenerator) GenerateTimeout() uint64 {
	return uint64(timeoutGen.MinimumTicks) + uint64(timeoutGen.randGen.Int63n(timeoutGen.MaximumTicks-timeoutGen.MinimumTicks+1))
}

type FixedTimeoutGenerator struct {
	Ticks uint64
}

func NewFixedTimeoutGenerator(ticks uint64) *FixedTimeoutGenerator {
	return &FixedTimeoutGenerator{Ticks: ticks}
}

func (timeoutGen *FixedTimeoutGenerator) GenerateTimeout() uint64 {
	return timeoutGen.Ticks
}
