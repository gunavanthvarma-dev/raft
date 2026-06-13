package raft

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
	Success bool

	//RequestVoteRequest
	CandidateId  NodeId
	LastLogIndex uint64
	LastLogTerm  uint64

	//RequestVoteResponse
	VoteGranted bool
}

type ClientRequest struct {
	Data []byte
}

type ClientResponse struct {
	Data []byte
}

type ServerTasks struct {
	Messages         []Message
	EntriesToPersist []LogEntry
	EntriesToApply   []LogEntry
}
