package shapes

import (
	"github.com/johnny-morrice/shapes/asm"
)

type loopStack struct {
	startAddress []int
}

type CompileVisitor struct {
	Process *Process
	Error   error
	Library *Library

	loopStack loopStack
}

func (c *CompileVisitor) VisitAST(ast *asm.AST) {
	c.Process = &Process{}
}

func (c *CompileVisitor) LeaveAST(ast *asm.AST) {

}

func (c *CompileVisitor) VisitLoop(l *asm.LoopStmt) {
	startAddress := len(c.Process.ByteCode)
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

func (c *CompileVisitor) VisitAdd(stmt *asm.AddStmt) {
	c.appendByteCode(twoOp(OP_ADD, stmt.TwoOperandStmt))
}

func (c *CompileVisitor) VisitPop(stmt *asm.PopStmt) {
	c.appendByteCode(twoOp(OP_POP, stmt.TwoOperandStmt))
}

func (c *CompileVisitor) VisitPush(stmt *asm.PushStmt) {
	c.appendByteCode(twoOp(OP_PUSH, stmt.TwoOperandStmt))
}

func (c *CompileVisitor) VisitWrite(stmt *asm.WriteStmt) {
	c.appendByteCode(oneOp(OP_WRITE, stmt.OneOperandStmt))
}

func (c *CompileVisitor) VisitRead(stmt *asm.ReadStmt) {
	c.appendByteCode(oneOp(OP_READ, stmt.OneOperandStmt))
}

func (c *CompileVisitor) VisitSet(stmt *asm.SetStmt) {
	c.appendByteCode(twoOp(OP_SET, stmt.TwoOperandStmt))
}

func (c *CompileVisitor) VisitCopy(stmt *asm.CopyStmt) {
	c.appendByteCode(twoOp(OP_COPY, stmt.TwoOperandStmt))
}

func (c *CompileVisitor) VisitSub(stmt *asm.SubStmt) {
	c.appendByteCode(twoOp(OP_SUB, stmt.TwoOperandStmt))
}

func (c *CompileVisitor) VisitJump(stmt *asm.JumpStmt) {
	c.appendByteCode(twoOp(OP_JMPNZ, stmt.TwoOperandStmt))
}

func (c *CompileVisitor) VisitCall(stmt *asm.CallStmt) {
	index, err := c.Library.GetFunctionIndex(stmt.VmFunc)

	if err != nil {
		c.Error = err
		return
	}

	twoOpStmt := asm.TwoOperandStmt{}
	twoOpStmt.Operand[0] = index
	twoOpStmt.Operand[1] = stmt.Operand
	c.appendByteCode(twoOp(OP_CALL, twoOpStmt))
}

func (c *CompileVisitor) appendByteCode(operations ...Operation) {
	c.Process.ByteCode = append(c.Process.ByteCode, operations...)
}

func twoOp(opCode OpCode, stmt asm.TwoOperandStmt) Operation {
	return Operation{
		OpCode:  opCode,
		Operand: [2]Operand{makeOperand(stmt.Operand[0]), makeOperand(stmt.Operand[1])},
	}
}

func oneOp(opCode OpCode, stmt asm.OneOperandStmt) Operation {
	return Operation{
		OpCode:  opCode,
		Operand: [2]Operand{makeOperand(stmt.Operand), 0},
	}
}

func makeOperand(op int) Operand {
	return Operand(uint64(op))
}

const MAX_BYTECODE = 65535
