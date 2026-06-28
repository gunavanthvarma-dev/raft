package raft

import "fmt"

type MessageType uint8

const (
	RequestVoteRequest MessageType = iota
	RequestVoteResponse
	AppendEntriesRequest
	AppendEntriesResponse
	//Client Messages
	ClientRedirect
	ClientError
)

type Message struct {
	Type       MessageType
	FromNodeId NodeId
	ToNodeId   NodeId

	Term uint64 //common for all messages
	//AppendEntriesRequest
	LeaderId     NodeId
	PrevLogIndex uint64
	PrevLogTerm  uint64
	Entries      []LogEntry
	LeaderCommit uint64

	//AppendEntriesResponse
	Success        bool
	LastEntryIndex uint64

	//RequestVoteRequest
	CandidateId  NodeId
	LastLogIndex uint64
	LastLogTerm  uint64

	//RequestVoteResponse
	VoteGranted bool
}

func (msg *Message) ToString() string {
	return fmt.Sprintf("MessageType:%d\nFromNodeId:%d\nToNodeId:%d\nTerm:%d\nLeaderId:%d\nPrevLogIndex:%d\nPrevLogTerm:%d\nEntries:\nLeaderCommit:%d\nSuccess:%d\nLastEntryIndex:%d\nCandidateId:%d\nLastLogIndex:%d\nLastLogTerm:%d\nVoteGranted:%d", msg.Type, msg.FromNodeId, msg.ToNodeId, msg.Term, msg.LeaderId, msg.PrevLogIndex, msg.PrevLogTerm, msg.LeaderCommit, msg.Success, msg.LastEntryIndex, msg.CandidateId, msg.LastLogIndex, msg.LastLogTerm, msg.VoteGranted)
}

type ClientRequest struct {
	Data []byte
}

type ClientResponse struct {
	Data []byte
}

type PersistentState struct {
	currentTerm uint64
	votedFor    NodeId
}

type ServerTasks struct {
	Messages         []Message
	EntriesToPersist []LogEntry
	EntriesToApply   []LogEntry
	State            PersistentState
}
