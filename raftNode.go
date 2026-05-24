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
	CurrentNodeId    NodeId
	Peers            []NodeId
	ElectionTimeout  uint64 // random election timeout threshold
	ElectionElapsed  uint64 //ticks since hearing from the leader
	NodeStatus       RaftState
	HeartbeatElapsed uint64 //ticks since last haertbeat; for leader
	HeartbeatTimeout uint64 // frequency of heartbeats; for leader
	//persistent state
	currentTerm uint64 // latest term server has seen(default 0 at boot)
	votedFor    NodeId //candidateID that received vote in current term(default 0 if none)
	log         []LogEntry
	//volatile state
	commitIndex uint64 //index of highest log entry known to be committed(default 0)
	lastApplied uint64 //index of highest log entry applied to state machine(default 0)

	//volatile state on all leaders; Reinitialized after each election
	nextIndex map[NodeId]uint64 //for each server, index of the next log entry to send to that server(initizzalised to leader last log index+1 )

	matchIndex map[NodeId]uint64 //for each server, index of the highest log entry known to be replicated on server
}

//Methods:

// 1. Create new RaftNode with config arguments NewRaftNode()
// 2. ProcessNetworkMessage ---> procesess the message received from the server based on the message type:-
//This will be a continuous running loop in the server that starts after Node creation
//should i put the ticker update here?
// 3. processAppendEntriesRequest()
// 4. processAppendEntriesResponse()
// 5. processRequestVoteRequest()
// 6. processRequestVoteResponse()
// 7. Tick()  ---> update time
// 8. ProcessClientRequest() ---> process client requests
// 9. ServerTasks() --- tells the server what it needs to do
// 10. Proceed() --- tells Raft that everything was done

func NewRaftNode(nodeid NodeId, peers []NodeId, electionTimeout uint64, heartbeatTimeout uint64) *RaftNode {

	return &RaftNode{
		CurrentNodeId:    nodeid,
		Peers:            peers,
		ElectionTimeout:  electionTimeout,
		ElectionElapsed:  0,
		NodeStatus:       Follower,
		HeartbeatElapsed: 0,
		HeartbeatTimeout: heartbeatTimeout,
		currentTerm:      0,
		votedFor:         0,
		commitIndex:      0,
		lastApplied:      0,
		nextIndex:        make(map[NodeId]uint64),
		matchIndex:       make(map[NodeId]uint64),
	}
}
