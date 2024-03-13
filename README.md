# tuple

A tuple "data type" for Go mimicking the API of Python tuples.

## Features

- [ ] to slice
- [ ] immutable once created
- [ ] fixed number of items
- [ ] ordered -> insertion order. this should be fine
- [ ] tuples can be stored inside tuples
- [ ] any number of items
- [ ] indexable
- [ ] threadsafe
- [ ] sliceable -> but does NOT change the original tuple ( copy then return )
- [ ] combinable -> concatenation operations
- [ ] hashable -> they can be map keys -> this one might be really tricky
- [ ] range over -> probably a function


## Hashing

Use a `gob.Encode` to create a byte stream from the struct than use https://pkg.go.dev/hash/maphash to turn it into a Sum64.
