package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunwenli/projectx/util"
)

func TestTXMaxLength(t *testing.T) {
	p := NewTxPool(1)
	p.Add(util.NewRandomTransaction(10))
	assert.Equal(t, 1, p.all.Count())
}
