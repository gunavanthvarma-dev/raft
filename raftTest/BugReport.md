June 28,2026

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