package shapes

import (
	"io"

	"github.com/pkg/errors"

	"github.com/johnny-morrice/shapes/asm"
)

func InterpretProgramAST(ast *asm.AST, input io.Reader, output io.Writer) error {
	const errMsg = "ExecuteProgramCode failed"

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
