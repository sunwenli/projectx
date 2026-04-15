package core

import "fmt"

type Instruction byte

const (
	InstrPush Instruction = 0x0a
	InstrAdd  Instruction = 0x0b
)

type VM struct {
	data  []byte
	ip    int //point of instruction
	stack []byte
	sp    int // point of stack
}

func NewVm(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: make([]byte, 10),
		sp:    -1,
	}
}

func (vm *VM) Run() error {
	for {
		instr := Instruction(vm.data[vm.ip])
		if err := vm.Exec(instr); err != nil {
			return err
		}
		vm.ip++
		if vm.ip > len(vm.data)-1 {
			break
		}
		fmt.Println(instr)
	}
	return nil
}

func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrPush:
		vm.pushstack(vm.data[vm.ip-1])
	case InstrAdd:
		a := vm.stack[0]
		b := vm.stack[1]
		c := a + b
		vm.pushstack(c)
	}
	return nil
}

func (vm *VM) pushstack(b byte) {
	vm.sp++
	vm.stack[vm.sp] = b
}
