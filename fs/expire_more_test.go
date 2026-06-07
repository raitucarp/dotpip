package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpireMore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_expire_more_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("mykey")
	_, _ = dotfs.Set(key, "val")

	// PExpire
	res, err := dotfs.PExpire(key, 1000)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	// ExpireAt
	res, err = dotfs.ExpireAt(key, int(time.Now().Add(time.Second).Unix()))
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	// PExpireAt
	res, err = dotfs.PExpireAt(key, int(time.Now().Add(time.Second).UnixMilli()))
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	// ExpireTime
	resTime, err := dotfs.ExpireTime(key)
	assert.NoError(t, err)
	assert.Greater(t, resTime, int64(0))

	// PExpireTime
	resPTime, err := dotfs.PExpireTime(key)
	assert.NoError(t, err)
	assert.Greater(t, resPTime, int64(0))

	// PTTL
	resPTTL, err := dotfs.PTTL(key)
	assert.NoError(t, err)
	assert.Greater(t, resPTTL, int64(0))

	// Persist
	resPersist, err := dotfs.Persist(key)
	assert.NoError(t, err)
	assert.Equal(t, true, resPersist)

	// Expire on non-existent key
	res, err = dotfs.Expire(dotpip.NewKey("nonexistent"), 1)
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	// PExpire on non-existent key
	res, err = dotfs.PExpire(dotpip.NewKey("nonexistent"), 1)
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	// ExpireAt on non-existent key
	res, err = dotfs.ExpireAt(dotpip.NewKey("nonexistent"), 1)
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	// PExpireAt on non-existent key
	res, err = dotfs.PExpireAt(dotpip.NewKey("nonexistent"), 1)
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	// ExpireTime on non-existent key
	resT2, err := dotfs.ExpireTime(dotpip.NewKey("nonexistent"))
	assert.NoError(t, err)
	assert.Equal(t, int64(-2), resT2)

	// PExpireTime on non-existent key
	resT3, err := dotfs.PExpireTime(dotpip.NewKey("nonexistent"))
	assert.NoError(t, err)
	assert.Equal(t, int64(-2), resT3)

	// PTTL on non-existent key
	resT4, err := dotfs.PTTL(dotpip.NewKey("nonexistent"))
	assert.NoError(t, err)
	assert.Equal(t, int64(-2), resT4)

	// Persist on non-existent key
	res, err = dotfs.Persist(dotpip.NewKey("nonexistent"))
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	// ExpireTime on key without expiration
	keyNoExp := dotpip.NewKey("noexp")
	_, _ = dotfs.Set(keyNoExp, "val")
	resT5, err := dotfs.ExpireTime(keyNoExp)
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), resT5)

	// PExpireTime on key without expiration
	resT6, err := dotfs.PExpireTime(keyNoExp)
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), resT6)

	// PTTL on key without expiration
	resT7, err := dotfs.PTTL(keyNoExp)
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), resT7)
}

func TestPExpireAtOptions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_expire_opts_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("optkey")

	// Missing key XX
	res, err := dotfs.PExpireAt(key, int(time.Now().Add(100*time.Hour).UnixMilli()), dotpip.WithExpireXX())
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	_, _ = dotfs.Set(key, "val")

	// Present key NX
	res, err = dotfs.PExpireAt(key, int(time.Now().Add(100*time.Hour).UnixMilli()), dotpip.WithExpireNX())
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	// Has expire NX (should fail)
	res, err = dotfs.PExpireAt(key, int(time.Now().Add(100*time.Hour).UnixMilli()), dotpip.WithExpireNX())
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	// Has expire XX (should succeed)
	res, err = dotfs.PExpireAt(key, int(time.Now().Add(99*time.Hour).UnixMilli()), dotpip.WithExpireXX())
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	// GT: Only if new expiry is > current expiry
	// Current is 99 hours from now
	res, err = dotfs.PExpireAt(key, int(time.Now().Add(98*time.Hour).UnixMilli()), dotpip.WithExpireGT())
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	res, err = dotfs.PExpireAt(key, int(time.Now().Add(100*time.Hour).UnixMilli()), dotpip.WithExpireGT())
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	// LT: Only if new expiry is < current expiry
	// Current is 100 hours from now
	res, err = dotfs.PExpireAt(key, int(time.Now().Add(101*time.Hour).UnixMilli()), dotpip.WithExpireLT())
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	res, err = dotfs.PExpireAt(key, int(time.Now().Add(99*time.Hour).UnixMilli()), dotpip.WithExpireLT())
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestPExpireAtPast(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotpip_expire_past_test_")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dotfs := fs.NewFileSystem(tmpDir)
	defer dotfs.Close()

	key := dotpip.NewKey("pastkey")
	_, _ = dotfs.Set(key, "val")

	// NX when has TTL
	_, _ = dotfs.PExpire(key, 1000000)                            // set far future ttl
	res, err := dotfs.PExpireAt(key, 1000, dotpip.WithExpireNX()) // past TTL
	assert.NoError(t, err)
	assert.Equal(t, false, res) // Failed to update because NX

	// XX when has no TTL
	key2 := dotpip.NewKey("pastkey2")
	_, _ = dotfs.Set(key2, "val")
	res, err = dotfs.PExpireAt(key2, 1000, dotpip.WithExpireXX()) // past TTL
	assert.NoError(t, err)
	assert.Equal(t, false, res) // Failed to update because XX

	// GT when new TTL <= current TTL
	_, _ = dotfs.Set(key, "val")
	_, _ = dotfs.PExpireAt(key, int(time.Now().Add(10*time.Second).UnixMilli())) // future TTL
	res, err = dotfs.PExpireAt(key, 1000, dotpip.WithExpireGT())                 // past TTL
	assert.NoError(t, err)
	assert.Equal(t, false, res)

	res, err = dotfs.PExpireAt(key, 1000)
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	// Re-verify the key deletion properly
	resExists, _ := dotfs.Exists(key)
	assert.Equal(t, []bool{false}, resExists)
}
