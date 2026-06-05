sed -i 's/assert.NoError(t, err) \/\/ empty map/assert.Error(t, err) \/\/ empty map/g' fs/vectors_test.go
sed -i '60s/assert.NoError(t, err)/assert.Error(t, err)/' fs/vectors_test.go
sed -i '78s/assert.NoError(t, err)/assert.Error(t, err)/' fs/vectors_test.go
sed -i '112s/assert.NoError(t, err)/assert.Error(t, err)/' fs/vectors_test.go
sed -i '116s/assert.NoError(t, err)/assert.Error(t, err)/' fs/vectors_test.go
sed -i '150s/assert.NoError(t, err)/assert.Error(t, err)/' fs/vectors_test.go
sed -i '153s/assert.NoError(t, err)/assert.Error(t, err)/' fs/vectors_test.go
