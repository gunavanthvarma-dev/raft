package raft

func main() {
	// get config details from CLI

	//create Ticker
	//initialize network layer
	//initialize storage layer

	// create a raft node based on config

	//main process loop

	//select stmt
	// <-ticker
} // call raftNode.Tick()
// msg <-networkChannel
// call raftNode.ProcessNetworkMessage(msg)
// clientRequest <-clientChannel
// call raftNode.ProcessClientRequest(clientRequest)
//end select

//call raftNode.ServerTasks()
//Persist to Disk
// Send messages to other Nodes
// Apply Changes to State machine
