package redis_class

import (
	"context"
	"testing"
)

func TestPubSub(t *testing.T) {
	pubSub(context.Background(), client)
}

// go test -v ./database/redis -run=^TestPubSub$ -count=1
