package shapes

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

func ExecuteProgramCode(source []byte, input io.Reader, output io.Writer) error {
	const errMsg = "ExecuteProgramCode failed"

	ast, err := Parse(source)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	process, err := Compile(ast)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	runtime := MakeRuntime(process, input, output)

	err = runtime.Execute()

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

type AST struct {
}

func Parse(source []byte) (*AST, error) {
	panic("not implemented")
}

type Operand int

type Operation struct {
	OpCode   OpCode
	Register int
	Operand  Operand
}

type OpCode byte

const (
	JMPNZ = OpCode(iota)
	ADD
	SUB
	PUSH
	POP
	READ
	WRITE
)

func Compile(ast *AST) (*Process, error) {
	panic("not implemented")
}

type Process struct {
	PC       int
	ByteCode []Operation
	Register [REGISTER_COUNT]byte
	Stack    []byte
	Error    error
}

func (process *Process) ExecuteStep(callTable []RuntimeCall) bool {
	if process.PC >= len(process.ByteCode) {
		return false
	}

	op := process.ByteCode[process.PC]

	impl := callTable[op.OpCode]

	impl(op)

	if process.Error != nil {
		return false
	}

	return true
}

func (process *Process) Peek() byte {
	if len(process.Stack) == 0 {
		process.failEmptyStack()
		return 0
	}

	return process.Stack[len(process.Stack)-1]
}

func (process *Process) Pop() byte {
	tip := process.Peek()

	if process.Error != nil {
		return 0
	}

	process.Stack = process.Stack[:len(process.Stack)-1]

	return tip
}

func (process *Process) failEmptyStack() {
	process.Error = errors.New("stack was empty")
}

func (process *Process) failNoSuchRegister(register int) {
	process.Error = fmt.Errorf("No such register '%d'", register)
}

func (process *Process) Push(tip byte) {
	process.Stack = append(process.Stack, tip)
}

func (process *Process) GetRegister(register int) byte {
	if register >= REGISTER_COUNT {
		process.failNoSuchRegister(register)
		return 0
	}

	return process.Register[register]
}

func (process *Process) SetRegister(register int, val byte) {
	if register >= REGISTER_COUNT {
		process.failNoSuchRegister(register)
		return
	}

	process.Register[register] = val
}

type RuntimeCall func(op Operation)

type Runtime struct {
	Process   *Process
	CallTable []RuntimeCall
	Input     io.Reader
	Output    io.Writer
}

func MakeRuntime(process *Process, input io.Reader, output io.Writer) *Runtime {
	runtime := &Runtime{
		Process: process,
		Input:   input,
		Output:  output,
	}

	runtime.CallTable = []RuntimeCall{
		runtime.jmpnz,
		runtime.add,
		runtime.sub,
		runtime.push,
		runtime.pop,
		runtime.read,
		runtime.write,
	}

	return runtime
}

func (runtime *Runtime) jmpnz(op Operation) {
	val := runtime.Process.GetRegister(op.Register)

	if val != 0 {
		runtime.Process.PC = int(op.Operand)
	}
}

func (runtime *Runtime) add(op Operation) {
	panic("not implemented")
}

func (runtime *Runtime) sub(op Operation) {
	panic("notimplemented")
}

func (runtime *Runtime) push(op Operation) {
	panic("notimplemented")
}

func (runtime *Runtime) pop(op Operation) {
	panic("notimplemented")
}

func (runtime *Runtime) read(op Operation) {
	panic("notimplemented")
}

func (runtime *Runtime) write(op Operation) {
	panic("notimplemented")
}

func (runtime *Runtime) Execute() error {
	for runtime.Process.ExecuteStep(runtime.CallTable) {
	}

	if runtime.Process.Error != nil {
		return runtime.Process.Error
	}

	return nil
}

const REGISTER_COUNT = 256
