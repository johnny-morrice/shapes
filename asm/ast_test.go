package asm

import (
	"reflect"
	"testing"
)

func TestASTBuilderAppend(t *testing.T) {
	stmt1 := &SetStmt{}
	stmt2 := &CopyStmt{}
	stmt3 := &WriteStmt{}

	expected := &AST{
		Statements: []Statement{
			stmt1,
			stmt2,
			stmt3,
		},
	}

	builder := &ASTBuilder{}
	builder.Append(stmt1)
	builder.Append(stmt2, stmt3)

	assertASTEqual(t, expected, builder.AST)
}

func TestASTBuilderOpenLoop_OneLoop(t *testing.T) {
	stmt1 := &SetStmt{}
	stmt2 := &CopyStmt{}
	stmt3 := &WriteStmt{}

	loop := &LoopStmt{
		Nest: []Statement{stmt2, stmt3},
	}
	loop.Operand = 2

	expected := &AST{
		Statements: []Statement{
			stmt1,
			loop,
		},
	}

	builder := &ASTBuilder{}
	builder.Append(stmt1)
	builder.OpenLoop(2)
	builder.Append(stmt2, stmt3)
	leaveBlock(t, builder)

	assertASTEqual(t, expected, builder.AST)
}

func TestASTBuilderOpenLoop_NestedLoop(t *testing.T) {
	stmt1 := &SetStmt{}
	stmt2 := &CopyStmt{}
	stmt3 := &WriteStmt{}
	stmt4 := &AddStmt{}

	loop2 := &LoopStmt{
		Nest: []Statement{stmt3, stmt4},
	}
	loop2.Operand = 3
	loop1 := &LoopStmt{
		Nest: []Statement{stmt2, loop2},
	}
	loop1.Operand = 2

	expected := &AST{
		Statements: []Statement{
			stmt1,
			loop1,
		},
	}

	builder := &ASTBuilder{}
	builder.Append(stmt1)
	builder.OpenLoop(2)
	builder.Append(stmt2)
	builder.OpenLoop(3)
	builder.Append(stmt3, stmt4)
	leaveBlock(t, builder)
	leaveBlock(t, builder)

	assertASTEqual(t, expected, builder.AST)
}

func leaveBlock(t *testing.T, builder *ASTBuilder) {
	t.Helper()

	err := builder.LeaveBlock()

	if err != nil {
		t.Errorf("Error in LeaveBlock: %s", err.Error())
		t.FailNow()
		return
	}
}

func assertASTEqual(t *testing.T, expected *AST, actual *AST) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v but received %v", expected, actual)
	}
}
