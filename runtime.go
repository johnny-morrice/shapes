package shapes

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type Address uint8

type LongAddress uint16

type Operation struct {
	OpCode  OpCode
	Operand [2]byte
}

func (op Operation) LongAddress() LongAddress {
	long := LongAddress(op.Operand[0])
	long = long << 8
	long = long | LongAddress(op.Operand[1])
	return long
}

func (op Operation) Address(operand Address) Address {
	return Address(op.Operand[operand])
}

type OpCode byte

const (
	OP_JMPNZ = OpCode(iota)
	OP_ADD
	OP_SUB
	OP_PUSH
	OP_POP
	OP_READ
	OP_WRITE
	OP_COPY
	OP_SET
)

type Process struct {
	PC       LongAddress
	ByteCode []Operation
	Register [MAX_WORD]byte
	Stack    [MAX_WORD][]byte
	Error    error
}

func MakeProcess(byteCode []Operation) *Process {
	return &Process{
		ByteCode: byteCode,
	}
}

func Compile(ast *AST) (*Process, error) {
	compiler := &CompileVisitor{}
	ast.Visit(compiler)

	return compiler.Process, compiler.Error
}

func (process *Process) IsSameByteCode(other *Process) bool {
	panic("not implemented")
}

func (process *Process) IsTerminated() bool {
	return process.PC >= LongAddress(len(process.ByteCode)) || process.Error != nil
}

func (process *Process) ExecuteStep(callTable []RuntimeCall) {
	op := process.ByteCode[process.PC]

	impl := callTable[op.OpCode]

	impl(op)
}

func (process *Process) Peek(stackAddr Address) byte {
	stack := process.Stack[stackAddr]

	if len(stack) == 0 {
		process.failEmptyStack(stackAddr)
		return 0
	}

	return stack[len(stack)-1]
}

func (process *Process) Pop(stackAddr Address) byte {
	tip := process.Peek(stackAddr)

	if process.Error != nil {
		return 0
	}

	stack := process.Stack[stackAddr]
	process.Stack[stackAddr] = stack[:len(stack)-1]

	return tip
}

func (process *Process) failEmptyStack(stackAddr Address) {
	process.Error = fmt.Errorf("stack was empty '%d'", stackAddr)
}

func (process *Process) Push(stackAddr Address, tip byte) {
	process.Stack[stackAddr] = append(process.Stack[stackAddr], tip)
}

func (process *Process) GetRegister(register Address) byte {
	return process.Register[register]
}

func (process *Process) SetRegister(register Address, val byte) {
	process.Register[register] = val
}

func (process *Process) IncrementPC() {
	process.PC++
}

type RuntimeCall func(op Operation)

type Runtime struct {
	Process     *Process
	Error       error
	CallTable   []RuntimeCall
	Input       io.Reader
	Output      io.Writer
	readBuffer  []byte
	writeBuffer []byte
}

func MakeRuntime(process *Process, input io.Reader, output io.Writer) *Runtime {
	runtime := &Runtime{
		Process:     process,
		Input:       input,
		Output:      output,
		readBuffer:  []byte{0},
		writeBuffer: []byte{0},
	}

	runtime.CallTable = []RuntimeCall{
		runtime.jmpnz,
		runtime.add,
		runtime.sub,
		runtime.push,
		runtime.pop,
		runtime.read,
		runtime.write,
		runtime.copy,
		runtime.set,
	}

	return runtime
}

func (runtime *Runtime) hasError() bool {
	return runtime.Process.Error != nil
}

func (runtime *Runtime) jmpnz(op Operation) {
	val := runtime.Process.GetRegister(op.Address(0))

	if runtime.hasError() {
		return
	}

	if val != 0 {
		runtime.Process.PC = op.LongAddress()
	} else {
		runtime.Process.IncrementPC()
	}
}

func (runtime *Runtime) onRegisters(op Operation, f func(valZero, valOne byte) byte) {
	valZero := runtime.Process.GetRegister(op.Address(0))

	if runtime.hasError() {
		return
	}

	valOne := runtime.Process.GetRegister(op.Address(1))

	if runtime.hasError() {
		return
	}

	newVal := f(valZero, valOne)
	runtime.Process.SetRegister(op.Address(0), newVal)
}

func (runtime *Runtime) add(op Operation) {
	runtime.onRegisters(op, func(valZero, valOne byte) byte {
		return valZero + valOne
	})
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) sub(op Operation) {
	runtime.onRegisters(op, func(valZero, valOne byte) byte {
		return valZero - valOne
	})
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) push(op Operation) {
	val := runtime.Process.GetRegister(op.Address(1))

	if runtime.hasError() {
		return
	}

	runtime.Process.Push(op.Address(0), val)
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) pop(op Operation) {
	tip := runtime.Process.Pop(op.Address(0))

	if runtime.hasError() {
		return
	}

	runtime.Process.SetRegister(op.Address(1), tip)
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) read(op Operation) {
	const errMsg = "Runtime.read failed"

	_, err := runtime.Input.Read(runtime.readBuffer)

	if err != nil {
		runtime.Process.Error = errors.Wrap(err, errMsg)
		return
	}

	runtime.Process.SetRegister(op.Address(0), runtime.readBuffer[0])
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) write(op Operation) {
	const errMsg = "Runtime.Write failed"

	val := runtime.Process.GetRegister(op.Address(0))
	runtime.writeBuffer[0] = val

	_, err := runtime.Output.Write(runtime.writeBuffer)

	if err != nil {
		runtime.Process.Error = errors.Wrap(err, errMsg)
		return
	}
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) copy(op Operation) {
	val := runtime.Process.GetRegister(op.Address(1))

	if runtime.hasError() {
		return
	}

	runtime.Process.SetRegister(op.Address(0), val)
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) set(op Operation) {
	runtime.Process.SetRegister(op.Address(0), byte(op.Address(1)))
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) Execute() error {
	for !runtime.Process.IsTerminated() {
		runtime.Process.ExecuteStep(runtime.CallTable)
	}

	if runtime.Process.Error != nil {
		return runtime.Process.Error
	}

	return nil
}

const MAX_WORD = 256
