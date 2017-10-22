package shapes

import (
	"bytes"
	"strconv"
	"testing"
)

func TestRuntimeExecute(t *testing.T) {
	addInput, addExpect := addTestData()
	subInput, subExpect := subTestData()
	pushInput, pushExpect := pushTestData()
	popInput, popExpect := popTestData()
	readInput, readExpect := readTestData()
	writeInput, writeExpect := writeTestData()
	copyInput, copyExpect := copyTestData()
	setInput, setExpect := setTestData()
	jmpnzInput, jmpnzExpect := jmpnzTestData()

	table := [][2]cannedProcess{
		[2]cannedProcess{addInput, addExpect},
		[2]cannedProcess{subInput, subExpect},
		[2]cannedProcess{pushInput, pushExpect},
		[2]cannedProcess{popInput, popExpect},
		[2]cannedProcess{readInput, readExpect},
		[2]cannedProcess{writeInput, writeExpect},
		[2]cannedProcess{copyInput, copyExpect},
		[2]cannedProcess{setInput, setExpect},
		[2]cannedProcess{jmpnzInput, jmpnzExpect},
	}

	for i, test := range table {
		ok := runtimeExecuteHelper(t, test[0], test[1])
		if !ok {
			t.Errorf("Failure in test case %d", i)
		}
	}
}

func jmpnzTestData() (cannedProcess, cannedProcess) {
	byteCode := []Operation{
		Operation{OpCode: OP_SET, Operand: [2]Operand{0, 10}},
		Operation{OpCode: OP_SET, Operand: [2]Operand{1, 1}},
		Operation{OpCode: OP_ADD, Operand: [2]Operand{2, 0}},
		Operation{OpCode: OP_SUB, Operand: [2]Operand{0, 1}},
		Operation{OpCode: OP_JMPNZ, Operand: [2]Operand{0, 2}},
	}

	input := makeInputProcess(byteCode, []byte{})

	register := [REGISTER_COUNT]byte{}
	register[1] = 1
	register[2] = 55
	expect := makeExpectProcess(byteCode, register, [REGISTER_COUNT][]byte{}, []byte{})

	return input, expect
}

func addTestData() (cannedProcess, cannedProcess) {
	byteCode := []Operation{
		Operation{OpCode: OP_SET, Operand: [2]Operand{0, 10}},
		Operation{OpCode: OP_SET, Operand: [2]Operand{1, 20}},
		Operation{OpCode: OP_ADD, Operand: [2]Operand{0, 1}},
	}

	input := makeInputProcess(byteCode, []byte{})

	register := [REGISTER_COUNT]byte{}
	register[0] = 30
	register[1] = 20
	expect := makeExpectProcess(byteCode, register, [REGISTER_COUNT][]byte{}, []byte{})

	return input, expect
}

func subTestData() (cannedProcess, cannedProcess) {
	byteCode := []Operation{
		Operation{OpCode: OP_SET, Operand: [2]Operand{0, 10}},
		Operation{OpCode: OP_SET, Operand: [2]Operand{1, 20}},
		Operation{OpCode: OP_SUB, Operand: [2]Operand{1, 0}},
	}

	input := makeInputProcess(byteCode, []byte{})

	register := [REGISTER_COUNT]byte{}
	register[0] = 10
	register[1] = 10
	expect := makeExpectProcess(byteCode, register, [REGISTER_COUNT][]byte{}, []byte{})

	return input, expect
}

func pushTestData() (cannedProcess, cannedProcess) {
	byteCode := []Operation{
		Operation{OpCode: OP_SET, Operand: [2]Operand{1, 20}},
		Operation{OpCode: OP_PUSH, Operand: [2]Operand{0, 1}},
	}

	input := makeInputProcess(byteCode, []byte{})

	register := [REGISTER_COUNT]byte{}
	register[1] = 20
	stack := [REGISTER_COUNT][]byte{}
	stack[0] = []byte{20}
	expect := makeExpectProcess(byteCode, register, stack, []byte{})

	return input, expect
}

func popTestData() (cannedProcess, cannedProcess) {
	byteCode := []Operation{
		Operation{OpCode: OP_SET, Operand: [2]Operand{0, 40}},
		Operation{OpCode: OP_SET, Operand: [2]Operand{1, 20}},
		Operation{OpCode: OP_PUSH, Operand: [2]Operand{0, 0}},
		Operation{OpCode: OP_PUSH, Operand: [2]Operand{0, 1}},
		Operation{OpCode: OP_POP, Operand: [2]Operand{0, 2}},
	}

	input := makeInputProcess(byteCode, []byte{})

	register := [REGISTER_COUNT]byte{}
	register[0] = 40
	register[1] = 20
	register[2] = 20
	stack := [REGISTER_COUNT][]byte{}
	stack[0] = []byte{40}
	expect := makeExpectProcess(byteCode, register, stack, []byte{})

	return input, expect
}

func readTestData() (cannedProcess, cannedProcess) {
	byteCode := []Operation{
		Operation{OpCode: OP_READ, Operand: [2]Operand{60, 0}},
	}

	input := makeInputProcess(byteCode, []byte{42})

	register := [REGISTER_COUNT]byte{}
	register[60] = 42
	expect := makeExpectProcess(byteCode, register, [REGISTER_COUNT][]byte{}, []byte{})

	return input, expect
}

