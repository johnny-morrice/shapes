package shapes

import (
	"errors"

	"github.com/johnny-morrice/shapes/asm"
)

type NameTable struct {
	Registers [MAX_WORD]string
	Stacks    [MAX_WORD]string
}

type CompileVisitor struct {
	NameTable
	Process *Process
	Error   error
}

func (c CompileVisitor) LeaveAST(ast *asm.AST) {
	panic(errors.New("CompileVisitor.LeaveAST not implemented"))
}

func (c CompileVisitor) LeaveLoop(l *asm.LoopStmt) {
	panic(errors.New("CompileVisitor.LeaveLoop not implemented"))
}

func (c CompileVisitor) VisitAST(ast *asm.AST) {
	panic(errors.New("CompileVisitor.VisitAST not implemented"))
}

func (c CompileVisitor) VisitAdd(a *asm.AddStmt) {
	panic(errors.New("CompileVisitor.VisitAdd not implemented"))
}

func (c CompileVisitor) VisitLoop(l *asm.LoopStmt) {
	panic(errors.New("CompileVisitor.VisitLoop not implemented"))
}

func (c CompileVisitor) VisitPop(p *asm.PopStmt) {
	panic(errors.New("CompileVisitor.VisitPop not implemented"))
}

func (c CompileVisitor) VisitPush(p *asm.PushStmt) {
	panic(errors.New("CompileVisitor.VisitPush not implemented"))
}

func (c CompileVisitor) VisitWrite(w *asm.WriteStmt) {
	panic(errors.New("CompileVisitor.VisitWrite not implemented"))
}

func (c CompileVisitor) VisitRead(r *asm.ReadStmt) {
	panic(errors.New("CompileVisitor.VisitRead not implemented"))
}

func (c CompileVisitor) VisitSet(s *asm.SetStmt) {
	panic(errors.New("CompileVisitor.VisitSet not implemented"))
}

func (c CompileVisitor) VisitCopy(s *asm.CopyStmt) {
	panic(errors.New("CompileVisitor.VisitCopy not implemented"))
}

func (c CompileVisitor) VisitSub(s *asm.SubStmt) {
	panic(errors.New("CompileVisitor.VisitSub not implemented"))
}
