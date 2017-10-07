package shapes

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

func ExecuteProgramCode(source []byte, input io.Reader, output io.Writer) error {
	const errMsg = "ExecuteProgramCode failed"

	ast, err := Parse(source)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	process, err := Compile(ast)

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	runtime := MakeRuntime(process, input, output)

	err = runtime.Execute()

	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

type AST struct {
	Statements []Statement
}

type Statement interface {
	Visit(visitor ASTVisitor)
}

type ASTVisitor interface {
	VisitLoop(loop *LoopStmt)
	VisitAdd(add *AddStmt)
	VisitSub(sub *SubStmt)
	VisitPush(push *PushStmt)
	VisitPop(pop *PopStmt)
	VisitRead(read *ReadStmt)
	VisitWrite(write *WriteStmt)
	VisitSet(set *SetStmt)
}

type LoopStmt struct {
	Operand string
	Nest    []Statement
}

func (stmt *LoopStmt) Visit(visitor ASTVisitor) {
	visitor.VisitLoop(stmt)
}

type TwoOperandStmt struct {
	Operands [2]string
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
	TwoOperandStmt
}

type WriteStmt struct {
	TwoOperandStmt
}

type SetStmt struct {
	TwoOperandStmt
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

type StatementType uint8

const (
	STMT_LOOP = iota
	STMT_ADD
	STMT_SUB
	STMT_PUSH
	STMT_POP
	STMT_READ
	STMT_WRITE
	STMT_SET
)

func Parse(source []byte) (*AST, error) {
	panic("not implemented")
}

type Address uint8

type LongAddress uint16

type Operand uint8

type Operation struct {
	OpCode  OpCode
	Operand [2]Operand
}

func (op Operation) LongAddress() LongAddress {
	long := LongAddress(op.Operand[0])
	long = long << 8
	long = long | LongAddress(op.Operand[1])
	return long
}

func (op Operation) Address(operand Address) Address {
	return Address(op.Operand[operand])
}

type OpCode byte

const (
	OP_JMPNZ = OpCode(iota)
	OP_ADD
	OP_SUB
	OP_PUSH
	OP_POP
	OP_READ
	OP_WRITE
	OP_COPY
	OP_SET
)

type Process struct {
	PC       LongAddress
	ByteCode []Operation
	Register [MAX_WORD]byte
	Stack    [MAX_WORD][]byte
	Error    error
}

func Compile(ast *AST) (*Process, error) {
	panic("not implemented")
}

func (process *Process) IsTerminated() bool {
	return process.PC >= LongAddress(len(process.ByteCode)) || process.Error != nil
}

func (process *Process) ExecuteStep(callTable []RuntimeCall) {
	op := process.ByteCode[process.PC]

	impl := callTable[op.OpCode]

	impl(op)
}

func (process *Process) Peek(stackAddr Address) byte {
	stack := process.Stack[stackAddr]

	if len(stack) == 0 {
		process.failEmptyStack(stackAddr)
		return 0
	}

	return stack[len(stack)-1]
}

func (process *Process) Pop(stackAddr Address) byte {
	tip := process.Peek(stackAddr)

	if process.Error != nil {
		return 0
	}

	stack := process.Stack[stackAddr]
	process.Stack[stackAddr] = stack[:len(stack)-1]

	return tip
}

func (process *Process) failEmptyStack(stackAddr Address) {
	process.Error = fmt.Errorf("stack was empty '%d'", stackAddr)
}

func (process *Process) Push(stackAddr Address, tip byte) {
	process.Stack[stackAddr] = append(process.Stack[stackAddr], tip)
}

func (process *Process) GetRegister(register Address) byte {
	return process.Register[register]
}

func (process *Process) SetRegister(register Address, val byte) {
	process.Register[register] = val
}

type RuntimeCall func(op Operation)

type Runtime struct {
	Process     *Process
	Error       error
	CallTable   []RuntimeCall
	Input       io.Reader
	Output      io.Writer
	readBuffer  []byte
	writeBuffer []byte
}

func MakeRuntime(process *Process, input io.Reader, output io.Writer) *Runtime {
	runtime := &Runtime{
		Process:     process,
		Input:       input,
		Output:      output,
		readBuffer:  []byte{0},
		writeBuffer: []byte{0},
	}

	runtime.CallTable = []RuntimeCall{
		runtime.jmpnz,
		runtime.add,
		runtime.sub,
		runtime.push,
		runtime.pop,
		runtime.read,
		runtime.write,
		runtime.copy,
		runtime.set,
	}

	return runtime
}

func (runtime *Runtime) hasError() bool {
	return runtime.Process.Error != nil
}

func (runtime *Runtime) jmpnz(op Operation) {
	val := runtime.Process.GetRegister(op.Address(0))

	if runtime.hasError() {
		return
	}

	if val != 0 {
		runtime.Process.PC = op.LongAddress()
	}
}

func (runtime *Runtime) onRegisters(op Operation, f func(valZero, valOne byte) byte) {
	valZero := runtime.Process.GetRegister(op.Address(0))

	if runtime.hasError() {
		return
	}

	valOne := runtime.Process.GetRegister(op.Address(1))

	if runtime.hasError() {
		return
	}

	newVal := f(valZero, valOne)
	runtime.Process.SetRegister(op.Address(0), newVal)
}

func (runtime *Runtime) add(op Operation) {
	runtime.onRegisters(op, func(valZero, valOne byte) byte {
		return valZero + valOne
	})
}

func (runtime *Runtime) sub(op Operation) {
	runtime.onRegisters(op, func(valZero, valOne byte) byte {
		return valZero - valOne
	})
}

func (runtime *Runtime) push(op Operation) {
	val := runtime.Process.GetRegister(op.Address(1))

	if runtime.hasError() {
		return
	}

	runtime.Process.Push(op.Address(0), val)
}

func (runtime *Runtime) pop(op Operation) {
	tip := runtime.Process.Pop(op.Address(0))

	if runtime.hasError() {
		return
	}

	runtime.Process.SetRegister(op.Address(1), tip)
}

func (runtime *Runtime) read(op Operation) {
	const errMsg = "Runtime.read failed"

	_, err := runtime.Input.Read(runtime.readBuffer)

	if err != nil {
		runtime.Process.Error = errors.Wrap(err, errMsg)
		return
	}

	runtime.Process.SetRegister(op.Address(0), runtime.readBuffer[0])
}

func (runtime *Runtime) write(op Operation) {
	const errMsg = "Runtime.Write failed"

	val := runtime.Process.GetRegister(op.Address(0))
	runtime.writeBuffer[0] = val

	_, err := runtime.Output.Write(runtime.writeBuffer)

	if err != nil {
		runtime.Process.Error = errors.Wrap(err, errMsg)
		return
	}
}

func (runtime *Runtime) copy(op Operation) {
	val := runtime.Process.GetRegister(op.Address(1))

	if runtime.hasError() {
		return
	}

	runtime.Process.SetRegister(op.Address(0), val)
}

func (runtime *Runtime) set(op Operation) {
	runtime.Process.SetRegister(op.Address(0), byte(op.Address(1)))
}

func (runtime *Runtime) Execute() error {
	for !runtime.Process.IsTerminated() {
		runtime.Process.ExecuteStep(runtime.CallTable)
	}

	if runtime.Process.Error != nil {
		return runtime.Process.Error
	}

	return nil
}

const MAX_WORD = 256
