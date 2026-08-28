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

func TestRequestVoteResponseTrue(t *testing.T) {

	timeoutGen := NewFixedTimeoutGenerator(5)
	testNode := NewRaftNode(1, []NodeId{2, 3}, 5, 3, timeoutGen)

	testNode.Tick()
	testNode.Tick()

	reqVote := &Message{Type: 0, FromNodeId: 2, ToNodeId: 1, Term: 1, CandidateId: 2, LastLogIndex: 0, LastLogTerm: 0}

	testNode.ProcessNetworkMessage(*reqVote)
	tasks := testNode.Ready()
	msgs := tasks.Messages
	require.Equal(t, 1, len(msgs))
	assert.Equal(t, RequestVoteResponse, msgs[len(msgs)-1].Type)
	assert.Equal(t, true, msgs[len(msgs)-1].VoteGranted)
	assert.Equal(t, Follower, testNode.NodeStatus)

}

func TestRequestVoteResponseFalse(t *testing.T){
	timeoutGen := NewFixedTimeoutGenerator(5)
	testNode := NewRaftNode(1, []NodeId{2, 3}, 5, 3, timeoutGen)
	testNode.currentTerm = 2

	testNode.Tick()
	testNode.Tick()

	reqVote := &Message{Type: 0, FromNodeId: 2, ToNodeId: 1, Term: 1, CandidateId: 2, LastLogIndex: 0, LastLogTerm: 0}

	testNode.ProcessNetworkMessage(*reqVote)
	tasks := testNode.Ready()
	msgs := tasks.Messages
	require.Equal(t, 1, len(msgs))
	assert.Equal(t, RequestVoteResponse, msgs[len(msgs)-1].Type)
	assert.Equal(t, false, msgs[len(msgs)-1].VoteGranted)
	assert.Equal(t, Follower, testNode.NodeStatus)
}

func TestCandidateReceiveVoteFromAllNodesAndConvertToLeaderAndSendInitialAppendEntriesRPC(t *testing.T){
	timeoutGen := NewFixedTimeoutGenerator(3)
	testNode := NewRaftNode(1, []NodeId{2, 3}, 2, 3, timeoutGen)

	testNode.Tick();
	// call ready
	testNode.Ready()
	// call Advance
	testNode.Advance()
	testNode.Tick();
	testNode.Ready()
	testNode.Advance()
	reqVoteResp1 := &Message{Type:RequestVoteResponse,FromNodeId: 2,ToNodeId: 1,Term: 0,VoteGranted: true}
	reqVoteResp2 := &Message{Type:RequestVoteResponse,FromNodeId: 3,ToNodeId: 1,Term: 0,VoteGranted: true}
	testNode.ProcessNetworkMessage(*reqVoteResp1)
	tasks:= testNode.Ready()
	msgs:= tasks.Messages
	require.Equal(t, 2, len(msgs))
	assert.Equal(t, AppendEntriesRequest, msgs[len(msgs)-1].Type)
	assert.Equal(t,0,len(msgs[len(msgs)-1].Entries))
	assert.Equal(t,uint64(1),msgs[len(msgs)-1].Term)
	assert.Equal(t,uint64(0),msgs[len(msgs)-1].LeaderCommit)
	testNode.Advance()
	testNode.Tick()
	tasks=testNode.Ready()
	assert.Equal(t,0,len(tasks.Messages))
	testNode.Advance()
	testNode.ProcessNetworkMessage(*reqVoteResp2)
	tasks = testNode.Ready()
	assert.Equal(t,0,len(tasks.Messages))
	testNode.Advance()
	// call Advance
	
}