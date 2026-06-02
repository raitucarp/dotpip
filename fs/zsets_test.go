package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestZSetsCommands(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}

	for _, encoding := range encodings {
		t.Run(string(encoding), func(t *testing.T) {
			pathRoot := filepath.Join(os.TempDir(), "dotpip_zsets_test_"+string(encoding))
			_ = os.MkdirAll(pathRoot, os.ModePerm)
			defer os.RemoveAll(pathRoot)

			db := fs.NewFileSystem(pathRoot)
			db.EncodeType(encoding)

			key1 := dotpip.NewKey("zset1")
			key2 := dotpip.NewKey("zset2")

			// ZAdd
			added, err := db.ZAdd(key1, []dotpip.Z{
				{Score: 1, Member: "one"},
				{Score: 2, Member: "two"},
				{Score: 3, Member: "three"},
			})
			if err != nil {
				t.Fatalf("ZAdd failed: %v", err)
			}
			if added != 3 {
				t.Errorf("Expected 3 added, got %d", added)
			}

			// ZCard
			card, _ := db.ZCard(key1)
			if card != 3 {
				t.Errorf("Expected card 3, got %d", card)
			}

			// ZScore
			score, _ := db.ZScore(key1, "two")
			if score != 2 {
				t.Errorf("Expected score 2, got %v", score)
			}

			// ZRank
			rank, _ := db.ZRank(key1, "two")
			if rank != 1 {
				t.Errorf("Expected rank 1, got %v", rank)
			}

			// ZRevRank
			revRank, _ := db.ZRevRank(key1, "two")
			if revRank != 1 {
				t.Errorf("Expected revRank 1, got %v", revRank)
			}

			// ZIncrBy
			newScore, _ := db.ZIncrBy(key1, 2, "one")
			if newScore != 3 {
				t.Errorf("Expected newScore 3, got %v", newScore)
			}

			// ZCount
			count, _ := db.ZCount(key1, 2, 3)
			if count != 3 {
				t.Errorf("Expected count 3, got %v", count)
			}

			// ZRem
			removed, _ := db.ZRem(key1, "two")
			if removed != 1 {
				t.Errorf("Expected 1 removed, got %v", removed)
			}
			card, _ = db.ZCard(key1)
			if card != 2 {
				t.Errorf("Expected card 2, got %d", card)
			}

			// ZPopMax
			poppedMax, _ := db.ZPopMax(key1, 1)
			if len(poppedMax) != 1 || poppedMax[0].Member != "three" {
				t.Errorf("Expected poppedMax 'three', got %v", poppedMax)
			}

			// ZPopMin
			poppedMin, _ := db.ZPopMin(key1, 1)
			if len(poppedMin) != 1 || poppedMin[0].Member != "one" {
				t.Errorf("Expected poppedMin 'one', got %v", poppedMin)
			}

			card, _ = db.ZCard(key1)
			if card != 0 {
				t.Errorf("Expected card 0, got %d", card)
			}

			// ZInter / ZUnion
			_, _ = db.ZAdd(key1, []dotpip.Z{
				{Score: 1, Member: "a"},
				{Score: 2, Member: "b"},
			})
			_, _ = db.ZAdd(key2, []dotpip.Z{
				{Score: 2, Member: "b"},
				{Score: 3, Member: "c"},
			})

			inter, _ := db.ZInter(key1, key2)
			if len(inter) != 1 || inter[0] != "b" {
				t.Errorf("Expected inter 'b', got %v", inter)
			}

			union, _ := db.ZUnionWithScores(key1, key2)
			if len(union) != 3 {
				t.Errorf("Expected union len 3, got %v", union)
			}
			for _, u := range union {
				if u.Member == "b" && u.Score != 4 {
					t.Errorf("Expected union score for 'b' to be 4, got %v", u.Score)
				}
			}
		})
	}
}

func TestZRangeCommands(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}

	for _, encoding := range encodings {
		t.Run(string(encoding), func(t *testing.T) {
			pathRoot := filepath.Join(os.TempDir(), "dotpip_zrange_test_"+string(encoding))
			_ = os.MkdirAll(pathRoot, os.ModePerm)
			defer os.RemoveAll(pathRoot)

			db := fs.NewFileSystem(pathRoot)
			db.EncodeType(encoding)

			key := dotpip.NewKey("zset_range")

			_, _ = db.ZAdd(key, []dotpip.Z{
				{Score: 1, Member: "a"},
				{Score: 2, Member: "b"},
				{Score: 3, Member: "c"},
				{Score: 4, Member: "d"},
				{Score: 5, Member: "e"},
			})

			// Test ZRange by index
			res, _ := db.ZRange(key, "0", "2")
			if len(res) != 3 || res[0] != "a" || res[2] != "c" {
				t.Errorf("Expected [a, b, c], got %v", res)
			}

			// Test ZRange by negative index
			res, _ = db.ZRange(key, "-2", "-1")
			if len(res) != 2 || res[0] != "d" || res[1] != "e" {
				t.Errorf("Expected [d, e], got %v", res)
			}

			// Test ZRange by Score
			res, _ = db.ZRange(key, "2", "4", dotpip.WithZRangeByScore())
			if len(res) != 3 || res[0] != "b" || res[2] != "d" {
				t.Errorf("Expected [b, c, d], got %v", res)
			}

			// Test ZRange by Lex
			res, _ = db.ZRange(key, "b", "d", dotpip.WithZRangeByLex())
			if len(res) != 3 || res[0] != "b" || res[2] != "d" {
				t.Errorf("Expected [b, c, d], got %v", res)
			}

			// Test ZRange Rev
			res, _ = db.ZRange(key, "0", "2", dotpip.WithZRangeRev())
			if len(res) != 3 || res[0] != "e" || res[2] != "c" {
				t.Errorf("Expected [e, d, c], got %v", res)
			}

			// Test ZRandMember
			res, _ = db.ZRandMember(key, 2)
			if len(res) != 2 {
				t.Errorf("Expected 2 rand members, got %v", len(res))
			}

			// Test ZLexCount
			count, _ := db.ZLexCount(key, "b", "d")
			if count != 3 {
				t.Errorf("Expected lex count 3, got %v", count)
			}
		})
	}
}
