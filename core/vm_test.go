package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	stack := NewStack(16)
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	stack.Pop()
	assert.Equal(t, 2, stack.data[0])
	stack.Pop()
	assert.Equal(t, 3, stack.data[0])
}

func TestVm(t *testing.T) {
	contractstat := NewState()
	// data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x0b}
	// vm := NewVm(data, contractstat)
	// assert.Nil(t, vm.Run())

	// assert.Equal(t, int(3), vm.stack.Pop())

	// data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d}
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	vm := NewVm(data, contractstat)
	assert.Nil(t, vm.Run())

	valuebyte, err := vm.contractState.Get([]byte("FOO"))
	assert.Nil(t, err)
	value := deserializeInt64(valuebyte)
	assert.Equal(t, int64(5), value)
}

func TestSub(t *testing.T) {
	contractstat := NewState()
	data := []byte{0x03, 0x0a, 0x02, 0x0a, 0x0e}
	vm := NewVm(data, contractstat)
	assert.Nil(t, vm.Run())

	assert.Equal(t, int(1), vm.stack.Pop().(int))
}
