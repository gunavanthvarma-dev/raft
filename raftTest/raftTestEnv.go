package raftTest

type RaftTestEnv struct{
	NumberOfNodes uint64
	TestNodesList []TestNode
	TickInterval  time.Duration
	TestNetwork	  Network
	TestStorage	  Storage
    CrashNodes	 Map[uint64][]TestNode // key - Tick interval, list of nodeIds to crash
	RestartNodes Map[uint64][]TestNode //key - Tick interval, list of nodeIds to restart 
}

// type NodeHarness struct{
// 	NumberOfNodes uint64
// 	TestNodesList []TestNode


// }

type TestNode struct{
	raftNode RaftNode
	Inbound  chan Message
	Outbound chan Message
}

type Network struct{
    //filter -- set of criteria that checks for network delay 
	// Inbound Queue map
	// Outbound queue map
}

type Storage struct{
   // log storage for all the nodes; could be a HashMap 
}

//Test Config file
	//numberofnodes
	//TickInterval
	//network filter criteria:
		//Network delay 
			// at what Tick
			// for which NodeId
			// Inbound or/And Outbound?
			// Map<Tick,NodeLinks> --- {FromNode,ToNode,Inbound/Outbound}
		//Network partition --- for now lets implement the simple model; messages are just dropped until the partition heals
			// PartitionNodesMap --- Map<Tick,NodeLinks> --- NodeLinks - {FromNode,ToNode(s)}
			// PartitionHealNodesMap --- Map<Tick,NodeLinks> --- NodeLinks - {FromNode,ToNode(s)}
	//crash nodes
		// at what tick
		// which node(s)
	//Restart nodes
		// at what tick
		// which node(s)
	//Election Timeout
		// fixed or randomi
	// Possible future improvements:
		// monitor --- collects the logs of particular nodes at particular tick