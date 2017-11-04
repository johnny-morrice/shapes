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
	popValue := &asm.PushStmt{}
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
	moveTape := &asm.CallStmt{
		VmFunc: asm.TAPE_MOVE_HEAD,
	}
	writeTape := &asm.CallStmt{
		VmFunc: asm.TAPE_WRITE_HEAD,
	}
	readTape := &asm.CallStmt{
		VmFunc: asm.TAPE_READ_HEAD,
	}

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
	input := &asm.WriteStmt{}
	input.Operand = __VALUE_REGISTER

	builder.Append(newTape, popTapeIndex)

	for _, chr := range source {
		switch chr {
		case '<':
			builder.Append(pushTapeLeft, pushTapeIndex, moveTape)
		case '>':
			builder.Append(pushTapeRight, pushTapeIndex, moveTape)
		case '+':
			builder.Append(readTape, popValue, increment, pushValue, pushTapeIndex, writeTape)
		case '-':
			builder.Append(readTape, popValue, decrement, pushValue, pushTapeIndex, writeTape)
		case '.':
			builder.Append(pushTapeIndex, readTape, popValue, output)
		case ',':
			builder.Append(input, pushValue, pushTapeIndex, writeTape)
		case '[':
			builder.OpenLoop(__VALUE_REGISTER)
		case ']':
			err := builder.LeaveBlock()
			if err != nil {
				return nil, errors.New("Closed non-existent loop")
			}
		}
	}

	if loopDepth > 0 {
		return nil, fmt.Errorf("Parse failed with %d unclosed loops", loopDepth)
	}

	return builder.AST, nil
}

const __STACK_INDEX = 0
const __TAPE_INDEX_REGISTER = 1
const __VALUE_REGISTER = 0
const __INCREMENT_REGISTER = 2
const __TAPE_LEFT_REGISTER = 3
const __TAPE_RIGHT_REGISTER = 4
const __RETURN_REGISTER = 5
