package shapes

import (
	"testing"
)

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

func TestCompile(t *testing.T) {
	addRegisterSuccess := compilation{}
	addLiteralSuccess := compilation{}
	subRegisterSuccess := compilation{}
	subLiteralSuccess := compilation{}
	pushRegisterSuccess := compilation{}
	pushLiteralSuccess := compilation{}
	popSuccess := compilation{}
	readSuccess := compilation{}
	writeSuccess := compilation{}
	setRegisterSuccess := compilation{}
	setLiteralSuccess := compilation{}

	successCases := []compilation{
		addRegisterSuccess,
		addLiteralSuccess,
		subRegisterSuccess,
		subLiteralSuccess,
		pushRegisterSuccess,
		pushLiteralSuccess,
		popSuccess,
		readSuccess,
		writeSuccess,
		setRegisterSuccess,
		setLiteralSuccess,
	}

	for i, test := range successCases {
		ok := compileHelper(t, test.ast, MakeProcess(test.expected))

		if !ok {
			t.Errorf("Test failure at success case %d", i)
		}
	}

	failureCases := []*AST{}

	for i, test := range failureCases {
		ok := compileFailureHelper(t, test)

		if !ok {
			t.Errorf("Test failure at failure case %d", i)
		}
	}
}
