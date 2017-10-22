package shapes

import (
	"testing"
)

func TestCompile(t *testing.T) {
	successCases := []compilation{
		addRegisterCompilation(),
		addLiteralCompilation(),
		subRegisterCompilation(),
		subLiteralCompilation(),
		pushRegisterCompilation(),
		pushLiteralCompilation(),
		popCompilation(),
		readCompilation(),
		writeCompilation(),
		setRegisterCompilation(),
		setLiteralCompilation(),
		loopCompilation(),
		severalStatementsCompilation(),
	}

	for i, test := range successCases {
		ok := compileHelper(t, test.ast, MakeProcess(test.expected))

		if !ok {
			t.Errorf("Test failure at success case %d", i)
		}
	}

	failureCases := []*AST{
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
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func addLiteralCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func subRegisterCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func subLiteralCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func pushRegisterCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func pushLiteralCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func popCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func readCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func writeCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func setRegisterCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func setLiteralCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func loopCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func severalStatementsCompilation() compilation {
	statements := []Statement{}
	expected := []Operation{}
	panic("not implemented")
	return makeCompilation(statements, expected)
}

func unknownRegisterFailure() *AST {
	panic("not implemented")
}
func unknownStackFailure() *AST {
	panic("not implemented")
}
func programTooBigFailure() *AST {
	panic("not implemented")
}

func compileHelper(t *testing.T, ast *AST, expected *Process) bool {
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

func compileFailureHelper(t *testing.T, ast *AST) bool {
	t.Helper()

	_, err := Compile(ast)

	if err == nil {
		t.Error("Expected compile failure")
		return false
	}

	return true
}

type compilation struct {
	ast      *AST
	expected []Operation
}

func makeCompilation(statements []Statement, expected []Operation) compilation {
	return compilation{
		ast:      &AST{Statements: statements},
		expected: expected,
	}
}
