package raft

type NodeId uint64

type LogEntry struct {
	Index   uint64
	Term    uint64
	Command []byte
	//add LeaderCommit?
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
	logIndex    uint64 //next Index to insert in the log
	commitIndex uint64 //index of highest log entry known to be committed(default 0)
	lastApplied uint64 //index of highest log entry applied to state machine(default 0)
	leaderId    NodeId //leaderId
	//volatile state on all leaders; Reinitialized after each election
	nextIndex map[NodeId]uint64 //for each server, index of the next log entry to send to that server(initizzalised to leader last log index+1 )

	matchIndex   map[NodeId]uint64 //for each server, index of the highest log entry known to be replicated on server
	timeoutGen   *TimeoutGenerator //generate timeout
	serverTasks  *ServerTasks
	majority     uint64 //tracks the number of AppendEntriesResponse received to advance commitIndex
	leaderCommit uint64 // leader's commitIndex
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
func NewRaftNode(nodeid NodeId, peers []NodeId, electionTimeout uint64, heartbeatTimeout uint64, timeoutGen *TimeoutGenerator) *RaftNode {

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
		logIndex:         0,
		commitIndex:      0,
		lastApplied:      0,
		nextIndex:        make(map[NodeId]uint64),
		matchIndex:       make(map[NodeId]uint64),
		timeoutGen:       timeoutGen,
		serverTasks:      new(ServerTasks),
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

	if node.NodeStatus == Leader {
		node.HeartbeatElapsed += 1
		if node.HeartbeatElapsed == node.HeartbeatTimeout {
			node.sendAppendEntries()
			node.HeartbeatElapsed = 0
		}
	} else {
		node.ElectionElapsed += 1
		if node.ElectionElapsed == node.ElectionTimeout {
			node.ElectionElapsed = 0
			node.ElectionVotes = 0
			node.ElectionTimeout = node.timeoutGen.GenerateTimeout() //generate new timeout
			node.currentTerm += 1
			node.NodeStatus = Candidate
			node.votedFor = node.CurrentNodeId
			node.ElectionVotes += 1
			//need a function that builds RequestVoteMessage to send to all peers
			//add Message to ServerTasks
			node.sendRequestVote()
		}
	}

}

func (node *RaftNode) sendRequestVote() {
	for _, peerId := range node.Peers {
		msg := &Message{Type: RequestVoteRequest,
			FromNodeId:   node.CurrentNodeId,
			ToNodeId:     peerId,
			Term:         node.currentTerm,
			CandidateId:  node.CurrentNodeId,
			LastLogIndex: node.log[len(node.log)-1].Index,
			LastLogTerm:  node.log[len(node.log)-1].Term,
		}
		node.serverTasks.Messages = append(node.serverTasks.Messages, *msg)
	}
}

func (node *RaftNode) sendAppendEntries() {
	if node.NodeStatus != Leader {
		//log error
		return
	}
	for _, peerId := range node.Peers {
		var msg *Message
		if node.log[len(node.log)-1].Index >= node.nextIndex[peerId] {

			msg = &Message{Type: AppendEntriesRequest,
				FromNodeId:   node.CurrentNodeId,
				ToNodeId:     peerId,
				Term:         node.currentTerm,
				LeaderId:     node.CurrentNodeId,
				PrevLogIndex: node.log[node.nextIndex[peerId]-1].Index,
				PrevLogTerm:  node.log[node.nextIndex[peerId]-1].Term,
				Entries:      node.log[node.nextIndex[peerId]:],
				LeaderCommit: node.commitIndex,
			}
		} else {
			msg = &Message{Type: AppendEntriesRequest,
				FromNodeId:   node.CurrentNodeId,
				ToNodeId:     peerId,
				Term:         node.currentTerm,
				LeaderId:     node.CurrentNodeId,
				PrevLogIndex: 0,
				PrevLogTerm:  0,
				Entries:      make([]LogEntry, 0),
				LeaderCommit: node.commitIndex,
			}
			//send Heartbeat
		}
		node.serverTasks.Messages = append(node.serverTasks.Messages, *msg) //append messages to serverTasks
	}

}

func (node *RaftNode) ProcessClientRequest(req ClientRequest) {
	//if current node is follower, reject request and send the latest leader info (need to decide on how to do it;add to ServerTasks)
	//create a ClientRedirect message, give leader Id and success as false
	//if current node is candidate, reject request(send an error response and tell to try again;add to ServerTasks)
	//create a ClientError message, and put success as false
	//if current node is leader, process the request
	//create new LogEntry
	//append logentry to leader's log
	//persist; add to ServerTasks
	//Build AppendEntriesRPC for all peers
	//send it to all peers(through serverTasks)
	//so the leader has to wait for majority replication before applyting the entry to the state machine(do we need a state variable to track whether the leader is waiting for majority replication and also track the majority count)

	switch node.NodeStatus {
	case Follower:
		msg := &Message{Type: ClientRedirect, FromNodeId: node.CurrentNodeId, LeaderId: node.leaderId, Success: false}
		node.serverTasks.Messages = append(node.serverTasks.Messages, *msg)

	case Candidate:
		msg := &Message{Type: ClientError, FromNodeId: node.CurrentNodeId, Success: false}
		node.serverTasks.Messages = append(node.serverTasks.Messages, *msg)

	case Leader:
		node.logIndex += 1
		logEntry := &LogEntry{Index: node.logIndex, Term: node.currentTerm, Command: req.Data}
		node.log = append(node.log, *logEntry)
		node.serverTasks.EntriesToPersist = append(node.serverTasks.EntriesToPersist, *logEntry)
		node.sendAppendEntries() //send AppendEntriesRPC to all followers
	}

}

func (node *RaftNode) appendEntriesResponseFalse(destNodeId NodeId) {
	msg := &Message{
		Type:       AppendEntriesResponse,
		FromNodeId: node.CurrentNodeId,
		ToNodeId:   destNodeId,
		Term:       node.currentTerm,
		Success:    false,
	}
	node.serverTasks.Messages = append(node.serverTasks.Messages, *msg)
}

func (node *RaftNode) appendEntriesResponseTrue(lastEntryIndex uint64) {
	msg := &Message{
		Type:           AppendEntriesResponse,
		FromNodeId:     node.CurrentNodeId,
		ToNodeId:       node.leaderId,
		Term:           node.currentTerm,
		Success:        true,
		LastEntryIndex: lastEntryIndex,
	}
	node.serverTasks.Messages = append(node.serverTasks.Messages, *msg)
}

func (node *RaftNode) RemoveEntriesAndAppend(entries []LogEntry) {

	//remove entries from entries[0].index till end
	tempLog := node.log[:entries[0].Index]
	for _, val := range entries {
		tempLog[val.Index] = val
	}
	node.log = tempLog
}

func (node *RaftNode) AddEntriesToPersist() {
	for _, val := range node.log[node.commitIndex:] {
		node.serverTasks.EntriesToPersist = append(node.serverTasks.EntriesToPersist, val)
	}
}

func (node *RaftNode) RequestVoteResponseFalse(destNodeId NodeId) {
	msg := &Message{
		Type:        RequestVoteResponse,
		FromNodeId:  node.CurrentNodeId,
		ToNodeId:    destNodeId,
		Term:        node.currentTerm,
		VoteGranted: false,
	}
	node.serverTasks.Messages = append(node.serverTasks.Messages, *msg)
}

func (node *RaftNode) RequestVoteResponseTrue(destNodeId NodeId) {
	msg := &Message{
		Type:        RequestVoteResponse,
		FromNodeId:  node.CurrentNodeId,
		ToNodeId:    destNodeId,
		Term:        node.currentTerm,
		VoteGranted: true,
	}
	node.serverTasks.Messages = append(node.serverTasks.Messages, *msg)
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
	//		node.leaderId = msg.FromId
	//		check = Entries[0]
	//		if check.term!= log[check.index].term
	//			remove everything from current entry to end of log
	//		append new entries to the log; add to ServerTasks                  ??PERSIST TO DISK FIRST BEFORE SENDING RESPONSE TO LEADER
	//		if leaderCommit>commitIndex                                           // do this in Advance()
	//			commitIndex = min(leaderCommit,index of last new entry)
	//		build AppendEntriesResponse; add to ServerTasks                       // this as well
	//1.3Candidate
	//  	1.3.1 if msg.term < currentTerm
	//			1.3.1.1 success = false
	//			1.3.1.2 term = currentTerm
	//			1.3.1.3 create AppendEntriesResponse and add to serverTasks
	//		1.3.2 else
	//			1.3.2.1 node.leaderId = msg.FromId
	//			1.3.2.2 change nodeStatus to Follower; and process the message as a follower   //redundant
	//1.4 Leader
	//1.4.1 if msg.term < currentTerm
	//		same as above
	//1.4.2 else
	//		convert to follower; process the message as a follower //redundant

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

	if msg.Term > node.currentTerm {
		node.NodeStatus = Follower
		node.leaderId = msg.FromNodeId // shoudnt this be like, the node that sent the msg should include the node it thinks is the leader? instead of assuming the message is sent by the leader
	}

	switch msg.Type {
	case AppendEntriesRequest:
		switch node.NodeStatus {
		case Follower:
			//func to process AppendEntriesRPC as a Follower
			// Step-1: check if the msg.term < currentTerm
			//if yes, send appendEntriesResponse with currentTerm and false
			//else
			//if log[prevLogIndex].term != msg.term
			//send appendEntriesResponse with currentTerm and false
			//break
			//clear log from log[msg.entries.start:]
			//append msg.entries
			// build entriestoPersist; add to ServerTasks;

			if len(msg.Entries) == 0 {
				//heartbeat
				node.ElectionElapsed = 0
				break
			}
			if msg.Term < node.currentTerm || msg.PrevLogIndex > uint64(len(node.log)) || node.log[msg.PrevLogIndex].Term != msg.Term {
				node.appendEntriesResponseFalse(msg.FromNodeId)
			} else {
				node.leaderCommit = msg.LeaderCommit
				node.RemoveEntriesAndAppend(msg.Entries)
				node.AddEntriesToPersist()
				node.ElectionElapsed = 0
			}
		case Candidate:
			if msg.Term < node.currentTerm {
				node.appendEntriesResponseFalse(msg.FromNodeId)
			} else {
				//A leader has been elected, convert to follower
				node.leaderCommit = msg.LeaderCommit
				node.NodeStatus = Follower
				node.ElectionElapsed = 0
				// do I have to rest anything?
			}
		case Leader:
			if msg.Term < node.currentTerm {
				node.appendEntriesResponseFalse(msg.FromNodeId)
			} else {
				//shoudnt happen as there will be 2 leaders in a single term. Invariant violation
				//LOG ERROR
			}

		}
	case AppendEntriesResponse:
		switch node.NodeStatus {
		case Follower:
			if msg.Term < node.currentTerm {
				node.appendEntriesResponseFalse(msg.FromNodeId)
			} else {
				//not leader; so Ignore
				//LOG Exception
			}

		case Candidate:
			if msg.Term < node.currentTerm {
				node.appendEntriesResponseFalse(msg.FromNodeId)
			} else {
				//not leader; so Ignore
				//LOG Exception
			}

		case Leader:
			if msg.Term < node.currentTerm {
				node.appendEntriesResponseFalse(msg.FromNodeId)
			} else {
				if msg.Success == true {
					node.nextIndex[msg.FromNodeId] += msg.LastEntryIndex + 1
					node.matchIndex[msg.FromNodeId] = msg.LastEntryIndex
				} else {
					node.nextIndex[msg.FromNodeId] -= 1
					node.sendAppendEntries()
				}
			}

		}
	case RequestVoteRequest:
		if msg.Term < node.currentTerm {
			node.RequestVoteResponseFalse(msg.FromNodeId)
		}
		switch node.NodeStatus {
		case Follower:
			lastLogEntry := node.log[len(node.log)-1]
			if (node.votedFor == 0 || node.votedFor == msg.FromNodeId) && lastLogEntry.Term <= msg.LastLogTerm {
				if lastLogEntry.Term == msg.LastLogTerm {
					if lastLogEntry.Index <= msg.LastLogIndex {
						node.RequestVoteResponseTrue(msg.FromNodeId)
					} else {
						node.RequestVoteResponseFalse(msg.FromNodeId)
					}
				} else {
					node.RequestVoteResponseTrue(msg.FromNodeId)
				}
			} else {
				node.RequestVoteResponseFalse(msg.FromNodeId)
			}
		case Candidate:
			node.RequestVoteResponseFalse(msg.FromNodeId)
			// what if the candidate requesting vote has a log that is more up-to-date?
		case Leader:
			// the message wont be processed as a leader
			// LOG error
		}

	case RequestVoteResponse:
		switch node.NodeStatus {
		case Follower:
		//Ignore or send a msg?
		case Candidate:
			if msg.VoteGranted == true {
				node.ElectionVotes += 1
			}
			if node.ElectionVotes >= uint64(len(node.Peers)/2) {
				node.NodeStatus = Leader
			}
		case Leader:
			//Ignore or send a msg?
		}

	}

}

func (node *RaftNode) Ready() ServerTasks {
	//if commitIndex>lastApplied
	//add log[lastApplied] to FSM; add to serverTasks;
	//send ServerTasks() to server

	if node.commitIndex > node.lastApplied {
		node.serverTasks.EntriesToApply = append(node.serverTasks.EntriesToApply, node.log[node.lastApplied+1])
	}
	for _, val := range node.serverTasks.EntriesToPersist {
		node.appendEntriesResponseTrue(val.Index)
	}
	return *node.serverTasks
}

// server tells Raft that everything is done
func (node *RaftNode) Advance() {
	//for every entry in EntriesToPersist
	// send AppendEntriesResponse to leader
	//		if leaderCommit>commitIndex
	//			commitIndex = min(leaderCommit,index of last new entry)
	//for every entry in EntriesToApply
	//++lastApplied

	clear(node.serverTasks.Messages)
	// for _, val := range node.serverTasks.EntriesToPersist {
	// 	node.appendEntriesResponseTrue(val.Index)
	// 	if node.leaderCommit > node.commitIndex {
	// 		node.commitIndex = min(node.leaderCommit, val.Index)
	// 	}
	// }
	if node.leaderCommit > node.commitIndex {
		node.commitIndex = min(node.leaderCommit, node.serverTasks.EntriesToPersist[len(node.serverTasks.EntriesToPersist)-1].Index)
	}
	clear(node.serverTasks.EntriesToPersist)

	node.lastApplied += uint64(len(node.serverTasks.EntriesToApply))
	clear(node.serverTasks.EntriesToApply)
}