func writeTestData() (cannedProcess, cannedProcess) {
	byteCode := []Operation{
		Operation{OpCode: OP_SET, Operand: [2]Operand{50, 10}},
		Operation{OpCode: OP_WRITE, Operand: [2]Operand{50, 0}},
	}

	input := makeInputProcess(byteCode, []byte{})

	register := [REGISTER_COUNT]byte{}
	register[50] = 10
	output := []byte{10}
	expect := makeExpectProcess(byteCode, register, [REGISTER_COUNT][]byte{}, output)

	return input, expect
}

func copyTestData() (cannedProcess, cannedProcess) {
	byteCode := []Operation{
		Operation{OpCode: OP_SET, Operand: [2]Operand{0, 10}},
		Operation{OpCode: OP_COPY, Operand: [2]Operand{1, 0}},
	}

	input := makeInputProcess(byteCode, []byte{})

	register := [REGISTER_COUNT]byte{}
	register[0] = 10
	register[1] = 10
	expect := makeExpectProcess(byteCode, register, [REGISTER_COUNT][]byte{}, []byte{})

	return input, expect
}

func setTestData() (cannedProcess, cannedProcess) {
	byteCode := []Operation{
		Operation{OpCode: OP_SET, Operand: [2]Operand{0, 10}},
	}

	input := makeInputProcess(byteCode, []byte{})

	register := [REGISTER_COUNT]byte{}
	register[0] = 10
	expect := makeExpectProcess(byteCode, register, [REGISTER_COUNT][]byte{}, []byte{})

	return input, expect
}

func binary(text string, size int) uint64 {
	num, err := strconv.ParseUint(text, 2, size)

	if err != nil {
		panic(err)
	}

	return num
}

type cannedProcess struct {
	process *Process
	input   *bytes.Buffer
	output  *bytes.Buffer
}

func makeInputProcess(byteCode []Operation, input []byte) cannedProcess {
	return cannedProcess{
		process: &Process{
			ByteCode: byteCode,
		},
		input:  bytes.NewBuffer(input),
		output: &bytes.Buffer{},
	}
}

func makeExpectProcess(byteCode []Operation, register [REGISTER_COUNT]byte, stack [REGISTER_COUNT][]byte, output []byte) cannedProcess {
	return cannedProcess{
		process: &Process{
			ByteCode: byteCode,
			Register: register,
			Stack:    stack,
		},
		input:  &bytes.Buffer{},
		output: bytes.NewBuffer(output),
	}
}

func (canned cannedProcess) makeRuntime() *Runtime {
	return MakeRuntime(canned.process, canned.input, canned.output)
}

func (canned cannedProcess) sameRegisters(t *testing.T, other cannedProcess) bool {
	t.Helper()

	for i, myReg := range canned.process.Register {
		theirReg := other.process.Register[i]
		if myReg != theirReg {
			t.Errorf("Mismatched register %d: expected %d but was %d", i, myReg, theirReg)
			return false
		}
	}

	return true
}
func (canned cannedProcess) sameStack(t *testing.T, other cannedProcess) bool {
	t.Helper()

	for i, myStack := range canned.process.Stack {
		theirStack := other.process.Stack[i]
		myStackLen := len(canned.process.Stack)
		theirStackLen := len(other.process.Stack)
		if myStackLen != theirStackLen {
			t.Errorf("Expected stack %d length %d but was %d", i, myStackLen, theirStackLen)
			return false
		}

		for j, myEntry := range myStack {
			theirEntry := theirStack[j]

			if myEntry != theirEntry {
				t.Errorf("Mismatched stack %d entry %d: expected %d but was %d", i, j, myEntry, theirEntry)
				return false
			}
		}

	}

	return true
}
func (canned cannedProcess) sameOutput(t *testing.T, other cannedProcess) bool {
	t.Helper()

	myOutput := canned.output.Bytes()
	theirOutput := other.output.Bytes()

	myOutputLen := len(myOutput)
	theirOutputLen := len(theirOutput)
	if myOutputLen != theirOutputLen {
		t.Errorf("Expected stack length %d but was %d", myOutputLen, theirOutputLen)
		return false
	}

	for i, myOut := range myOutput {
		theirOut := theirOutput[i]

		if myOut != theirOut {
			t.Errorf("Mismatched stack entry %d: expected %d but was %d", i, myOut, theirOut)
			return false
		}
	}

	return true
}

func runtimeExecuteHelper(t *testing.T, input cannedProcess, expect cannedProcess) bool {
	t.Helper()

	runtime := input.makeRuntime()
	err := runtime.Execute()

	ok := (err == nil) == (input.process.Error == nil) == (expect.process.Error == nil)

	if !ok {
		t.Error("Unexpected error state")
	}

	ok = ok && expect.sameRegisters(t, input)
	ok = ok && expect.sameStack(t, input)
	ok = ok && expect.sameOutput(t, input)

	return ok
}
