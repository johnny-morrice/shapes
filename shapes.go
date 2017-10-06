package shapes

import (
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
	OpCode  OpCode
	Operand Operand
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
	Register []byte
	Stack    []byte
	Error    error
}

func (process *Process) ExecuteStep(callTable []RuntimeCall) bool {
	if process.PC >= len(process.ByteCode) {
		return false
	}

	op := process.ByteCode[process.PC]

	impl := callTable[op.OpCode]

	impl(op.Operand)

	if process.Error != nil {
		return false
	}

	return true
}

type RuntimeCall func(operand Operand)

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

func (runtime *Runtime) jmpnz(operand Operand) {
	panic("notimplemented")
}

func (runtime *Runtime) add(operand Operand) {
	panic("notimplemented")
}

func (runtime *Runtime) sub(operand Operand) {
	panic("notimplemented")
}

func (runtime *Runtime) push(operand Operand) {
	panic("notimplemented")
}

func (runtime *Runtime) pop(operand Operand) {
	panic("notimplemented")
}

func (runtime *Runtime) read(operand Operand) {
	panic("notimplemented")
}

func (runtime *Runtime) write(operand Operand) {
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
