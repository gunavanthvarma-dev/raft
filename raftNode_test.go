package raft

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	//"github.com/stretchr/testify/require"
)

func TestElectionTimeoutStartsElection(t *testing.T) {

	timeoutGen := NewFixedTimeoutGenerator(3)
	testNode := NewRaftNode(1, []NodeId{2, 3}, 3, 3, timeoutGen)

	for range 3 {
		testNode.Tick()
	}

	tasks := testNode.Ready()
	msgs := tasks.Messages
	require.Equal(t, 2, len(msgs))
	assert.Equal(t, RequestVoteRequest, msgs[len(msgs)-1].Type)
	assert.Equal(t, uint64(1), testNode.currentTerm)
	assert.Equal(t, Candidate, testNode.NodeStatus)

}
