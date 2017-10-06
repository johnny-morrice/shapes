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

	code, err := Compile(ast)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	runtime := &Runtime{
		ByteCode: code,
		Input:    input,
		Output:   output,
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

type ByteCode struct {
}

func Compile(ast *AST) ([]ByteCode, error) {
	panic("not implemented")
}

type Runtime struct {
	ByteCode []ByteCode
	Input    io.Reader
	Output   io.Writer
}

func (runtime *Runtime) Execute() error {
	panic("not implemented")
}
