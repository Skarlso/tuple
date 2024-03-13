package tuple

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"sync"
)

// Tuple contains the values within a tuple.
type Tuple struct {
	values []any

	mu sync.RWMutex
}

// New will create a new Tuple. The Tuple can be access by its API.
// Doing so is thread safe. You can embed Tuples into a Tuple.
// t := New(1, 2, 3, New(4, 5, 6))
// t2 := t.Value(3)[*Tuple] -> a Tuple.
func New(values ...any) *Tuple {
	return &Tuple{
		values: values,
	}
}

// Contract doesn't include Value function so the value getter can be generic and
// used for accessing elements of the Tuple in a generic way.
// In order to prevent having to define a type during creation and other operations
// Value will not be on the tuple but on the Package.
type Contract interface {
	// ToSlice sadly in this case, they will have to do some type assertions.
	ToSlice() []any
	Len() int
	Slice(from, to int) *Tuple
	Sum() uint64
}

var _ Contract = &Tuple{}

// Concat adds up tuples and creates a new resulting tuple with all values in
// order as they were created.
func Concat(tuple ...*Tuple) *Tuple {
	result := &Tuple{}

	for _, v := range tuple {
		result.values = append(result.values, v.ToSlice()...)
	}

	return result
}

// Value will panic if the index is out of range. This is to keep in-line with tuple logic.
// Defining a type makes sure we get the right type when accessing values.
func Value[T any](t *Tuple, index int) (value T) {
	t.mu.RLock()
	val := t.values[index]
	t.mu.RUnlock()

	value = val.(T)

	return value
}

func (t *Tuple) ToSlice() []any {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := append([]any{}, t.values...)

	return result
}

func (t *Tuple) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return len(t.values)
}

// Slice will panic on case of index out of bounds.
func (t *Tuple) Slice(from, to int) *Tuple {
	t.mu.RLock()
	defer t.mu.RUnlock()

	dst := make([]any, len(t.values))
	copy(dst, t.values)
	slice := dst[from:to]
	newTuple := &Tuple{values: slice}

	return newTuple
}

func (t *Tuple) Sum() uint64 {
	h := fnv.New64a()
	v := reflect.ValueOf(t.values)

	// Iterate over the elements of the slice
	for i := 0; i < v.Len(); i++ {
		element := v.Index(i).Interface()
		// Hash each element
		if _, err := h.Write([]byte(fmt.Sprintf("%v", element))); err != nil {
			panic(err)
		}
	}

	return h.Sum64()
}
