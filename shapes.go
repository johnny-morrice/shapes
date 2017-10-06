package shapes

import (
	"io"

	"github.com/pkg/errors"
)

func ExecuteProgramCode(source string, input io.Reader, output io.Writer) error {
	const errMsg = "ExecuteProgramCode failed"

	ast, err := Parse(source)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	process, err := Compile(ast)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	runtime := &Runtime{
		Process: process,
		Input:   input,
		Output:  output,
	}

	err = runtime.Execute()

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

type AST struct {
}

func Parse(source string) (*AST, error) {
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

func (runtime *Runtime) Execute() error {
	panic("not implemented")
}
