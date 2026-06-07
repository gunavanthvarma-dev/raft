package raft

import (
	"math/rand"
	"time"
)

type TimeoutGenerator struct {
	MinimumTicks int64
	MaximumTicks int64
	randGen      *rand.Rand
}

func NewTimeoutGenerator(MinimumTicks int64, MaximumTicks int64) *TimeoutGenerator {
	return &TimeoutGenerator{MinimumTicks: MinimumTicks, MaximumTicks: MaximumTicks, randGen: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (timeoutGen *TimeoutGenerator) GenerateTimeout() uint64 {
	return uint64(timeoutGen.MinimumTicks) + uint64(timeoutGen.randGen.Int63n(timeoutGen.MaximumTicks-timeoutGen.MinimumTicks+1))
}
