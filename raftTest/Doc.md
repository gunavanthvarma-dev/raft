June 25, 2026

--> Need to setup a Test harness to deterministically simulate the RAFT scenarios
--> What do we need?
    1. A configuration file
        1.1 Number of nodes with node ids
        1.2 Set Initial state of the nodes
        1.3 Set up Ticker
        1.4 Setup a Transport Queue
        1.5 Setup a Persistent Storage interface that simulates a disk for all nodes. Should decide on to share same storage or create a goroutine for each node
        1.5 What all scenarios can we set up?
            1.5.1 Given a set of initial state, check the outcome after a specified duration
            1.5.2 Given a set of initial state, define messages and/Or order them based on Tick in the Transport queue, check the outcome after a specific duration or at each cycle/Tick
            1.5.3 Update Transport/Storage interface after each Tick
            1.5.4 Check output after each Tick or at the End of the scenario


# Test Infra Goals:

1. Assert the values of ServerTasks of a particular node/multiple nodes at a particular/multiple Tick(s)
2. Assert the state of a particular/multiple nodes at a particular/multiple Tick(s)
3. Provide a mechanism to set/modify state of a particular/multiple Node(s) at a particular Tick(s)

## Transport Layer:

1. Setup a Network Inbound/Outbound queue for every Node
2. Provide a mechanism to add a filter in the transport queue of a particular/multiple Node(s), filter parameters need to be determined
3. Provide a mechanism to add a delay in the transport queue of a particular/multiple Node(s) to simulate network partition/network congestion/node processing delay

## Storage Layer:

1. Setup a local storage for Node(s) 
2. Simulate the functionalities of the storage layer( persist entries and state)


## Future Goals:

1. Provide a mechanism to setup a config file to create pre-defined test scenarios; parameters to be determined later
2. Provide a mechanism to modify test parameters interactively 



Create a base TestInfraEnvironment
Parameters:
    1. Number of Nodes
    2. List of Nodes
    3. Tick Interval
    4. Standard network delay
    5. Abnormal network delay
    6. Standard process delay 
    7. Abnormal process delay
    8. Standard packet drop rate  --- Need to investigate further, accounting for TCP
    9. Abnormal packet drop rate  --- Need to investigate further, accounting for TCP
    10. Nodes that crash and do not restart
    11. Nodes that crash and restart
    12. Abnormal nodes list --- Each node has to specify the abnormal network delay, abnormal process delay, abnormal packet drop rate, crash rate, crash and restart rate, also specify the Tick interval where it occurs
    13. Transportlayer --- central switch ; has list of all inbound and outbound queues of all nodes 
    14. Storagelayer --- List of all persistent storage of all nodes

    Network delay can be set as inbound link delay or outbound link delay
    process delay can be set as a sleep()
    packet drop rate 

TestNode:
Parameters:
    1. RaftNode
    2. Tick interval
    3. Abnormal/Normal
    4. Abnormal network delay, Tick[]
    5. Abnormal process delay, Tick[]
    6. Abnormal packet drop rate, Tick[]
    7. Inbound Queue
    8. Outbound Queue
    9. Storage Layer; log

Test Env:

func init(TestEnvConfig):
    // create TestNodes based on TestEnvConfig file with an option of default/custom Node config --- need to look into it

    // Create Ticker based on Tick interval

    //  


Tests:-

1. Reach ElectionTimeout and Trigger Election, node changes status to Candidate and sends request vote to peers -- Done
2. Node receives valid requestVote, checks if its valid and sends response True -- Done
3. Node receives invalid requestVote, checks if its valid and sends response False  -- Done
4. Node receives all valid requestVote response, gets majority
5. Node receives a duplicate vote from a peer; check if its idempotent 
6. Node receives non majority valid requestVote response, check if it does not have a majority, it still stays in Candidate 
7. Node receives some invalid requestVote response and moslty valid requestVote, check if it got majority   
8. Node receives mostly invalid requestVote response, check if it does not have a majority, it still stays in Candidate
9. Node in Candidate mode, hits electionTimeout if it does not get enough votes in time
10. Node receives RequestVoteResponse with a different term; check if it rejects it
11. Node receives RequestVoteResponse after election completion; check if it handles it
12. Node elected Leader; check if it sends initial appendEntries to all the peers; check matchIndex and nextIndex initialization
11. 
