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

//call raftNode.Ready()
//Persist to Disk
// call raftNode.Advance()
// call raftNode.Ready()
// Send messages to other Nodes
//call raftNode.Advance()
// call raftNode.Read()
// Apply Changes to State machine
// call raftNode.Advance()
