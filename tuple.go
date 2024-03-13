package tuple

import (
	"bytes"
	"encoding/binary"
	"hash/maphash"
	"reflect"
	"sync"
)

// Tuple contains the values within a tuple.
type Tuple struct {
	values []any

	mu   sync.RWMutex
	hash maphash.Hash
}

// New will create a new Tuple. The Tuple can be access by its API.
// Doing so is thread safe. You can embed Tuples into a Tuple.
// t := New(1, 2, 3, New(4, 5, 6))
// t2 := t.Value(3)[*Tuple] -> a Tuple.
func New(values ...any) *Tuple {
	return &Tuple{
		values: values,
		hash:   maphash.Hash{},
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
	var sum bytes.Buffer

	for _, v := range t.values {
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:

		case reflect.Invalid:
			panic("invalid type")
		case reflect.Bool:
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			fallthrough
		case reflect.Complex64:
			fallthrough
		case reflect.Complex128:
			fallthrough
		case reflect.Uint64:
			switch ty := v.(type) {
			case int: // must be fixed size
				if err := binary.Write(&sum, binary.BigEndian, int64(ty)); err != nil {
					panic(err)
				}
			case uint: // must be fixed size
				if err := binary.Write(&sum, binary.BigEndian, uint64(ty)); err != nil {
					panic(err)
				}
			case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128:
				if err := binary.Write(&sum, binary.BigEndian, ty); err != nil {
					panic(err)
				}
			}
		case reflect.Uintptr:
			panic("unsupported type")
		case reflect.Chan:
			panic("unsupported type")
		case reflect.Func:
			panic("unsupported type")
		case reflect.Interface:
		case reflect.Map:
		case reflect.Pointer:
		case reflect.String:
			if err := binary.Write(&sum, binary.BigEndian, []byte(v.(string))); err != nil {
				panic(err)
			}
		case reflect.Struct:
			if err := binary.Write(&sum, binary.BigEndian, v); err != nil {
				panic(err)
			}
		case reflect.UnsafePointer:
			panic("unsupported type")
		}
	}

	var hash maphash.Hash
	hash.SetSeed(t.hash.Seed())
	if _, err := hash.Write(sum.Bytes()); err != nil {
		panic(err)
	}

	return hash.Sum64()

	//data := binary.BigEndian.Uint64(sum.Bytes())
	//var data uint64
	//binary.Read(sum, binary.BigEndian, &data)
	//return data
}
