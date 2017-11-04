package shapes

import (
	"io"

	"github.com/pkg/errors"

	"github.com/johnny-morrice/shapes/asm"
)

func InterpretProgramAST(ast *asm.AST, input io.Reader, output io.Writer) error {
	const errMsg = "ExecuteProgramCode failed"

	process, err := Compile(ast, StdLib())

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	builder := &RuntimeBuilder{
		Process:   process,
		Input:     input,
		Output:    output,
		Functions: StdLib().Functions,
	}

	runtime := builder.Build()

	err = runtime.Execute()

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}
