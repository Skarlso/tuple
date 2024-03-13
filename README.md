# tuple

A tuple "data type" for Go mimicking the API of Python tuples.

The aim was to keep the interface as simple as possible. On top of that
to provide a type safe way to access elements without having to resort
to struct types. Further, this library supports any number of items in
the tuple. After the tuple has been created, it's immutable. It cannot
be changed again.

## Examples

Get a simple value.

```go
t1 := tuple.New(1, 2, 3, 4)
v := tuple.Value[int](t1, 1)
fmt.Println(v) // 2 (int)
```

Note, that `Value` uses `index` to access elements and typed parameters to
fetch the right type. It will panic if the value is NOT of the right type.

Let's see with nested Tuple.

```go
t1 := tuple.New(1, 2, 3, 4, tuple.New("1", "2"))
v := tuple.Value[*Tuple](t1, 4)
fmt.Println(v) // &Tuple{} with values "1", "2"
```

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

Once a tuple is created it can be used as a key in a hash map. To use it as such do the following:

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

## Panics

This library _panics_ on errors. This is to keep it simple. Also, it only panics
on actual panics, such as, un-hashable types and index out of bounds errors.

## Why use this at all?

That's an interesting question. I would argue, why not? The biggest sell point here,
in my humble opinion, is that it has a hash value that can be used in a Map.
I used this for various interesting operations on [Advent Of Code](https://adventofcode.com/) problems.
