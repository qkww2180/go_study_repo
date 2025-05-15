package redis_class

import "testing"

func TestLockRace(t *testing.T) {
	LockRace(client)
}

// go test -v ./g_database/redis -run=^TestLockRace$ -count=1
