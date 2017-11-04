package shapes

import (
	"log"
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
		callVmFuncCompilation(),
	}

	for i, test := range successCases {
		ok := compileHelper(t, test.ast, MakeProcess(test.expected))

		if !ok {
			t.Errorf("Test failure at success case %d", i)
		}
	}

	failureCases := []*asm.AST{
		unknownFunctionFailure(),
	}

	for i, test := range failureCases {
		ok := compileFailureHelper(t, test)

		if !ok {
			t.Errorf("Test failure at failure case %d", i)
		}
	}
}

func addRegisterCompilation() compilation {
	addStmt := &asm.AddStmt{}
	addStmt.Operand[0] = 68
	addStmt.Operand[1] = 99

	statements := []asm.Statement{
		addStmt,
	}
	expected := []Operation{
		Operation{
			OpCode:  OP_ADD,
			Operand: [2]Operand{68, 99},
		},
	}
	return makeCompilation(statements, expected)
}

func subRegisterCompilation() compilation {
	subStmt := &asm.SubStmt{}
	subStmt.Operand[0] = 31
	subStmt.Operand[1] = 243

	statements := []asm.Statement{
		subStmt,
	}
	expected := []Operation{
		Operation{
			OpCode:  OP_SUB,
			Operand: [2]Operand{31, 243},
		},
	}
	return makeCompilation(statements, expected)
}

func pushRegisterCompilation() compilation {
	pushStmt := &asm.PushStmt{}
	pushStmt.Operand[0] = 67
	pushStmt.Operand[1] = 123

	statements := []asm.Statement{
		pushStmt,
	}
	expected := []Operation{
		Operation{
			OpCode:  OP_PUSH,
			Operand: [2]Operand{67, 123},
		},
	}
	return makeCompilation(statements, expected)
}

func popCompilation() compilation {
	popStmt := &asm.PopStmt{}
	popStmt.Operand[0] = 90
	popStmt.Operand[1] = 80

	statements := []asm.Statement{
		popStmt,
	}
	expected := []Operation{
		Operation{
			OpCode:  OP_POP,
			Operand: [2]Operand{90, 80},
		},
	}
	return makeCompilation(statements, expected)
}

func readCompilation() compilation {
	readStmt := &asm.ReadStmt{}
	readStmt.Operand = 8

	statements := []asm.Statement{
		readStmt,
	}
	expected := []Operation{
		Operation{
			OpCode:  OP_READ,
			Operand: [2]Operand{8},
		},
	}
	return makeCompilation(statements, expected)
}

func writeCompilation() compilation {
	writeStmt := &asm.WriteStmt{}
	writeStmt.Operand = 8

	statements := []asm.Statement{
		writeStmt,
	}
	expected := []Operation{
		Operation{
			OpCode:  OP_WRITE,
			Operand: [2]Operand{8},
		},
	}
	return makeCompilation(statements, expected)
}

func setRegisterCompilation() compilation {
	setStmt := &asm.SetStmt{}
	setStmt.Operand[0] = 1
	setStmt.Operand[1] = 20

	statements := []asm.Statement{
		setStmt,
	}
	expected := []Operation{
		Operation{
			OpCode:  OP_SET,
			Operand: [2]Operand{1, 20},
		},
	}
	return makeCompilation(statements, expected)
}

func callVmFuncCompilation() compilation {
	callStmtA := &asm.CallStmt{}
	callStmtA.VmFunc = VM_FUNC_A
	callStmtA.Operand = 42

	callStmtB := &asm.CallStmt{}
	callStmtB.VmFunc = VM_FUNC_B
	callStmtB.Operand = 65

	statements := []asm.Statement{
		callStmtA,
		callStmtB,
	}

	expected := []Operation{
		Operation{
			OpCode:  OP_CALL,
			Operand: [2]Operand{0, 42},
		},
		Operation{
			OpCode:  OP_CALL,
			Operand: [2]Operand{1, 65},
		},
	}
	return makeCompilation(statements, expected)
}

func copyRegisterCompilation() compilation {
	copyStmt := &asm.CopyStmt{}
	copyStmt.Operand[0] = 255
	copyStmt.Operand[1] = 34

	statements := []asm.Statement{
		copyStmt,
	}
	expected := []Operation{
		Operation{
			OpCode:  OP_COPY,
			Operand: [2]Operand{255, 34},
		},
	}
	return makeCompilation(statements, expected)
}

