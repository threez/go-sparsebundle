# go-sparsebundle [![Go Reference](https://pkg.go.dev/badge/github.com/threez/go-sparsebundle.svg)](https://pkg.go.dev/github.com/threez/go-sparsebundle) 

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

There are a few test files

* `tests/empty.sparsebundle` an almost empty file with only 10 bytes written
* `tests/FAT10.sparsebundle` a 10 MB GPT FAT file
* `tests/test.sparsebundle` an encrypted sparse bundle with a single file in it
