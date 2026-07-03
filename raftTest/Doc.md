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