package asm

import (
	"bytes"
	"errors"
	"fmt"
)

type block struct {
	statements *[]Statement
}

type ASTBuilder struct {
	AST   *AST
	stack []block
}

func (builder *ASTBuilder) Append(statements ...Statement) {
	builder.prepare()

	for _, stmt := range statements {
		builder.appendToBlock(stmt)
	}
}

func (builder *ASTBuilder) OpenLoop(operand int) {
	builder.prepare()

	loop := &LoopStmt{}
	loop.Operand = operand
	builder.appendToBlock(loop)
	blk := block{statements: &loop.Nest}
	builder.stack = append(builder.stack, blk)
}

func (builder *ASTBuilder) appendToBlock(stmt Statement) {
	tip := builder.stack[len(builder.stack)-1]
	*tip.statements = append(*tip.statements, stmt)
}

func (builder *ASTBuilder) LeaveBlock() error {
	if len(builder.stack) == 0 {
		return errors.New("Tried to pop stack of length 0")
	}

	builder.stack = builder.stack[:len(builder.stack)-1]
	return nil
}

func (builder *ASTBuilder) prepare() {
	if builder.AST == nil {
		builder.AST = &AST{}
	}

	if len(builder.stack) == 0 {
		builder.stack = []block{
			block{statements: &builder.AST.Statements},
		}
	}
}

type AST struct {
	Statements []Statement
}

func (ast *AST) Visit(visitor ASTVisitor) {
	visitor.VisitAST(ast)

	for _, stmt := range ast.Statements {
		stmt.Visit(visitor)
	}

	visitor.LeaveAST(ast)
}

// FIXME use proper pretty printer with indentation.
func (ast *AST) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("AST { ")

	for _, stmt := range ast.Statements {
		buff.WriteString(fmt.Sprint(stmt))
		buff.WriteRune(' ')
	}
	buff.WriteRune('}')

	return buff.String()
}

type Statement interface {
	Visit(visitor ASTVisitor)
}

type ASTVisitor interface {
	VisitAST(ast *AST)
	LeaveAST(ast *AST)
	VisitLoop(loop *LoopStmt)
	LeaveLoop(loop *LoopStmt)
	VisitAdd(add *AddStmt)
	VisitSub(sub *SubStmt)
	VisitPush(push *PushStmt)
	VisitPop(pop *PopStmt)
	VisitRead(read *ReadStmt)
	VisitWrite(write *WriteStmt)
	VisitSet(set *SetStmt)
	VisitCopy(copy *CopyStmt)
	VisitJump(jump *JumpStmt)
	VisitCall(call *CallStmt)
}

type OneOperandStmt struct {
	Operand int
}

type TwoOperandStmt struct {
	Operand [2]int
}

type LoopStmt struct {
	OneOperandStmt
	Nest []Statement
}

type AddStmt struct {
	TwoOperandStmt
}

type SubStmt struct {
	TwoOperandStmt
}

type PushStmt struct {
	TwoOperandStmt
}

type PopStmt struct {
	TwoOperandStmt
}

type ReadStmt struct {
	OneOperandStmt
}

type WriteStmt struct {
	OneOperandStmt
}

type SetStmt struct {
	TwoOperandStmt
}

type CopyStmt struct {
	TwoOperandStmt
}

type JumpStmt struct {
	TwoOperandStmt
}

type CallStmt struct {
	VmFunc string
	OneOperandStmt
}

func (stmt *LoopStmt) Visit(visitor ASTVisitor) {
	visitor.VisitLoop(stmt)

	for _, nestedStmt := range stmt.Nest {
		nestedStmt.Visit(visitor)
	}

	visitor.LeaveLoop(stmt)
}

func (stmt *AddStmt) Visit(visitor ASTVisitor) {
	visitor.VisitAdd(stmt)
}

func (stmt *SubStmt) Visit(visitor ASTVisitor) {
	visitor.VisitSub(stmt)
}

func (stmt *PushStmt) Visit(visitor ASTVisitor) {
	visitor.VisitPush(stmt)
}

func (stmt *PopStmt) Visit(visitor ASTVisitor) {
	visitor.VisitPop(stmt)
}

func (stmt *ReadStmt) Visit(visitor ASTVisitor) {
	visitor.VisitRead(stmt)
}

func (stmt *WriteStmt) Visit(visitor ASTVisitor) {
	visitor.VisitWrite(stmt)
}

func (stmt *SetStmt) Visit(visitor ASTVisitor) {
	visitor.VisitSet(stmt)
}

func (stmt *CopyStmt) Visit(visitor ASTVisitor) {
	visitor.VisitCopy(stmt)
}

func (stmt *JumpStmt) Visit(visitor ASTVisitor) {
	visitor.VisitJump(stmt)
}

func (stmt *CallStmt) Visit(visitor ASTVisitor) {
	visitor.VisitCall(stmt)
}

func (loop *LoopStmt) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString(fmt.Sprintf("Loop(%d) { ", loop.Operand))

	for _, stmt := range loop.Nest {
		buff.WriteString(fmt.Sprint(stmt))
		buff.WriteRune(' ')
	}
	buff.WriteString("};")

	return buff.String()
}

func (stmt *AddStmt) String() string {
	return fmt.Sprintf("Add(%d, %d);", stmt.Operand[0], stmt.Operand[1])
}
func (stmt *SubStmt) String() string {
	return fmt.Sprintf("Sub(%d, %d);", stmt.Operand[0], stmt.Operand[1])
}
func (stmt *PushStmt) String() string {
	return fmt.Sprintf("Push(%d, %d);", stmt.Operand[0], stmt.Operand[1])
}
func (stmt *PopStmt) String() string {
	return fmt.Sprintf("Pop(%d, %d);", stmt.Operand[0], stmt.Operand[1])
}
func (stmt *ReadStmt) String() string {
	return fmt.Sprintf("Read(%d);", stmt.Operand)
}
func (stmt *WriteStmt) String() string {
	return fmt.Sprintf("Write(%d);", stmt.Operand)
}
func (stmt *SetStmt) String() string {
	return fmt.Sprintf("Set(%d, %d);", stmt.Operand[0], stmt.Operand[1])
}
func (stmt *CopyStmt) String() string {
	return fmt.Sprintf("Copy(%d, %d);", stmt.Operand[0], stmt.Operand[1])
}
func (stmt *JumpStmt) String() string {
	return fmt.Sprintf("Jump(%d, %d);", stmt.Operand[0], stmt.Operand[1])
}
func (stmt *CallStmt) String() string {
	return fmt.Sprintf("Call(%s, %d);", stmt.VmFunc, stmt.Operand)
}
