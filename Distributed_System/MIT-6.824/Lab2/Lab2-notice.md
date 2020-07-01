# Lab2 - Raft

## Lab 2A

实现领导人选举和心跳包（ 没有 *log entry* 的 `AppendEntries` ）

- `AppendEntries RPC`：By accepting the RPC, the follower is implicitly telling the leader that their log matches the leader’s log up to and including the `prevLogIndex` included in the `AppendEntries` arguments.

  Upon receiving the reply, the leader might then decide that some entry has been replicated to a majority of servers, and start committing it.

- If the follower has all the entries the leader sent, the follower **MUST NOT** truncate its log. Any elements *following* the entries sent by the leader **MUST** be kept. This is because we could be receiving an outdated `AppendEntries` RPC from the leader, and truncating the log would mean “taking back” entries that we may have already told the leader that we have in our log.

- Specifically, you should *only* restart your election timer if a) you get an `AppendEntries` RPC from the *current* leader (i.e., if the term in the `AppendEntries` arguments is outdated, you should *not* reset your timer); b) you are starting an election; or c) you *grant* a vote to another peer.

- Ensure that you follow the second rule in “Rules for Servers” *before* handling an incoming RPC. The second rule states:

  > If RPC request or response contains term `T > currentTerm`: set `currentTerm = T`, convert to follower 

- if you have already voted in the current term, and an incoming `RequestVote` RPC has a higher term that you, you should *first* step down and adopt their term (thereby resetting `votedFor`), and *then* handle the RPC, which will result in you granting the vote !

- If you get an `AppendEntries` RPC with a `prevLogIndex` that points beyond the end of your log, you should handle it the same as if you did have that entry but the term did not match (i.e., reply false).

- It is important to implement the “up-to-date log” check *exactly* as described in section 5.4. No cheating and just checking the length!

- check for `commitIndex > lastApplied`  after `commitIndex` is updated (i.e., after `matchIndex` is updated)

- A leader is not allowed to update `commitIndex` to somewhere in a *previous* term (or, for that matter, a future term). Thus, as the rule says, you specifically need to check that `log[N].term == currentTerm`.

- If a leader sends out an `AppendEntries` RPC, and it is rejected, but *not because of log inconsistency* (this can only happen if our term has passed), then you should immediately step down, and *not* update `nextIndex`. If you do, you could race with the resetting of `nextIndex` if you are re-elected immediately.

- first record the term in the reply (it may be higher than your current term), and then to compare the current term with the term you sent in your original RPC. If the two are different, drop the reply and return.

- the correct thing to do is update `matchIndex` to be `prevLogIndex + len(entries[])` from the arguments you sent in the RPC originally.

- The leader has to be careful when processing replies; it must check that the term hasn't changed since sending the RPC, and must account for the possibility that replies from concurrent RPCs to the same follower have changed the leader's state (e.g. nextIndex).

**snapshot**

- If, when the server comes back up, it reads the updated snapshot, but the outdated log, it may end up applying some log entries *that are already contained within the snapshot*. This happens since the `commitIndex` and `lastApplied` are not persisted, and so Raft doesn’t know that those log entries have already been applied. The fix for this is to introduce a piece of persistent state to Raft that records what “real” index the first entry in Raft’s persisted log corresponds to. This can then be compared to the loaded snapshot’s `lastIncludedIndex` to determine what elements at the head of the log to discard.
