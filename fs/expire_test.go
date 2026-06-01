package fs

import (
	"dotpip"
	"testing"
	"time"
)

func TestExpire(t *testing.T) {
	dotfs.FlushAll()
	k := dotpip.NewKey("expire_test")

	dotfs.Set(k, "value")

	res, err := dotfs.Expire(k, 1)
	if err != nil || !res {
		t.Errorf("Expire should succeed on existing key")
	}

	ttl, _ := dotfs.TTL(k)
	if ttl <= 0 || ttl > 1 {
		t.Errorf("TTL should be 1 second, got %d", ttl)
	}

	time.Sleep(1200 * time.Millisecond)

	exists, _ := dotfs.Exists(k)
	if exists[0] {
		t.Errorf("Key should be expired")
	}
}

func TestPExpire(t *testing.T) {
	dotfs.FlushAll()
	k := dotpip.NewKey("pexpire_test")

	dotfs.Set(k, "value")

	res, err := dotfs.PExpire(k, 100)
	if err != nil || !res {
		t.Errorf("PExpire should succeed")
	}

	pttl, _ := dotfs.PTTL(k)
	if pttl <= 0 || pttl > 100 {
		t.Errorf("PTTL should be <= 100ms, got %d", pttl)
	}

	time.Sleep(150 * time.Millisecond)

	_, err = dotfs.Get(k)
	if err == nil {
		t.Errorf("Get should fail on expired key")
	}
}

func TestPersist(t *testing.T) {
	dotfs.FlushAll()
	k := dotpip.NewKey("persist_test")

	dotfs.Set(k, "value")
	dotfs.Expire(k, 1)

	res, err := dotfs.Persist(k)
	if err != nil || !res {
		t.Errorf("Persist should succeed on key with TTL")
	}

	ttl, _ := dotfs.TTL(k)
	if ttl != -1 {
		t.Errorf("TTL should be -1 for persisted key, got %d", ttl)
	}

	time.Sleep(1200 * time.Millisecond)

	exists, _ := dotfs.Exists(k)
	if !exists[0] {
		t.Errorf("Key should not be expired after Persist")
	}
}

func TestExpireOptions(t *testing.T) {
	dotfs.FlushAll()
	k := dotpip.NewKey("expire_options_test")
	dotfs.Set(k, "value")

	// NX on key without TTL should work
	res, _ := dotfs.Expire(k, 10, dotpip.WithExpireNX())
	if !res {
		t.Errorf("Expire NX should succeed on key without TTL")
	}

	// NX on key with TTL should fail
	res, _ = dotfs.Expire(k, 20, dotpip.WithExpireNX())
	if res {
		t.Errorf("Expire NX should fail on key with TTL")
	}

	// XX on key with TTL should work
	res, _ = dotfs.Expire(k, 20, dotpip.WithExpireXX())
	if !res {
		t.Errorf("Expire XX should succeed on key with TTL")
	}
}
