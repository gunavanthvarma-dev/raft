package raftTest

// type RaftTestEnv struct{
// 	NumberOfNodes uint64
// 	TestNodesList []TestNode
// 	TestTime      *Time
// 	TestNetwork	  *Network
// 	TestStorage	  *Storage
//     TestFailures  *Failures
// }

// func RaftTestEnvInit() *RaftTestEnv{
//     // NumberOfNodes
// 	// create time struct
// 	// create network struct
// 	// create failures struct
// 	// func to create TestNodes --- default config or custom config
//        // within this create Inbound and outbound channel
// }

// type TestNode struct{
// 	Id       NodeId
// 	raftNode RaftNode
// 	Inbound  chan Message
// 	Outbound chan Message
// 	NodeStorage []DiskEntry
// 	NodeTimeoutGen TimeoutGenerator
// }

// type Time struct{
// 	TickInterval time.Duration
// 	NodeTimeoutGen  []*TimeoutGenerator
// 	NodeTickers     []*time.Ticker
// 	FixedTimeout  bool  // true --- fixed, false --- random timeout

// }

// type Network struct{
//     //filter -- set of criteria that checks for network delay
// 	//NetworkPartition map[uint64]map[uint64]map[uint64]bool    // Map<Tick,Map<FromNode, Set of ToNodes>>
// 	// Inbound Queue map
// 	NodeInbound map[uint64]chan Message
// 	// Outbound queue map
// 	NodeOutbound map[uint64]chan Message
// }

// type Failures struct{
// 	NetworkPartition map[uint64]map[uint64]map[uint64]bool    // Map<Tick,Map<FromNode, Set of ToNodes>>

// 	CrashNodes map[uint64][]uint64   // Map<Tick,List of Nodes to crash>

// 	RestartNodes map[uint64][]uint64 // Map<Tick, List of Nodes to restart>
// }

// type Storage struct{
//    // log storage for all the nodes; could be a HashMap
//    	Disk map[uint64]*DiskEntry

// }

// //Test Config file
// 	//numberofnodes
// 	//TickInterval
// 	//network filter criteria:
// 		//Network delay
// 			// at what Tick
// 			// for which NodeId
// 			// Inbound or/And Outbound?
// 			// Map<Tick,NodeLinks> --- {FromNode,ToNode,Inbound/Outbound}
// 		//Network partition --- for now lets implement the simple model; messages are just dropped until the partition heals
// 			// PartitionNodesMap --- Map<Tick,NodeLinks> --- NodeLinks - {FromNode,ToNode(s)}
// 			// PartitionHealNodesMap --- Map<Tick,NodeLinks> --- NodeLinks - {FromNode,ToNode(s)}
// 	//crash nodes
// 		// at what tick
// 		// which node(s)
// 	//Restart nodes
// 		// at what tick
// 		// which node(s)
// 	//Election Timeout
// 		// fixed or randomi
// 	// Possible future improvements:
// 		// monitor --- collects the logs of particular nodes at particular tick
