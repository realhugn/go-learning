package snowflake

import (
	"testing"
	"time"
)

func TestNewSnowflake(t *testing.T) {
	tests := []struct {
		name         string
		workerID     int64
		datacenterID int64
		wantErr      bool
	}{
		{"Valid IDs", 1, 1, false},
		{"Max Valid IDs", maxWorkerID, maxDatacenterID, false},
		{"Invalid Worker ID", maxWorkerID + 1, 1, true},
		{"Invalid Datacenter ID", 1, maxDatacenterID + 1, true},
		{"Negative Worker ID", -1, 1, true},
		{"Negative Datacenter ID", 1, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf, err := NewSnowflake(tt.workerID, tt.datacenterID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSnowflake() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && sf == nil {
				t.Errorf("NewSnowflake() returned nil Snowflake for valid input")
			}
		})
	}
}

func TestSnowflake_NextID(t *testing.T) {
	sf, err := NewSnowflake(1, 1)
	if err != nil {
		t.Fatalf("Failed to create Snowflake: %v", err)
	}

	ids := make(map[int64]bool)
	for i := 0; i < 1000000; i++ {
		id := sf.NextID()
		if id == 0 {
			t.Errorf("NextID() returned 0")
		}
		if ids[id] {
			t.Errorf("NextID() returned duplicate ID: %d", id)
		}
		ids[id] = true
	}
}

func TestExtractInfoFromID(t *testing.T) {
	sf, _ := NewSnowflake(1, 1)
	id := sf.NextID()
	info := ExtractInfoFromID(id)

	if info.DatacenterID != 1 {
		t.Errorf("ExtractInfoFromID() datacenterID = %d, want 1", info.DatacenterID)
	}
	if info.WorkerID != 1 {
		t.Errorf("ExtractInfoFromID() workerID = %d, want 1", info.WorkerID)
	}
	if time.Since(info.Timestamp) > time.Second {
		t.Errorf("ExtractInfoFromID() timestamp too old: %v", info.Timestamp)
	}
}

func TestSnowflake_Concurrency(t *testing.T) {
	sf, _ := NewSnowflake(1, 1)
	concurrency := 10
	idsPerRoutine := 10000

	ids := make(chan int64, concurrency*idsPerRoutine)

	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < idsPerRoutine; j++ {
				ids <- sf.NextID()
			}
		}()
	}

	uniqueIDs := make(map[int64]bool)
	for i := 0; i < concurrency*idsPerRoutine; i++ {
		id := <-ids
		if uniqueIDs[id] {
			t.Errorf("Duplicate ID generated: %d", id)
		}
		uniqueIDs[id] = true
	}
}

func BenchmarkSnowflake_NextID(b *testing.B) {
	sf, _ := NewSnowflake(1, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sf.NextID()
	}
}
