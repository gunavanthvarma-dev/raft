May 19, 2026:

- Defined main RaftNode structure accroding to the Raft paper
- Defined the 4 Messages that interact with the Raft algorithm
- The current pattern being followed is to isolate the Raft logic into a separate library and
  the actual networking, storage layer will be implemented by the server. This is done for easy testability and also the flexibility to swap betweeen different network or storage services/infra.


May 20,2026:

- Because RAFT does not rely on wall clock, I can decide on the progression of time. Initially, The frequency of the Clock ticks is going to be 100ms for simplicity, so 1 Tick = 100ms. So heartbeats are going to be one in 100ms, and election timeouts will be a multiple of Tick. in Section 9.3 of the RAFT paper, the authors recommend a timeout duration of 150-300ms so that's going to be 1-3 Ticks. Quite a low number, so we are going to have a timeout duration of 500-1500ms. So the Ticks ranges from 5 to 15 Ticks. The election timeout for each node will be a random number from 5-15 Ticks. 

- This unit of measurement is temporary, after I build a working version of the RAFT protocol I will perform experiments to find the optimal unit of frequency and adjust the election timeout accordingly.

