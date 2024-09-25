package snowflake

import "time"

type SnowflakeInfo struct {
	Timestamp    time.Time
	DatacenterID int64
	WorkerID     int64
	Sequence     int64
}

func ExtractInfoFromID(id int64) SnowflakeInfo {
	timestamp := (id >> timeShift) + epochStart
	datacenterID := (id >> datacenterShift) & maxDatacenterID
	workerID := (id >> workerShift) & maxWorkerID
	sequence := id & maxSequence

	return SnowflakeInfo{
		Timestamp:    time.Unix(0, timestamp*int64(time.Millisecond)),
		DatacenterID: datacenterID,
		WorkerID:     workerID,
		Sequence:     sequence,
	}
}
