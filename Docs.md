May 19, 2026:

- Defined main RaftNode structure accroding to the Raft paper
- Defined the 4 Messages that interact with the Raft algorithm
- The current pattern being followed is to isolate the Raft logic into a separate library and
  the actual networking, storage layer will be implemented by the server. This is done for easy testability and also the flexibility to swap betweeen different network or storage services/infra.


May 20,2026:

- Because RAFT does not rely on wall clock, I can decide on the progression of time. Initially, The frequency of the Clock ticks is going to be 100ms for simplicity, so 1 Tick = 100ms. So heartbeats are going to be one in 100ms, and election timeouts will be a multiple of Tick. in Section 9.3 of the RAFT paper, the authors recommend a timeout duration of 150-300ms so that's going to be 1-3 Ticks. Quite a low number, so we are going to have a timeout duration of 500-1500ms. So the Ticks ranges from 5 to 15 Ticks. The election timeout for each node will be a random number from 5-15 Ticks. 

- This unit of measurement is temporary, after I build a working version of the RAFT protocol I will perform experiments to find the optimal unit of frequency and adjust the election timeout accordingly.


June 6,2026:-

- Implemented Timeout generator, sendRequestVote()
- TODO next: implement processClientRequest()


June 7, 2026:

 - So for AppendEntriesRPC the leader sends it and continues processing other messages; when it gets AppendEntriesResponse how does it know for which logEntry it received ack for? 

 - maybe based on NodeId it got the ack from? like check the nextIndex; cause the leader does not send the next log entry until it got for the previous one

 - Continue thinking about how does the leader know that majority is replicated


June 13, 2026:

 - So we can add multiple log entries in a single RPC, need to create a variable to adjust batch size

 - we don't have to track the index, because we won't send the next batch until we get an ack

 - FUTURE OPTIMIZATION: SEND MULTIPLE BATCHES AT ONCE TO A SINGLE FOLLOWER --- For this we need the last index of the batch in the follower ack

 - FUTURE OPTIMIZATION: add a variable to track the latest entry received by the follower so leader could update itself




June 24, 2026:

 - Raft leader election prototype implementation is done. I need to setup a test harness to allow me to create different scenarios to test the correctness of my implementation

 