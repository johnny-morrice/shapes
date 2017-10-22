package shapes

import (
	"testing"

	"github.com/johnny-morrice/shapes/asm"
)

func TestCompile(t *testing.T) {
	successCases := []compilation{
		addRegisterCompilation(),
		subRegisterCompilation(),
		pushRegisterCompilation(),
		popCompilation(),
		readCompilation(),
		writeCompilation(),
		setRegisterCompilation(),
		loopCompilation(),
		severalStatementsCompilation(),
	}

	for i, test := range successCases {
		ok := compileHelper(t, test.ast, MakeProcess(test.expected))

		if !ok {
			t.Errorf("Test failure at success case %d", i)
		}
	}

	failureCases := []*asm.AST{
		unknownRegisterFailure(),
		unknownStackFailure(),
		programTooBigFailure(),
	}

	for i, test := range failureCases {
		ok := compileFailureHelper(t, test)

		if !ok {
			t.Errorf("Test failure at failure case %d", i)
		}
	}
}

func addRegisterCompilation() compilation {
	statements := []asm.Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func subRegisterCompilation() compilation {
	statements := []asm.Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func pushRegisterCompilation() compilation {
	statements := []asm.Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func popCompilation() compilation {
	statements := []asm.Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func readCompilation() compilation {
	statements := []asm.Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func writeCompilation() compilation {
	statements := []asm.Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func setRegisterCompilation() compilation {
	statements := []asm.Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func copyRegisterCompilation() compilation {
	statements := []asm.Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func loopCompilation() compilation {
	statements := []asm.Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func severalStatementsCompilation() compilation {
	statements := []asm.Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func unknownRegisterFailure() *asm.AST {
	panic("not implemented")
}
func unknownStackFailure() *asm.AST {
	panic("not implemented")
}
func programTooBigFailure() *asm.AST {
	panic("not implemented")
}

func compileHelper(t *testing.T, ast *asm.AST, expected *Process) bool {
	t.Helper()

	actual, err := Compile(ast)

	if err != nil {
		t.Errorf("Expected compile success but got error: %s", err.Error())
		return false
	}

	if !expected.IsSameByteCode(actual) {
		t.Error("Expected same bytecode")
		return false
	}

	return true
}

func compileFailureHelper(t *testing.T, ast *asm.AST) bool {
	t.Helper()

	_, err := Compile(ast)

	if err == nil {
		t.Error("Expected compile failure")
		return false
	}

	return true
}

type compilation struct {
	ast      *asm.AST
	expected []Operation
}

func makeCompilation(statements []asm.Statement, expected []Operation) compilation {
	return compilation{
		ast:      &asm.AST{Statements: statements},
		expected: expected,
	}
}
