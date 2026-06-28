package main

import (
	"gunavanthvarma-dev/raft"
	"log"
)

func TriggerElectionTimeout() {
	//raftTest := newRaftTest(3)
	//tickDuration,_ := time.ParseDuration("100ms")
	//Ticker := time.NewTicker(tickDuration)
	//var networkChannel chan raft.Message
	// for{
	// 	select{
	// 	case msg<-networkChannel:

	// 	}
	// }
	timeoutGen := raft.NewTimeoutGenerator(5, 15)
	testNode := raft.NewRaftNode(1, []raft.NodeId{2, 3, 4}, 6, 4, timeoutGen)
	for val := range 7 {
		testNode.Tick()
		log.Printf("\nTick %d", val)
		serverTasks := testNode.Ready()
		log.Printf("\nServerTasks from Ready():")
		log.Printf("\nMessages:")
		for idx, msg := range serverTasks.Messages {
			log.Printf("\n%d -> %s", idx, msg.ToString())
		}
		//Persist to Disk
		//Send messages
		// apply entries
		testNode.Advance()
	}

}

func main() {
	TriggerElectionTimeout()
}
