package asm

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
}

type LoopStmt struct {
	OneOperandStmt
	Nest []Statement
}

type OneOperandStmt struct {
	Operand byte
}

type TwoOperandStmt struct {
	Operands [2]byte
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
