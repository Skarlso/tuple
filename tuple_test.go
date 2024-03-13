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

	first := tuple1.Sum()
	second := tuple1.Sum()
	assert.Equal(t, first, second)

	assert.Equal(t, tuple1.Sum(), tuple2.Sum())
}

func TestHashWithWeirdValues(t *testing.T) {
	tuple1 := New(1, 2, "string", &Tuple{values: []any{3, 4}})
	tuple2 := New(1, 2, "string", &Tuple{values: []any{3, 4}})
	tuple3 := New(1, 2, "string", &Tuple{values: []any{5, 6}})

	first := tuple1.Sum()
	second := tuple1.Sum()
	assert.Equal(t, first, second)
	fmt.Println(first, second)
	assert.Equal(t, tuple1.Sum(), tuple2.Sum())
	assert.NotEqual(t, tuple1.Sum(), tuple3.Sum())
}