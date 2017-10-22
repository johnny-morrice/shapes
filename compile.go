package shapes

import (
	"errors"

	"github.com/johnny-morrice/shapes/asm"
)

type loopStack struct {
	startAddress []int
}

type CompileVisitor struct {
	Process *Process
	Error   error

	loopStack loopStack
}

func (c *CompileVisitor) VisitAST(ast *asm.AST) {
}

func (c *CompileVisitor) LeaveAST(ast *asm.AST) {

}

func (c *CompileVisitor) VisitLoop(l *asm.LoopStmt) {
	startAddress := len(c.Process.ByteCode) + 1
	c.loopStack.startAddress = append(c.loopStack.startAddress, startAddress)
}

func (c *CompileVisitor) LeaveLoop(l *asm.LoopStmt) {
	tipIndex := len(c.loopStack.startAddress) - 1
	startAddress := c.loopStack.startAddress[tipIndex]
	c.loopStack.startAddress = c.loopStack.startAddress[:tipIndex]
	c.appendByteCode(Operation{
		OpCode:  OP_JMPNZ,
		Operand: [2]Operand{Operand(l.Operand), Operand(startAddress)},
	})
}

func (c *CompileVisitor) VisitAdd(add *asm.AddStmt) {
	c.appendByteCode(Operation{
		OpCode:  OP_ADD,
		Operand: [2]Operand{Operand(add.Operand[0]), Operand(add.Operand[1])},
	})
}

func (c *CompileVisitor) VisitPop(p *asm.PopStmt) {
	panic(errors.New("CompileVisitor.VisitPop not implemented"))
}

func (c *CompileVisitor) VisitPush(p *asm.PushStmt) {
	panic(errors.New("CompileVisitor.VisitPush not implemented"))
}

func (c *CompileVisitor) VisitWrite(w *asm.WriteStmt) {
	panic(errors.New("CompileVisitor.VisitWrite not implemented"))
}

func (c *CompileVisitor) VisitRead(r *asm.ReadStmt) {
	panic(errors.New("CompileVisitor.VisitRead not implemented"))
}

func (c *CompileVisitor) VisitSet(s *asm.SetStmt) {
	panic(errors.New("CompileVisitor.VisitSet not implemented"))
}

func (c *CompileVisitor) VisitCopy(s *asm.CopyStmt) {
	panic(errors.New("CompileVisitor.VisitCopy not implemented"))
}

func (c *CompileVisitor) VisitSub(s *asm.SubStmt) {
	panic(errors.New("CompileVisitor.VisitSub not implemented"))
}

func (c *CompileVisitor) appendByteCode(operations ...Operation) {
	c.Process.ByteCode = append(c.Process.ByteCode, operations...)
}

const MAX_BYTECODE = 65535
