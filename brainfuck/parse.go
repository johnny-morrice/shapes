package brainfuck

import (
	"errors"
	"fmt"

	"github.com/johnny-morrice/shapes/asm"
)

func Parse(source []byte) (*asm.AST, error) {
	loopDepth := 0
	builder := &asm.ASTBuilder{}

	pushValue := &asm.PushStmt{}
	pushValue.Operand[0] = __STACK_INDEX
	pushValue.Operand[1] = __VALUE_REGISTER
	popValue := &asm.PopStmt{}
	popValue.Operand[0] = __STACK_INDEX
	popValue.Operand[1] = __VALUE_REGISTER
	increment := &asm.AddStmt{}
	increment.Operand[0] = __VALUE_REGISTER
	increment.Operand[1] = __INCREMENT_REGISTER
	decrement := &asm.SubStmt{}
	decrement.Operand[0] = __VALUE_REGISTER
	decrement.Operand[1] = __INCREMENT_REGISTER

	newTape := &asm.CallStmt{
		VmFunc: asm.TAPE_NEW,
	}
	newTape.Operand = __STACK_INDEX
	moveTape := &asm.CallStmt{
		VmFunc: asm.TAPE_MOVE_HEAD,
	}
	moveTape.Operand = __STACK_INDEX
	writeTape := &asm.CallStmt{
		VmFunc: asm.TAPE_WRITE_HEAD,
	}
	writeTape.Operand = __STACK_INDEX
	readTape := &asm.CallStmt{
		VmFunc: asm.TAPE_READ_HEAD,
	}
	readTape.Operand = __STACK_INDEX

	pushTapeLeft := &asm.PushStmt{}
	pushTapeLeft.Operand[0] = __STACK_INDEX
	pushTapeLeft.Operand[1] = __TAPE_LEFT_REGISTER
	pushTapeRight := &asm.PushStmt{}
	pushTapeRight.Operand[0] = __STACK_INDEX
	pushTapeRight.Operand[1] = __TAPE_RIGHT_REGISTER
	pushTapeIndex := &asm.PushStmt{}
	pushTapeIndex.Operand[0] = __STACK_INDEX
	pushTapeIndex.Operand[1] = __TAPE_INDEX_REGISTER
	popTapeIndex := &asm.PopStmt{}
	popTapeIndex.Operand[0] = __STACK_INDEX
	popTapeIndex.Operand[1] = __TAPE_INDEX_REGISTER

	output := &asm.WriteStmt{}
	output.Operand = __VALUE_REGISTER
	input := &asm.ReadStmt{}
	input.Operand = __VALUE_REGISTER

	setTapeLeft := &asm.SetStmt{}
	setTapeLeft.Operand[0] = __TAPE_LEFT_REGISTER
	setTapeLeft.Operand[1] = -1
	setTapeRight := &asm.SetStmt{}
	setTapeRight.Operand[0] = __TAPE_RIGHT_REGISTER
	setTapeRight.Operand[1] = 1
	setIncrement := &asm.SetStmt{}
	setIncrement.Operand[0] = __INCREMENT_REGISTER
	setIncrement.Operand[1] = 1
	resetValue := &asm.SetStmt{}
	resetValue.Operand[0] = __VALUE_REGISTER
	resetValue.Operand[1] = 0

	builder.Append(setTapeLeft, setTapeRight, setIncrement)
	builder.Append(newTape, popTapeIndex)

	for _, chr := range source {
		switch chr {
		case '<':
			builder.Append(pushTapeLeft, pushTapeIndex, moveTape, pushTapeIndex, readTape, popValue)
		case '>':
			builder.Append(pushTapeRight, pushTapeIndex, moveTape, pushTapeIndex, readTape, popValue)
		case '+':
			builder.Append(increment, pushValue, pushTapeIndex, writeTape)
		case '-':
			builder.Append(decrement, pushValue, pushTapeIndex, writeTape)
		case '.':
			builder.Append(output)
		case ',':
			builder.Append(input, pushValue, pushTapeIndex, writeTape)
		case '[':
			builder.OpenLoop(__VALUE_REGISTER)
			loopDepth++
		case ']':
			err := builder.LeaveBlock()
			if err != nil {
				return nil, errors.New("Closed non-existent loop")
			}
			loopDepth--
		}
	}

	if loopDepth != 0 {
		return nil, fmt.Errorf("Unexpected loop nesting depth %d", loopDepth)
	}

	return builder.AST, nil
}

const __STACK_INDEX = 0

const (
	__VALUE_REGISTER = iota
	__TAPE_INDEX_REGISTER
	__INCREMENT_REGISTER
	__TAPE_LEFT_REGISTER
	__TAPE_RIGHT_REGISTER
)
