package shapes

import (
	"errors"
)

type CompileVisitor struct {
	Process *Process
	Error   error
}

func (c CompileVisitor) LeaveAST(ast *AST) {
	panic(errors.New("CompileVisitor.LeaveAST not implemented"))
}

func (c CompileVisitor) LeaveLoop(l *LoopStmt) {
	panic(errors.New("CompileVisitor.LeaveLoop not implemented"))
}

func (c CompileVisitor) VisitAST(ast *AST) {
	panic(errors.New("CompileVisitor.VisitAST not implemented"))
}

func (c CompileVisitor) VisitAdd(a *AddStmt) {
	panic(errors.New("CompileVisitor.VisitAdd not implemented"))
}

func (c CompileVisitor) VisitLoop(l *LoopStmt) {
	panic(errors.New("CompileVisitor.VisitLoop not implemented"))
}

func (c CompileVisitor) VisitPop(p *PopStmt) {
	panic(errors.New("CompileVisitor.VisitPop not implemented"))
}

func (c CompileVisitor) VisitPush(p *PushStmt) {
	panic(errors.New("CompileVisitor.VisitPush not implemented"))
}

func (c CompileVisitor) VisitRead(r *ReadStmt) {
	panic(errors.New("CompileVisitor.VisitRead not implemented"))
}

func (c CompileVisitor) VisitSet(s *SetStmt) {
	panic(errors.New("CompileVisitor.VisitSet not implemented"))
}

func (c CompileVisitor) VisitSub(s *SubStmt) {
	panic(errors.New("CompileVisitor.VisitSub not implemented"))
}

func (c CompileVisitor) VisitWrite(w *WriteStmt) {
	panic(errors.New("CompileVisitor.VisitWrite not implemented"))
}
