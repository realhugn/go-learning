package uidgenerator

import (
	"fmt"
	"sync"
	"time"
)

const (
	nodeBits     = 10                        // 10 bits for Node ID (0-1023 nodes)
	sequenceBits = 12                        // 12 bits for Sequence number (0-4095 IDs per millisecond)
	maxNodeID    = -1 ^ (-1 << nodeBits)     // Max Node ID (1023)
	maxSequence  = -1 ^ (-1 << sequenceBits) // Max Sequence (4095)

	// Shift for the timestamp, node, and sequence
	nodeShift      = sequenceBits
	timestampShift = nodeBits + sequenceBits

	// Custom epoch (Unix time in milliseconds)
	epoch = 1609459200000 // January 1, 2021 in milliseconds
)

// IDGenerator struct to hold state
type IDGenerator struct {
	mutex     sync.Mutex
	timestamp int64
	nodeID    int64
	sequence  int64
}

func NewIDGenerator(nodeID int64) *IDGenerator {
	if nodeID < 0 || nodeID > maxNodeID {
		panic(fmt.Sprintf("Node ID must be between 0 and %d", maxNodeID))
	}
	return &IDGenerator{
		nodeID: nodeID,
	}
}

func (g *IDGenerator) GenerateID() int64 {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	now := time.Now().UnixNano() / int64(time.Millisecond)

	// If the current timestamp is the same as the last one, increment the sequence
	if g.timestamp == now {
		g.sequence = (g.sequence + 1) & maxSequence
		// If sequence overflows, wait for the next millisecond
		if g.sequence == 0 {
			for now <= g.timestamp {
				now = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		// Reset sequence for new millisecond
		g.sequence = 0
	}

	g.timestamp = now

	// Generate ID by shifting and combining timestamp, node ID, and sequence
	id := ((now - epoch) << timestampShift) | (g.nodeID << nodeShift) | g.sequence
	return id
}