func loopCompilation() compilation {
	writeStmt := &asm.WriteStmt{}
	writeStmt.Operand = 8

	subStmt := &asm.SubStmt{}
	subStmt.Operand[0] = 8
	subStmt.Operand[1] = 1

	setStmt := &asm.SetStmt{}
	setStmt.Operand[0] = 8
	setStmt.Operand[1] = 10

	loopStmt := &asm.LoopStmt{}
	loopStmt.Operand = 8
	loopStmt.Nest = []asm.Statement{
		writeStmt,
		subStmt,
	}

	statements := []asm.Statement{
		setStmt,
		loopStmt,
	}

	expected := []Operation{
		Operation{
			OpCode:  OP_SET,
			Operand: [2]Operand{8, 10},
		},
		Operation{
			OpCode:  OP_WRITE,
			Operand: [2]Operand{8, 0},
		},
		Operation{
			OpCode:  OP_SUB,
			Operand: [2]Operand{8, 1},
		},
		Operation{
			OpCode:  OP_JMPNZ,
			Operand: [2]Operand{8, 1},
		},
	}

	return makeCompilation(statements, expected)
}

func severalStatementsCompilation() compilation {
	writeStmt := &asm.WriteStmt{}
	writeStmt.Operand = 90

	subStmt := &asm.SubStmt{}
	subStmt.Operand[0] = 32
	subStmt.Operand[1] = 15

	setStmt := &asm.SetStmt{}
	setStmt.Operand[0] = 56
	setStmt.Operand[1] = 22

	statements := []asm.Statement{
		writeStmt,
		subStmt,
		setStmt,
	}

	expected := []Operation{
		Operation{
			OpCode:  OP_WRITE,
			Operand: [2]Operand{90, 0},
		},
		Operation{
			OpCode:  OP_SUB,
			Operand: [2]Operand{32, 15},
		},
		Operation{
			OpCode:  OP_SET,
			Operand: [2]Operand{56, 22},
		},
	}
	return makeCompilation(statements, expected)
}

func unknownFunctionFailure() *asm.AST {
	ast := &asm.AST{}
	ast.Statements = []asm.Statement{
		&asm.CallStmt{
			VmFunc: VM_FUNC_NO_EXIST,
		},
	}
	return ast
}

func compileHelper(t *testing.T, ast *asm.AST, expected *Process) bool {
	t.Helper()

	actual, err := Compile(ast, LibTest())

	if err != nil {
		t.Errorf("Expected compile success but got error: %s", err.Error())
		return false
	}

	if !expected.IsSameByteCode(actual) {
		t.Error("Expected same bytecode")

		log.Print("Expected bytecode")
		logByteCode(expected)

		log.Print("Actual bytecode")
		logByteCode(actual)

		return false
	}

	return true
}

func logByteCode(process *Process) {
	for i, byteCode := range process.ByteCode {
		log.Printf("%d %v", i, byteCode)
	}
}

func compileFailureHelper(t *testing.T, ast *asm.AST) bool {
	t.Helper()

	_, err := Compile(ast, LibTest())

	if err == nil {
		t.Error("Expected compile failure")
		return false
	}

	return true
}

func LibTest() *Library {
	lib := &Library{}
	lib.AddFunction(VM_FUNC_A, vmFuncA)
	lib.AddFunction(VM_FUNC_B, vmFuncB)
	return lib
}

func vmFuncA(runtime *Runtime, stackAddr Address) {
	runtime.Process.Pop(stackAddr)
	runtime.Process.Push(stackAddr, VM_FUNC_A_RETURN)
	runtime.Process.IncrementPC()
}

func vmFuncB(runtime *Runtime, stackAddr Address) {
	runtime.Process.Pop(stackAddr)
	runtime.Process.Push(stackAddr, VM_FUNC_B_RETURN)
	runtime.Process.IncrementPC()
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

const VM_FUNC_A = "vm_func_a"
const VM_FUNC_B = "vm_func_b"
const VM_FUNC_NO_EXIST = "vm_func_no_exist"
const VM_FUNC_A_RETURN = 1
const VM_FUNC_B_RETURN = 2
