May 19, 2026:

- Defined main RaftNode structure accroding to the Raft paper
- Defined the 4 Messages that interact with the Raft algorithm
- The current pattern being followed is to isolate the Raft logic into a separate library and
  the actual networking, storage layer will be implemented by the server. This is done for easy testability and also the flexibility to swap betweeen different network or storage services/infra.