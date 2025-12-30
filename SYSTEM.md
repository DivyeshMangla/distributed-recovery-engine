# Distributed Node Protocol

## Core Concepts

**Identity**: `(NodeID, LifeID)` where `LifeID` is monotonic. Higher `LifeID` always wins.

**States**: `ALIVE` → `SUSPECT` → `DEAD` → `QUARANTINED` → `ALIVE`
- Missed heartbeats: `ALIVE` → `SUSPECT`
- Heartbeat OK: `SUSPECT` → `ALIVE`
- Confirmed failure: `SUSPECT` → `DEAD`
- Higher `LifeID` hello: `DEAD` → `QUARANTINED`
- Realive accepted: `QUARANTINED` → `ALIVE`

`DEAD` is final per `LifeID`. All packets ignored if `LifeID` is lower.

---

## Packets

**Hello** - Introduce identity on startup/reconnect
- Higher `LifeID` → override state to `QUARANTINED`, reply `HelloAck(accepted=true)`
- Same `LifeID` + different address → `HelloAck(accepted=false)`

**Heartbeat** - Liveness signal (periodic in `ALIVE`/`SUSPECT`)
- Ignored if `DEAD`/`QUARANTINED`

**Gossip** - Distribute metadata (never introduces new nodes)

**Realive** - Request re-admission from `QUARANTINED`
- Reply `RealiveAck(admitted=true/false)`
- Admitted → `QUARANTINED` → `ALIVE`

---

## Timing & Parameters

`T_base` = base interval (200ms for datacenters, 1s standard, 5s for WAN)  
`N` = cluster size

**Heartbeat**: Every `T_base ± 25% jitter` to 1 random peer
- Traffic: O(N) msgs/sec cluster-wide

**Gossip**: Every `T_base ± 25% jitter` to `⌈log₂(N)⌉ + 1` random peers
- Batch multiple updates per message
- Traffic: O(N log N) msgs/sec
- Convergence: O(log N) rounds

**Suspicion**: `ALIVE` → `SUSPECT` after `max(5 * T_base, 3 * RTT_p99)`
- Accounts for probabilistic heartbeat coverage

**Failure**: `SUSPECT` → `DEAD` after `2 * suspicion_timeout` AND gossip confirmation from `⌈log₂(N)⌉` peers

**Realive**: Send to `⌈log₂(N)⌉ + 1` random peers (3-11 typically)
- Retry every `2 * T_base ± 50% jitter`
- Max 10 attempts before backoff

**Example (N=1000, T_base=1s)**:
- Heartbeats: 1k/sec, Gossip: 10k/sec
- Suspicion: ~5s, Failure: ~10s
- Realive: 11 peers (vs 500 with N/2)

---

## Rules

- Higher `LifeID` overrides all prior state
- Silence alone never causes `DEAD` (requires gossip quorum)
- `ALIVE` status requires `Realive` acceptance
- Each packet has one responsibility
- Jitter prevents synchronized bursts