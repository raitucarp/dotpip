package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestVectors(t *testing.T) {
	os.RemoveAll("test_data_vectors")
	db := fs.NewFileSystem("test_data_vectors")
	defer func() {
		db.Close()
		os.RemoveAll("test_data_vectors")
	}()

	key := dotpip.NewKey("my_vector")

	// Options Tests
	opt := dotpip.WithVAddCas(true)
	opt2 := dotpip.WithVAddEF(1)
	opt3 := dotpip.WithVAddM(1)
	opt4 := dotpip.WithVAddQuant("Q8")
	opt5 := dotpip.WithVAddReduceDim(1)
	_ = opt
	_ = opt2
	_ = opt3
	_ = opt4
	_ = opt5

	o1 := dotpip.WithVSimCount(1)
	o2 := dotpip.WithVSimEF(1)
	o3 := dotpip.WithVSimEpsilon(0.5)
	o4 := dotpip.WithVSimFilter("test")
	o5 := dotpip.WithVSimFilterEF(1)
	o6 := dotpip.WithVSimNoThread(true)
	o7 := dotpip.WithVSimTruth(true)
	o8 := dotpip.WithVSimWithAttribs(true)
	o9 := dotpip.WithVSimWithScores(true)
	_ = o1
	_ = o2
	_ = o3
	_ = o4
	_ = o5
	_ = o6
	_ = o7
	_ = o8
	_ = o9

	// Test Error key not exist
	_, err := db.VCard(key)
	assert.NoError(t, err) // readVectorSet should return empty for not exist

	_, err = db.VDim(key)
	assert.Error(t, err) // empty map

	_, err = db.VInfo(key)
	assert.Error(t, err)

	_, err = db.VLinks(key, "a")
	assert.NoError(t, err)

	_, err = db.VRandMember(key, 1)
	assert.NoError(t, err)

	_, err = db.VRange(key, "a", "b")
	assert.NoError(t, err)

	// Test VAdd
	n, err := db.VAdd(key, "item1", []float32{1.0, 2.0, 3.0})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	// Try VAdd with different dim
	_, err = db.VAdd(key, "item2", []float32{1.0, 2.0})
	assert.Error(t, err)

	// VAdd exist
	n, err = db.VAdd(key, "item1", []float32{1.0, 2.0, 3.0})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	// Test VCard
	card, err := db.VCard(key)
	assert.NoError(t, err)
	assert.Equal(t, 1, card)

	// Test VDim
	dim, err := db.VDim(key)
	assert.NoError(t, err)
	assert.Equal(t, 3, dim)

	// Test VIsMember
	isMember, err := db.VIsMember(key, "item1", "item2")
	assert.NoError(t, err)
	assert.Equal(t, []bool{true, false}, isMember)

	// Test VEmb
	emb, err := db.VEmb(key, "item1", "item2")
	assert.NoError(t, err)
	assert.Equal(t, [][]float32{{1.0, 2.0, 3.0}, nil}, emb)

	// Test VSetAttr and VGetAttr
	n, err = db.VSetAttr(key, "item1", `{"color":"red"}`)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	// Error setting attr on not exist
	_, err = db.VSetAttr(key, "item2", `{"color":"red"}`)
	assert.Error(t, err)

	// Error setting attr with invalid JSON
	_, err = db.VSetAttr(key, "item1", `invalid`)
	assert.Error(t, err)

	attr, err := db.VGetAttr(key, "item1", "item2")
	assert.NoError(t, err)
	assert.Equal(t, `{"color":"red"}`, *attr[0])
	assert.Nil(t, attr[1])

	// Test VInfo
	info, err := db.VInfo(key)
	assert.NoError(t, err)
	assert.Equal(t, 3, info["dimensions"])
	assert.Equal(t, 1, info["count"])

	// Test VRandMember
	randM, err := db.VRandMember(key, 5)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(randM))

	// Test VRange
	_, _ = db.VAdd(key, "item2", []float32{1.1, 2.1, 3.1})
	rng, err := db.VRange(key, "item1", "item2")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(rng))

	// Test VSim
	sim, err := db.VSim(key, []float32{1.0, 2.0, 3.0})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(sim))

	sim, err = db.VSim(key, "item1")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(sim))

	_, err = db.VSim(key, "nonexist")
	assert.Error(t, err)

	_, err = db.VSim(key, 123)
	assert.Error(t, err)

	// Test VRem
	n, err = db.VRem(key, "item1", "item2")
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	// Should delete key
	card, err = db.VCard(key)
	assert.NoError(t, err)
	assert.Equal(t, 0, card)
}

func TestVectorsExtraCoverage(t *testing.T) {
    // cover untested lines
	os.RemoveAll("test_data_vectors_extra")
	db := fs.NewFileSystem("test_data_vectors_extra")
	defer func() {
		db.Close()
		os.RemoveAll("test_data_vectors_extra")
	}()

	key := dotpip.NewKey("my_vector2")

	// VCard empty
	c, _ := db.VCard(key)
    assert.Equal(t, 0, c)

	// VEmb non existent
	e, err := db.VEmb(key, "item1")
	assert.NoError(t, err)
	assert.Equal(t, [][]float32{nil}, e)

	_, _ = db.VAdd(key, "item1", []float32{1, 2})

	// VGetAttr non existent attributes
	a, err := db.VGetAttr(key, "item1")
	assert.NoError(t, err)
	assert.Nil(t, a[0])

	_, _ = db.VRem(key, "item1") // empty it

	// writeVectorSet empty
	_, _ = db.VRem(key, "item1")
	_, _ = db.VAdd(key, "item1", []float32{1})
	_, _ = db.VRem(key, "item1")

}

func TestVLinks(t *testing.T) {
	os.RemoveAll("test_data_vlinks")
	db := fs.NewFileSystem("test_data_vlinks")
	defer func() {
		db.Close()
		os.RemoveAll("test_data_vlinks")
	}()

	key := dotpip.NewKey("my_vector")
	res, err := db.VLinks(key, "a")
	assert.NoError(t, err)
	assert.Equal(t, make(map[int][]string), res)
}
