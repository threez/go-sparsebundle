package partition

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threez/go-sparsebundle"
)

func TestEmpty(t *testing.T) {
	b, err := sparsebundle.Open("../tests/empty.sparsebundle", 16)
	assert.NoError(t, err)
	defer b.Close()

	_, err = Partitions(b)
	assert.Error(t, err)
}

func TestTestFile(t *testing.T) {
	b, err := sparsebundle.Open("../tests/test.sparsebundle", 16)
	assert.NoError(t, err)
	defer b.Close()

	_, err = Partitions(b)
	assert.Error(t, err)
}

func TestFat10File(t *testing.T) {
	b, err := sparsebundle.Open("../tests/FAT10.sparsebundle", 16)
	assert.NoError(t, err)
	defer b.Close()

	tab, err := Partitions(b)
	assert.NoError(t, err)
	log.Println(tab)
}
