package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestSetsCommands(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}

	for _, encoding := range encodings {
		t.Run(string(encoding), func(t *testing.T) {
			pathRoot := filepath.Join(os.TempDir(), "dotpip_sets_test_"+string(encoding))
			_ = os.MkdirAll(pathRoot, os.ModePerm)
			defer func() { _ = os.RemoveAll(pathRoot) }()

			db := fs.NewFileSystem(pathRoot)
			db.EncodeType(encoding)

			key1 := dotpip.NewKey("set1")
			key2 := dotpip.NewKey("set2")
			key3 := dotpip.NewKey("set3")

			// SAdd and SCard
			added, err := db.SAdd(key1, "a", "b", "c")
			if err != nil {
				t.Fatalf("SAdd failed: %v", err)
			}
			if added != 3 {
				t.Errorf("Expected 3 added, got %d", added)
			}

			// Add duplicate
			added, _ = db.SAdd(key1, "a", "d")
			if added != 1 {
				t.Errorf("Expected 1 added, got %d", added)
			}

			card, _ := db.SCard(key1)
			if card != 4 {
				t.Errorf("Expected card 4, got %d", card)
			}

			// SIsMember
			isMember, _ := db.SIsMember(key1, "a")
			if !isMember {
				t.Errorf("Expected 'a' to be member")
			}
			isMember, _ = db.SIsMember(key1, "z")
			if isMember {
				t.Errorf("Expected 'z' not to be member")
			}

			// SMembers
			members, _ := db.SMembers(key1)
			if len(members) != 4 {
				t.Errorf("Expected 4 members, got %d", len(members))
			}

			// SMIsMember
			smIsMember, _ := db.SMIsMember(key1, "a", "z", "b")
			if len(smIsMember) != 3 || !smIsMember[0] || smIsMember[1] || !smIsMember[2] {
				t.Errorf("SMIsMember failed: %v", smIsMember)
			}

			// SRem
			removed, _ := db.SRem(key1, "c", "d", "z")
			if removed != 2 {
				t.Errorf("Expected 2 removed, got %d", removed)
			}
			card, _ = db.SCard(key1)
			if card != 2 {
				t.Errorf("Expected card 2, got %d", card)
			}

			// SMove
			_, _ = db.SAdd(key2, "x", "y")
			moved, _ := db.SMove(key1, key2, "a")
			if !moved {
				t.Errorf("SMove failed")
			}
			isMember, _ = db.SIsMember(key1, "a")
			if isMember {
				t.Errorf("Expected 'a' removed from key1")
			}
			isMember, _ = db.SIsMember(key2, "a")
			if !isMember {
				t.Errorf("Expected 'a' added to key2")
			}

			// SDiff
			_, _ = db.SAdd(key1, "1", "2", "3")
			_, _ = db.SAdd(key2, "2", "3", "4")
			diff, _ := db.SDiff(key1, key2)
			if len(diff) != 2 { // '1', 'b' (b was left from previous ops)
				t.Errorf("Expected diff length 2, got %d: %v", len(diff), diff)
			}

			// SDiffStore
			diffCount, _ := db.SDiffStore(key3, key1, key2)
			if diffCount != len(diff) {
				t.Errorf("Expected diffStore %d, got %d", len(diff), diffCount)
			}

			// SInter
			inter, _ := db.SInter(key1, key2)
			if len(inter) != 2 { // '2', '3'
				t.Errorf("Expected inter length 2, got %d", len(inter))
			}

			// SInterCard
			interCard, _ := db.SInterCard(1, key1, key2)
			if interCard != 1 {
				t.Errorf("Expected interCard 1, got %d", interCard)
			}

			// SInterStore
			interCount, _ := db.SInterStore(key3, key1, key2)
			if interCount != 2 {
				t.Errorf("Expected interStore 2, got %d", interCount)
			}

			// SUnion
			union, _ := db.SUnion(key1, key2)
			// key1: 'b', '1', '2', '3'
			// key2: 'x', 'y', 'a', '2', '3', '4'
			// union: 'b', '1', '2', '3', 'x', 'y', 'a', '4' (8 items)
			if len(union) != 8 {
				t.Errorf("Expected union length 8, got %d: %v", len(union), union)
			}

			// SUnionStore
			unionCount, _ := db.SUnionStore(key3, key1, key2)
			if unionCount != 8 {
				t.Errorf("Expected unionStore 8, got %d", unionCount)
			}

			// SPop
			popped, _ := db.SPop(key3, 2)
			if len(popped) != 2 {
				t.Errorf("Expected 2 popped, got %d", len(popped))
			}
			card, _ = db.SCard(key3)
			if card != 6 {
				t.Errorf("Expected card 6 after pop, got %d", card)
			}

			// SRandMember
			randMembers, _ := db.SRandMember(key3, 3)
			if len(randMembers) != 3 {
				t.Errorf("Expected 3 randMembers, got %d", len(randMembers))
			}
			card, _ = db.SCard(key3)
			if card != 6 {
				t.Errorf("Expected card 6 after randMember, got %d", card)
			}
		})
	}
}
