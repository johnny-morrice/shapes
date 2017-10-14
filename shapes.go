package shapes

import (
	"io"

	"github.com/pkg/errors"
)

func InterpretProgram(source []byte, input io.Reader, output io.Writer) error {
	const errMsg = "ExecuteProgramCode failed"

	ast, err := Parse(source)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	process, err := Compile(ast, defaultNameTable())

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

func defaultNameTable() NameTable {
	panic("not implemented")
}
