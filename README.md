# go-sparsebundle

Direct access to the sparse bundle format using golang using cgo.

This library is a wrapper on top of the already existing

https://github.com/gyf304/sparsebundle-fuse due to that the software is GPLv2.


## How to use

```go
bbundle, err := Open("./tests/test.sparsebundle", 16)
if err != nil {
    panic(err)
}
defer b.Close()

// bundle implements:
// * io.ReadWriteSeeker
// * io.ReaderAt
// * io.WriterAt
// * io.Closer
```

## Test

There are two test files

* `tests/empty.sparsebundle` an empty sparse bundle
* `tests/test.sparsebundle` an encrypted sparse bundle with a single file in it
