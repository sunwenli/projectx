package core

import (
	"fmt"
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
	data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x0b}
	vm := NewVm(data)
	assert.Nil(t, vm.Run())

	assert.Equal(t, int(3), vm.stack.Pop())

	data = []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d}
	vm = NewVm(data)
	assert.Nil(t, vm.Run())

	result := vm.stack.Pop().([]byte)
	fmt.Println("result :", string(result))
	assert.Equal(t, "FOO", string(result))
}

func TestSub(t *testing.T) {
	data := []byte{0x03, 0x0a, 0x02, 0x0a, 0x0e}
	vm := NewVm(data)
	assert.Nil(t, vm.Run())

	assert.Equal(t, int(1), vm.stack.Pop().(int))
}
