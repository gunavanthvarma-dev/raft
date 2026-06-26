package rafttest

import "gunavanthvarma-dev/raft"

type raftTest struct {
	NumberOfNodes uint64
	NodesList     []raft.RaftNode
	NodeState     map[raft.RaftState]NodeStateVariables
}

type NodeStateVariables struct {
	ElectionTimeout  uint64
	ElectionElapsed  uint64
	HeartbeatElapsed uint64
	HeartbeatTimeout uint64
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
