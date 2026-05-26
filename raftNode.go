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
	ElectionVotes    uint64 //count votes during election
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
// This will be a continuous running loop in the server that starts after Node creation
// should i put the ticker update here?
// 3. processAppendEntriesRequest()
// 4. processAppendEntriesResponse()
// 5. processRequestVoteRequest()
// 6. processRequestVoteResponse()
// 7. Tick()  ---> update time
// 8. ProcessClientRequest() ---> process client requests
// 9. ServerTasks() --- tells the server what it needs to do
// 10. Proceed() --- tells Raft that everything was done
// 11. candidateActions() ---
func NewRaftNode(nodeid NodeId, peers []NodeId, electionTimeout uint64, heartbeatTimeout uint64) *RaftNode {

	return &RaftNode{
		CurrentNodeId:    nodeid,
		Peers:            peers,
		ElectionTimeout:  electionTimeout,
		ElectionElapsed:  0,
		NodeStatus:       Follower,
		HeartbeatElapsed: 0,
		HeartbeatTimeout: heartbeatTimeout,
		ElectionVotes:    0,
		currentTerm:      0,
		votedFor:         0,
		commitIndex:      0,
		lastApplied:      0,
		nextIndex:        make(map[NodeId]uint64),
		matchIndex:       make(map[NodeId]uint64),
	}
}

func (node *RaftNode) Tick() {
	//increment ElectionElapsed
	//check if ElectionElapsed==ElectionTimeout
	//IF YES
	//ElectionElapsed=0
	//ElectionVotes=0
	//ElectionTimeout = Random ticks
	//increment currentTerm
	//change state to Candidate
	//votedFor = CurrentNodeId
	//increment ElectionVotes
	//Build RequestVoteRequest message
	//send it to all peers(through serverTasks)
	//IF NO

}

func (node *RaftNode) ProcessClientRequest(req ClientRequest) {
	//if current node is follower, reject request and send the latest leader info (need to decide on how to do it;add to ServerTasks)
	//if current node is candidate, reject request(send an error response and tell to try again;add to ServerTasks)
	//if current node is leader, process the request
	//create new LogEntry
	//append logentry to leader's log
	//Build AppendEntriesRPC for all peers
	//send it to all peers(through serverTasks)
	//so the leader has to wait for majority replication before applyting the entry to the state machine(do we need a state variable to track whether the leader is waiting for majority replication and also track the majority count)

}

func (node *RaftNode) ProcessNetworkMessage(msg Message) {
	// whatver message you get, if msg.term > currentTerm;  porcess the msg as a follower
	//	if entries[] is empty ---> HEARTBEAT
	//		ElectionElapsed=0
	//1.AppendEntriesRequest
	//1.1check NodeStatus before processing msg
	//1.2Follower
	//	1.2.1if msg.term < currentTerm
	//		1.2.1.1 success = false
	//		1.2.1.2 term = currentTerm
	//		1.2.1.3 create AppendEntriesResponse and add to serverTasks
	//		1.2.1.4 return
	//	1.2.2 if log[msg.prevLogIndex] not exist || log[msg.prevLogIndex].term!=prevLogTerm
	//	same as above
	//
	//		check = Entries[0]
	//		if check.term!= log[check.index].term
	//			remove everything from current entry to end of log
	//		append new entries to the log; add to ServerTasks
	//		if leaderCommit>commitIndex
	//			commitIndex = min(leaderCommit,index of last new entry)
	//		build AppendEntriesResponse; add to ServerTasks
	//1.3Candidate
	//  	1.3.1 if msg.term < currentTerm
	//			1.3.1.1 success = false
	//			1.3.1.2 term = currentTerm
	//			1.3.1.3 create AppendEntriesResponse and add to serverTasks
	//		1.3.2 else
	//			1.3.2.1 change nodeStatus to Follower; and process the message as a follower
	//1.4 Leader
	//1.4.1 if msg.term < currentTerm
	//		same as above
	//1.4.2 else
	//		convert to follower; process the message as a follower

	// 2. AppendEntriesResponse

	//	2.1 Leader
	//			2.1.1 if msg.success == true
	//					2.1.1.1 nextIndex[msg.FromNodeId]++
	//					2.1.1.2 matchIndex[msg.FromNodeId]++
	//			2.1.2 if msg.success == false
	//					2.1.2.1 nextIndex[msg.FromNodeId]--
	//					2.1.2.2 Build AppendEntriesRPC with log[nextIndex[msg.FromNodeId]]; add to ServerTasks
	//	2.2 Candidate
	//			2.2.1 Ignore the msg; cause its stale if the term is older than current term
	//  2.3 Follower
	//			2.3.1 Same as above

	// 3. RequestVoteRequest [Some restrictions will have to be placed for worse case scenarios]

	//	3.1  Leader
	//		if msg.Term > currentTerm; process as follower [restriction might be needed here]
	//
	//  3.2  Candidate
	//		if msg.Term == currentTerm
	//			build RequestVoteResponse and reject; add to ServerTasks
	//
	//  3.3  Follower
	//		if (votedFor = Null || msg.candidateID) && msg.lastLogIndex == lastIndex && msg.lastLogTerm == lastTerm
	//			build RequestVoteResponse as true; add to ServerTasks
	//		else
	//			build RequestVoteResponse as false; add to ServerTasks

	// 4. RequestVoteResponse

	//		3.1  Leader
	//				Ignore message unless msg.term>currentTerm
	//
	//		3.2  Candidate
	//				if msg.term == currentTerm && msg.voteGranted == true
	//					ElectionVotes++
	//				if ElectionVotes > peers.length/2
	//					Election won, convert to Leader
	//					send heartbeat to all peers
	//					Build AppendEntriesRPC; add to ServerTasks
	//		3.3  Follower
	//				Ignore it

}
