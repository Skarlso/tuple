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
	key    uint64

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
	Key() uint64
	Range(done <-chan struct{}) chan any
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

// ToSlice creates a slice out of the tuple values.
func (t *Tuple) ToSlice() []any {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := append([]any{}, t.values...)

	return result
}

// Len returns how many elements there are in the Tuple.
func (t *Tuple) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return len(t.values)
}

// Slice returns a new tuple with the syntax [from: to].
// It will panic in case of index out of bounds.
func (t *Tuple) Slice(from, to int) *Tuple {
	t.mu.RLock()
	defer t.mu.RUnlock()

	dst := make([]any, len(t.values))
	copy(dst, t.values)
	slice := dst[from:to]
	newTuple := &Tuple{values: slice}

	return newTuple
}

// Key returns a unique value for a given Tuple that can be used as a Key.
// The value is cached on the Tuple after the first call of this function
// and cannot be changed again.
// _Note_: This is a trivial approach to hashing and the usage of %v depends
// on the output of %v. If that changes, the Key changes as well.
// Further, there are cases where this might produce a similar key or a
// different key even though values didn't change. It is good for trivial
// types, but can cause some problems for more complex one.
func (t *Tuple) Key() uint64 {
	if t.key != 0 {
		return t.key
	}

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

	t.key = h.Sum64()

	return t.key
}

// Range provides a channel from which to fetch values of a tuple.
func (t *Tuple) Range(done <-chan struct{}) chan any {
	result := make(chan any)

	t.mu.Lock()
	values := make([]any, len(t.values))
	copy(values, t.values)
	t.mu.Unlock()

	go func() {
		defer close(result)
		for _, v := range values {
			select {
			case <-done:
				return
			case result <- v:
			}
		}
	}()

	return result
}
