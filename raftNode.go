package raft

type NodeId uint64

type LogEntry struct {
	Index   uint64
	Term    uint64
	Command []byte
}

type RaftState uint8

const (
	Follower  RaftState = iota //0
	Candidate                  //1
	Leader                     //2
)

type RaftNode struct {
	currentNodeId    NodeId
	peers            []NodeId
	electionTimeout  uint64 // random election timeout threshold
	electionElapsed  uint64 //ticks since hearing from the leader
	nodeStatus       RaftState
	heartbeatElapsed uint64 //ticks since last haertbeat; for leader
	heartbeatTimeout uint64 // frequency of heartbeats; for leader

	//persistent state
	currentTerm uint64 // latest term server has seen(default 0 at boot)
	votedFor    NodeId //candidateID that received vote in current term(default 0 if none)
	log         []LogEntry
	//volatile state
	commitIndex uint64 //index of highest log entry known to be committed
	lastApplied uint64 //index of highest log entry applied to state machine

	//volatile state on all leaders; Reinitialized after each election
	nextIndex map[NodeId]uint64 //for each server, index of the next log entry to send to that server(initizzalised to leader last log index+1 )

	matchIndex map[NodeId]uint64 //for each server, index of the highest log entry known to be replicated on server
}

//
