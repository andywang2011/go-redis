package goredis

import "testing"

func TestClientInterface(t *testing.T) {
	var _ Client = (*fakeRedis)(nil)
}
