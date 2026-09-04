

1. TestName: TriggerElectionTimeout()

        Problem Found: 
            --> In sendRequestVote(), LastLogIndex is calculated as node.log[len(node.log)-1].Index, simlar for LastLogTerm
            --> Problem is when the log is empty, the len() returns 0 and then you will be accessing a slice with an index of -1, which throws a runtime panic.

        Solution:

            --> return 0 when length of the log is 0

2. TestName: TriggerElectionTimeout()

        Problem Found:
            --> when Advance() is called, it clears out the Messages field in ServerTasks, the func clear() does not clear the slice but just resets the values in it to the default value.

        Solution:

            --> Reslicing


3. TestName: TestRequestVoteResponseTrue()

        Problem Found:
            --> so &5.2 in the paper where we check the last log entry while processing the RequestVote RPC, we check the last element in the receiver node's log. what if the log is empty? like in the scenario where its the first election, the log is empty, so there will be a panic.

        Solution:
            --> include a default log entry in all nodes


4. TestName: TestRequestVoteResponseFalse()

        Problem Found:
            --> so after receiving requestVote rpc, we check the term in the message; if its less than the receiving node's current term, we send a false msg. Problem was that condition should be if/else. but I gave it just as an if condition.
                                            if msg.Term < node.currentTerm {
                                            node.RequestVoteResponseFalse(msg.FromNodeId)
                                            }
                                            switch node.NodeStatus {                     THIS SHOULD BE IN ELSE CONDITION!!!!!
        Solution:
            --> Add an else condition

5. TestName: TestCandidateReceivesAllValidRequestVoteGetsMajority()

        Problem Found: 
            --> Majority vote check condition was >=len(Peers)/2 which was wrong as Peers did not account for the the candidate node itself as part of the vote; so it appeared to have majority vote when in reality it got triggered at 50% and not >50% of the vote

        Solution:
            --> changed it to > (len(Peers)+1)/2

6. TestName: TestCandidateReceivesDuplicateRequestVoteIsIdempotent()

        Problem Found:
            --> so when the same peerNode sends a duplicate vote again; its not idempotent and the electionVotes increases again

        Solution:
            --> Make ElectionVotes variable in RaftNode struct a Set to enable idempotency

5. TestName: TestCandidateReceiveVoteFromAllNodesAndConvertToLeaderAndSendInitialAppendEntriesRPC()

        Problem Found:
            --> So in Rules for Leaders: &5.2 A newely elected leader needs to send an initial AppendEntriesRPC to all followers but currently in the code the initial RPC is not sent until the leader hits the heartbeat timeout. so the problem is any other node can trigger election again if the leader does not announce itself in time.

        Solution:
            --> as soon as the node changes its status to leader; send a heartbeat

        
        Problem Found:
            --> so nextIndex, matchIndex are empty map, so when a new leader tries to send appendEntriesRPC; it tries to access and empty map and throws panic

        Solution:
            --> initailze values with last log entry of the leader+1 upon election [mentioned in volatile state of leaders] 
                                        


