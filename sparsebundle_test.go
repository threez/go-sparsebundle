package sparsebundle

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenFailed(t *testing.T) {
	_, err := Open("foo", 16)
	assert.Error(t, err)
}

func TestEmpty(t *testing.T) {
	b, err := Open("./tests/empty.sparsebundle", 16)
	assert.NoError(t, err)
	defer b.Close()

	// size
	assert.Equal(t, int64(134217728), b.Size())

	// read
	assert.Equal(t, int64(0), b.offset)
	data := make([]byte, 10)
	n, err := b.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, n, 10)
	assert.Equal(t, int64(10), b.offset)
	assert.Equal(t, []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x10}, data)

	// seek
	s, err := b.Seek(0, io.SeekStart)
	assert.Equal(t, int64(0), b.offset)
	assert.Equal(t, int64(0), s)
	assert.NoError(t, err)

	// write
	n, err = b.Write([]byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x10})
	assert.NoError(t, err)
	assert.Equal(t, n, 10)
	assert.Equal(t, int64(10), b.offset)

	// flush
	err = b.Flush()
	assert.NoError(t, err)

	// read again
	_, err = b.Seek(0, io.SeekStart)
	assert.NoError(t, err)
	n, err = b.Read(data[:])
	assert.Equal(t, n, 10)
	assert.Equal(t, []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x10}, data)
}

func TestTestFile(t *testing.T) {
	b, err := Open("./tests/test.sparsebundle", 16)
	assert.NoError(t, err)
	defer b.Close()

	// size
	assert.Equal(t, int64(100020736), b.Size())

	// read
	assert.Equal(t, int64(0), b.offset)
	data := make([]byte, 10)
	n, err := b.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, n, 10)
	assert.Equal(t, int64(10), b.offset)
	assert.Equal(t, []byte{0xf7, 0xcb, 0x34, 0x58, 0x41, 0x6b, 0x3e, 0x28, 0x40, 0x76}, data)

	// seek
	_, err = b.Seek(0, io.SeekStart)
	assert.NoError(t, err)
}

func TestFat10File(t *testing.T) {
	b, err := Open("./tests/FAT10.sparsebundle", 16)
	assert.NoError(t, err)
	defer b.Close()

	// seek
	_, err = b.Seek(440, io.SeekStart)
	assert.NoError(t, err)

	// read
	assert.Equal(t, int64(440), b.offset)
	data := make([]byte, 10)
	n, err := b.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, n, 10)
	assert.Equal(t, int64(450), b.offset)
	assert.Equal(t, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xfe, 0xff, 0xff}, data)

}
