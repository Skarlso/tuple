package tuple

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tuple := New(1, 2, 3, "string")
	v := Value[int](tuple, 2)
	assert.Equal(t, 3, v)

	v2 := Value[string](tuple, 3)
	assert.Equal(t, "string", v2)
}

func TestValue(t *testing.T) {
	tuple := New(1, 2, 3, "string", New("1", "2", "3"))

	assert.Equal(t, 2, Value[int](tuple, 1))
	assert.Equal(t, &Tuple{values: []any{"1", "2", "3"}}, Value[*Tuple](tuple, 4))
}

func TestLen(t *testing.T) {
	tuple := New(1, 2, 3, "string", New("1", "2", "3"))

	assert.Equal(t, 5, tuple.Len())
}

func TestToSlice(t *testing.T) {
	tuple := New(1, 2, 3, "string")

	assert.Equal(t, []any{1, 2, 3, "string"}, tuple.ToSlice())
}

func TestConcat(t *testing.T) {
	tuple1 := New(1, 2, 3, "string")
	tuple2 := New(4, 5, 6)

	result := Concat(tuple1, tuple2)

	assert.Equal(t, &Tuple{values: []any{1, 2, 3, "string", 4, 5, 6}}, result)
}

func TestSlice(t *testing.T) {
	tuple1 := New(1, 2, 3, "string")

	result := tuple1.Slice(1, 3)

	assert.Equal(t, &Tuple{values: []any{2, 3}}, result)
	// make sure we didn't change the original
	assert.Equal(t, []any{1, 2, 3, "string"}, tuple1.values)
}

func TestHash(t *testing.T) {
	tuple1 := New(1, 2)
	tuple2 := New(1, 2)

	first := tuple1.Key()
	second := tuple1.Key()
	assert.Equal(t, first, second)
	assert.Equal(t, tuple1.Key(), tuple2.Key())

	m := make(map[uint64]struct{})
	t1 := New(1, 2, 3, 4, New("5", "6", "7", "8"))
	t2 := New(1, 2, 3, 4, New("10", "11", "12", "13"))
	m[t1.Key()] = struct{}{}
	m[t2.Key()] = struct{}{}

	fmt.Println(m)
}

func TestHashWithWeirdValues(t *testing.T) {
	tuple1 := New(1, 2, "string", &Tuple{values: []any{3, 4}})
	tuple2 := New(1, 2, "string", &Tuple{values: []any{3, 4}})
	tuple3 := New(1, 2, "string", &Tuple{values: []any{5, 6}})

	first := tuple1.Key()
	second := tuple1.Key()
	assert.Equal(t, first, second)
	assert.Equal(t, tuple1.Key(), tuple2.Key())
	assert.NotEqual(t, tuple1.Key(), tuple3.Key())
}

func TestRange(t *testing.T) {
	tuple1 := New(1, 2, "string")
	done := make(chan struct{})

	var result []any
	for v := range tuple1.Range(done) {
		result = append(result, v)
	}

	assert.Equal(t, []any{1, 2, "string"}, result)
}
