// sparsebundle implements a reader writer seeker closer
// for macOS sparse bundle format.
package sparsebundle

// #cgo CFLAGS: -D_FILE_OFFSET_BITS=64 -g -Wall
// #cgo LDFLAGS: -lpthread
// #include <stdlib.h>
// #include <sparsebundle.h>
import "C"

import (
	"errors"
	"fmt"
	"io"
	"unsafe"
)

// ensure to implement io interfaces
var _ io.ReadWriteSeeker = (*Bundle)(nil)
var _ io.Closer = (*Bundle)(nil)
var _ io.ReaderAt = (*Bundle)(nil)
var _ io.WriterAt = (*Bundle)(nil)

var (
	ErrWhence       = errors.New("Seek: invalid whence")
	ErrOffset       = errors.New("Seek: invalid offset")
	ErrMaxOpenBands = errors.New("maxOpenBands have to be at least 10")
	ErrNotOpen      = errors.New("bundle not opened")
)

// Implements the
type Bundle struct {
	handle C.sparse_handle_t
	path   *C.char
	offset int64
}

func Open(path string, maxOpenBands int) (*Bundle, error) {
	var bundle Bundle
	var opts C.struct_sparse_options

	opts.path = C.CString(path)
	bundle.path = opts.path
	if maxOpenBands < 10 {
		return nil, ErrMaxOpenBands
	}
	opts.max_open_bands = C.int(maxOpenBands)
	e := C.sparse_open(&bundle.handle, &opts)
	if e != 0 {
		return nil, fmt.Errorf("failed to open %q: %w", path, bundle.Err())
	}

	return &bundle, nil
}

// Seek as in io.Seeker
func (b *Bundle) Seek(offset int64, whence int) (int64, error) {
	if b.handle == nil {
		return 0, ErrNotOpen
	}
	switch whence {
	case io.SeekStart:
		// leaf offset as is
	case io.SeekCurrent:
		offset += b.offset
	case io.SeekEnd:
		offset += int64(b.Size())
	default:
		return 0, ErrWhence
	}

	if offset < 0 {
		return 0, ErrOffset
	}

	b.offset = offset

	return offset, nil

}

// Read as in io.Reader
func (b *Bundle) Read(p []byte) (n int, err error) {
	n, err = b.ReadAt(p, b.offset)
	if err != nil {
		return
	}
	b.offset += int64(n)
	return
}

// ReadAt as in io.ReaderAt
func (b *Bundle) ReadAt(p []byte, off int64) (n int, err error) {
	if b.handle == nil {
		return 0, ErrNotOpen
	}
	ptr := (*C.char)(unsafe.Pointer(&p[0]))
	s := C.size_t(cap(p))
	offset := C.off_t(off)

	e := C.sparse_pread(b.handle, ptr, s, offset)
	if e < 0 {
		return 0, fmt.Errorf("failed to read: %w", b.Err())
	}

	return int(e), nil
}

// Write as in io.Writer
func (b *Bundle) Write(p []byte) (n int, err error) {
	n, err = b.WriteAt(p, b.offset)
	if err != nil {
		return
	}
	b.offset += int64(n)
	return
}

// WriteAt as in io.WriterAt
func (b *Bundle) WriteAt(p []byte, off int64) (n int, err error) {
	if b.handle == nil {
		return 0, ErrNotOpen
	}
	ptr := (*C.char)(unsafe.Pointer(&p[0]))
	s := C.size_t(len(p))
	offset := C.off_t(b.offset)

	e := C.sparse_pwrite(b.handle, ptr, s, offset)
	if e < 0 {
		return 0, fmt.Errorf("failed to write: %w", b.Err())
	}

	return int(e), nil
}

// Trim will remove size data at seeked position
func (b *Bundle) Trim(size int) (n int, err error) {
	if b.handle == nil {
		return 0, ErrNotOpen
	}
	return 0, nil
}

// Close closes the bundle and frees the resources
func (b *Bundle) Close() error {
	if b.handle == nil {
		return ErrNotOpen
	}
	e := C.sparse_close(&b.handle)
	if e != 0 {
		return fmt.Errorf("failed to close: %w", b.Err())
	}
	C.free(unsafe.Pointer(b.path))
	return nil
}

// Size returns the byte size for the bundle
func (b *Bundle) Size() int64 {
	if b.handle == nil {
		return -1
	}
	return int64(C.sparse_get_size(b.handle))
}

// Flush flushes everything to disk
func (b *Bundle) Flush() error {
	if b.handle == nil {
		return ErrNotOpen
	}
	e := C.sparse_flush(b.handle)
	if e != 0 {
		return b.Err()
	}
	return nil
}

func (b *Bundle) Err() error {
	cstr := C.sparse_get_error(b.handle)
	return errors.New(C.GoString(cstr))
}
