package rafttest

import "gunavanthvarma-dev/raft"

type raftTest struct {
	NumberOfNodes uint64
	NodesList     []raft.RaftNode
	NodeState     map[raft.RaftState]NodeStateVariables
	// add a log file
}

//IDEA:
//create a config file, read from that config file. implement later

func newRaftTest(noOfNodes uint64) *raftTest {
	peerslist := make([]raft.NodeId, 0)
	NodesList := make([]raft.RaftNode, 0)
	timeoutGen := raft.NewTimeoutGenerator(5, 15)
	for v := range noOfNodes {
		peerslist = append(peerslist, raft.NodeId(v+1))
		node := raft.NewRaftNode(
			raft.NodeId(v+1),
			make([]raft.NodeId, 0),
			timeoutGen.GenerateTimeout(),
			4,
			timeoutGen,
		)
		NodesList = append(NodesList, *node)
	}
	return &raftTest{NumberOfNodes: noOfNodes, NodesList: NodesList}
}

type NodeStateVariables struct {
	ElectionTimeout  uint64 // initially unique for each node
	ElectionElapsed  uint64
	HeartbeatElapsed uint64
	HeartbeatTimeout uint64 //unique for each node?
	ElectionVotes    uint64

	currentTerm uint64
	votedFor    raft.NodeId
	log         []raft.LogEntry

	logIndex    uint64
	commitIndex uint64
	lastApplied uint64
	leaderId    raft.NodeId
	nextIndex   map[raft.NodeId]uint64

	matchIndex   map[raft.NodeId]uint64
	timeoutGen   *raft.TimeoutGenerator
	serverTasks  *raft.ServerTasks
	majority     uint64
	leaderCommit uint64
}

type TransportLayer struct {
	NodeTransportLayers map[raft.NodeId]NodeTransport
}

type NodeTransport struct {
	Inbound  chan *raft.Message
	Outbound chan *raft.Message
}

type StorageLayer struct {
	NodeStorageLayers map[raft.NodeId]NodeStorage
}

type NodeStorage struct {
	log         []raft.LogEntry
	currentTerm uint64
	votedFor    raft.NodeId
}
