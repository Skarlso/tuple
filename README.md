# tuple

![logo](./logo.jpeg)

A tuple "data type" for Go mimicking the API of Python tuples.

The aim was to keep the interface as simple as possible. On top of that
to provide a type safe way to access elements without having to resort
to struct types. Further, this library supports any number of items in
the tuple. After the tuple has been created, it's immutable. It cannot
be changed again.

## Examples

### Values

```go
    t1 := tuple.New(1, 2, 3, 4)
    v := tuple.Value[int](t1, 1)
    fmt.Println(v) // 2 (int)
```

Note, that `Value` uses `index` to access elements and typed parameters to
fetch the right type. It will panic if the value is NOT of the right type.

Also _note_ that `Value` ISN'T on the Tuple! It's a package level function.
The reason for that is to keep the interface and the tuple creation type free.
If the Tuple creation had generics definitions it would mess with the
rest of the functions. Thus, it was decided to leave the Value as a package
function to have the benefit of generics.

If someone has a better ideas, PRs are welcome. ;)

Let's see Values with a nested Tuple.

```go
    t1 := tuple.New(1, 2, 3, 4, tuple.New("1", "2"))
    v := tuple.Value[*Tuple](t1, 4)
    fmt.Println(v) // &Tuple{} with values "1", "2"
```

### Range

```go
	tuple1 := New(1, 2, "string")
	done := make(chan struct{})

	var result []any
	for v := range tuple1.Range(done) {
		result = append(result, v)
	}

	fmt.Println(result) // []any{1, 2, "string"}
```

### Len

```go
    tuple1 := New(1, 2, "string")
	fmt.Println(tuple1.Len()) // 3
```

### Concatenate

You can add Tuples together.

```go
	tuple1 := New(1, 2, 3, "string")
	tuple2 := New(4, 5, 6)

	result := Concat(tuple1, tuple2)

	fmt.Println(result) // []any{1, 2, 3, "string", 4, 5, 6}
```

### Slice

Get a new Tuple for a subset of a Tuple. The original is unchanged.

```go
	tuple1 := New(1, 2, 3, "string")

	result := tuple1.Slice(1, 3)

    fmt.Println(result) // &Tuple{values: []any{2, 3}}
```

### ToSlice

Create a slice out of a Tuple.

```go
	tuple := New(1, 2, 3, "string")

	fmt.Println(tuple.ToSlice()) // []any{1, 2, 3, "string"}
```

## Concurrent safe

All tuple operations are using a Read sync mutex, thus should be safe for concurrent operations.

## Python tuple feature parity

Python tuple features that have been implemented:

- [x] to slice
- [x] immutable once created
- [x] fixed number of items
- [x] ordered -> insertion order
- [x] tuples can be stored inside tuples
- [x] any number of items
- [x] indexable
- [x] threadsafe
- [x] sliceable
- [x] combinable
- [x] can be used as hash keys
- [x] range over

## Hashing

Once a tuple is created, it can be used as a key in a hash map. To use it as such do the following:

```go
    t1 := tuple.New(1, 2, 3 ,4, tuple.New("5", "6", "7", "8"))
    t2 := tuple.New(1, 2, 3 ,4, tuple.New("5", "6", "7", "8"))
    m := make(map[uint64]struct{})
    m[t1.Key()] = struct{}{}
    m[t2.Key()] = struct{}{}
    // map[5075198087340781659:{}]

    t1 = tuple.New(1, 2, 3 ,4, tuple.New("5", "6", "7", "8"))
    t2 = tuple.New(1, 2, 3 ,4, tuple.New("10", "11", "12", "13"))
    m = make(map[uint64]struct{})
    m[t1.Key()] = struct{}{}
    m[t2.Key()] = struct{}{}
    // map[5075198087340781659:{} 15816401886143238155:{}]
```

_Note_: The hashing "algorithm" for generating the keys is super trivial. It wouldn't stand against millions of
values and could cause collisions quickly if values are only slightly different. The key generation depends on
`%v` to use as a clutch.

Further, because of the string representation, the output of `Key()` could change if the struct order
is modified. Thus, it is advised to avoid serializing the output of `Key()`.


## Panics

This library _panics_ on errors. This is to keep it simple. Also, it only panics
on actual panics, such as, un-hashable types and index out of bounds errors.

## Why use this at all?

That's an interesting question. I would argue, why not? The biggest sell point here,
in my humble opinion, is that it has a hash value that can be used in a Map.
I used this for various interesting operations on [Advent Of Code](https://adventofcode.com/) problems.
