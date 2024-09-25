package snowflake

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake is a distributed unique ID generator inspired by Twitter Snowflake.
// epochStart is the timestamp when the Snowflake epoch starts (2010-11-04 01:42:54 UTC).
// workerIDBits is the number of bits for worker ID. (5 bits)
// datacenterBits is the number of bits for datacenter ID. (5 bits)
// sequenceBits is the number of bits for sequence number. (12 bits)
const (
	epochStart     = int64(1288834974657) // Twitter epoch (2010-11-04 01:42:54 UTC)
	workerIDBits   = uint(5)              // 5 bits for worker ID
	datacenterBits = uint(5)              // 5 bits for datacenter ID
	sequenceBits   = uint(12)             // 12 bits for sequence number

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits)   // -1 ^ (-1 << 5) = 31
	maxDatacenterID = int64(-1) ^ (int64(-1) << datacenterBits) // -1 ^ (-1 << 5) = 31
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)   // -1 ^ (-1 << 12) = 4095

	timeShift       = workerIDBits + datacenterBits + sequenceBits // 5 + 5 + 12 = 22
	workerShift     = datacenterBits + sequenceBits                // 5 + 12 = 17
	datacenterShift = sequenceBits                                 // 12
)

type Snowflake struct {
	mutex         sync.Mutex // Mutex to protect data
	lastTimestamp int64
	workerID      int64
	datacenterID  int64
	sequence      int64
}

func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, fmt.Errorf("datacenter ID must be between 0 and %d", maxDatacenterID)
	}

	return &Snowflake{
		workerID:      workerID,
		datacenterID:  datacenterID,
		lastTimestamp: -1,
		sequence:      0,
	}, nil
}

func (s *Snowflake) NextID() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	if timestamp < s.lastTimestamp {
		return 0 // Handle clock moving backwards
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	return ((timestamp - epochStart) << timeShift) |
		(s.datacenterID << datacenterShift) |
		(s.workerID << workerShift) |
		s.sequence
}
