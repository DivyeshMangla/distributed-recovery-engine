package membership

import (
	"math/rand"
	"time"

	"github.com/divyeshmangla/distributed-recovery-engine/internal/protocol"
)

func (m *Membership) PickRandomPeer(selfID protocol.NodeID) *Member {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.members) <= 1 {
		return nil
	}

	peers := make([]*Member, 0, len(m.members))
	for id, member := range m.members {
		if id == selfID || member.Status == Dead {
			continue
		}
		peers = append(peers, member)
	}

	if len(peers) == 0 {
		return nil
	}

	return peers[rand.Intn(len(peers))]
}

// SelectGossipTargets returns a hybrid selection: self + recently active + random sample
func (m *Membership) SelectGossipTargets(selfID protocol.NodeID, maxCount int, recentWindow time.Duration) []Member {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	result := make([]Member, 0, maxCount)
	included := make(map[protocol.NodeID]bool)

	if self, ok := m.members[selfID]; ok {
		result = append(result, *self)
		included[selfID] = true
	}

	m.addRecentlySeenMembers(&result, included, now, recentWindow, maxCount)
	m.addRecentlyUpdatedMembers(&result, included, now, recentWindow, maxCount)
	m.fillWithRandomPeers(&result, included, maxCount)

	return result
}

func (m *Membership) addRecentlySeenMembers(result *[]Member, included map[protocol.NodeID]bool, now time.Time, window time.Duration, maxCount int) {
	for id, member := range m.members {
		if len(*result) >= maxCount {
			return
		}
		if included[id] || member.Status == Dead {
			continue
		}
		if now.Sub(member.LastSeen) <= window {
			*result = append(*result, *member)
			included[id] = true
		}
	}
}

func (m *Membership) addRecentlyUpdatedMembers(result *[]Member, included map[protocol.NodeID]bool, now time.Time, window time.Duration, maxCount int) {
	for id, member := range m.members {
		if len(*result) >= maxCount {
			return
		}
		if included[id] || member.Status == Dead {
			continue
		}
		if now.Sub(member.UpdatedAt) <= window {
			*result = append(*result, *member)
			included[id] = true
		}
	}
}

func (m *Membership) fillWithRandomPeers(result *[]Member, included map[protocol.NodeID]bool, maxCount int) {
	if len(*result) >= maxCount {
		return
	}

	var candidates []*Member
	for id, member := range m.members {
		if !included[id] && member.Status != Dead {
			candidates = append(candidates, member)
		}
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for _, member := range candidates {
		if len(*result) >= maxCount {
			return
		}
		*result = append(*result, *member)
	}
}
