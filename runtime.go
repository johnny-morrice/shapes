package shapes

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/johnny-morrice/shapes/asm"
)

type Address uint64

type Operand uint64

type Operation struct {
	OpCode  OpCode
	Operand [2]Operand
}

func (op Operation) String() string {
	return fmt.Sprintf("%v %d %d", op.OpCode, op.Operand[0], op.Operand[1])
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
	OP_CALL
)

var __OPCODE_STRING = []string{
	"JMPNZ",
	"ADD",
	"SUB",
	"PUSH",
	"POP",
	"READ",
	"WRITE",
	"COPY",
	"SET",
	"CALL",
}

func (opCode OpCode) String() string {
	return __OPCODE_STRING[opCode]
}

type Process struct {
	PC       Address
	ByteCode []Operation
	Register [REGISTER_COUNT]uint64
	Stack    [REGISTER_COUNT][]uint64
	Error    error
}

func MakeProcess(byteCode []Operation) *Process {
	return &Process{
		ByteCode: byteCode,
	}
}

func Compile(ast *asm.AST, lib *Library) (*Process, error) {
	compiler := &CompileVisitor{
		Library: lib,
	}
	ast.Visit(compiler)

	return compiler.Process, compiler.Error
}

func (process *Process) DumpBytecode(w io.Writer) {
	intPC := int(process.PC)
	for i, op := range process.ByteCode {
		fmt.Fprint(w, i)
		if i == intPC {
			fmt.Fprint(w, " *")
		}
		fmt.Fprintf(w, " %v", op)
		fmt.Fprint(w, "\n")
	}
}

func (process *Process) DumpRegisters(w io.Writer) {
	fmt.Println("REGISTER DUMP:\n")

	for i, reg := range process.Register {
		fmt.Fprintf(w, "%d %d\n", i, reg)
	}

	fmt.Println("STACK DUMP:\n")
	for i, stk := range process.Stack {
		if len(stk) == 0 {
			fmt.Fprintf(w, "%d EMPTY\n", i)
		} else {
			for j, val := range stk {
				fmt.Fprintf(w, "%d %d %d\n", i, j, val)
			}
		}
	}
}

func (process *Process) IsSameByteCode(other *Process) bool {
	if len(process.ByteCode) != len(other.ByteCode) {
		return false
	}

	for i, myByteCode := range process.ByteCode {
		theirByteCode := other.ByteCode[i]

		if myByteCode != theirByteCode {
			return false
		}
	}

	return true
}

func (process *Process) IsTerminated() bool {
	return int(process.PC) >= len(process.ByteCode) || process.Error != nil
}

func (process *Process) ExecuteStep(callTable []RuntimeCall) {
	op := process.ByteCode[process.PC]

	impl := callTable[op.OpCode]

	impl(op)
}

func (process *Process) Peek(stackAddr Address) uint64 {
	stack := process.Stack[stackAddr]

	if len(stack) == 0 {
		process.failEmptyStack(stackAddr)
		return 0
	}

	return stack[len(stack)-1]
}

func (process *Process) Pop(stackAddr Address) uint64 {
	tip := process.Peek(stackAddr)

	if process.Error != nil {
		return 0
	}

	stack := process.Stack[stackAddr]
	process.Stack[stackAddr] = stack[:len(stack)-1]

	return tip
}

func (process *Process) failEmptyStack(stackAddr Address) {
	process.Error = fmt.Errorf("stack %d was empty", stackAddr)
}

func (process *Process) Push(stackAddr Address, tip uint64) {
	process.Stack[stackAddr] = append(process.Stack[stackAddr], tip)
}

func (process *Process) GetRegister(register Address) uint64 {
	return process.Register[register]
}

func (process *Process) SetRegister(register Address, val uint64) {
	process.Register[register] = val
}

func (process *Process) IncrementPC() {
	process.PC++
}

type RuntimeCall func(op Operation)

type Runtime struct {
	RuntimeBuilder
	Error       error
	callTable   []RuntimeCall
	readBuffer  []byte
	writeBuffer []byte
}

type RuntimeBuilder struct {
	Process *Process
	Library *Library
	Input   io.Reader
	Output  io.Writer
}

func (builder *RuntimeBuilder) Build() *Runtime {
	runtime := &Runtime{
		RuntimeBuilder: *builder,
		readBuffer:     []byte{0},
		writeBuffer:    []byte{0},
	}

	runtime.callTable = []RuntimeCall{
		runtime.jmpnz,
		runtime.add,
		runtime.sub,
		runtime.push,
		runtime.pop,
		runtime.read,
		runtime.write,
		runtime.copy,
		runtime.set,
		runtime.call,
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
		runtime.Process.PC = op.Address(1)
	} else {
		runtime.Process.IncrementPC()
	}
}

func (runtime *Runtime) onRegisters(op Operation, f func(valZero, valOne uint64) uint64) {
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
	runtime.onRegisters(op, func(valZero, valOne uint64) uint64 {
		return valZero + valOne
	})
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) sub(op Operation) {
	runtime.onRegisters(op, func(valZero, valOne uint64) uint64 {
		return valZero - valOne
	})
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) push(op Operation) {
	val := runtime.Process.GetRegister(op.Address(1))

	if runtime.hasError() {
		return
	}

	runtime.Process.Push(op.Address(0), val)
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) pop(op Operation) {
	tip := runtime.Process.Pop(op.Address(0))

	if runtime.hasError() {
		return
	}

	runtime.Process.SetRegister(op.Address(1), tip)
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) read(op Operation) {
	const errMsg = "Runtime.read failed"

	_, err := runtime.Input.Read(runtime.readBuffer)

	if err != nil {
		runtime.Process.Error = errors.Wrap(err, errMsg)
		return
	}

	runtime.Process.SetRegister(op.Address(0), uint64(runtime.readBuffer[0]))
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) write(op Operation) {
	const errMsg = "Runtime.Write failed"

	val := runtime.Process.GetRegister(op.Address(0))
	runtime.writeBuffer[0] = byte(val)

	_, err := runtime.Output.Write(runtime.writeBuffer)

	if err != nil {
		runtime.Process.Error = errors.Wrap(err, errMsg)
		return
	}

	runtime.Process.IncrementPC()
}

func (runtime *Runtime) copy(op Operation) {
	val := runtime.Process.GetRegister(op.Address(1))

	if runtime.hasError() {
		return
	}

	runtime.Process.SetRegister(op.Address(0), val)
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) set(op Operation) {
	runtime.Process.SetRegister(op.Address(0), uint64(op.Operand[1]))
	runtime.Process.IncrementPC()
}

func (runtime *Runtime) call(op Operation) {
	callee := runtime.Library.Functions[op.Operand[0]]
	runtime.Process.Push(op.Address(1), uint64(runtime.Process.PC))
	callee(runtime, op.Address(1))
	// Callee moves PC.
}

func (runtime *Runtime) Execute() error {
	for !runtime.Process.IsTerminated() {
		runtime.Process.ExecuteStep(runtime.callTable)
	}

	if runtime.Process.Error != nil {
		// Maybe we should expose an error buffer to the process too.
		runtime.DebugDump(os.Stderr)
		return runtime.Process.Error
	}

	return nil
}

func (runtime *Runtime) DebugDump(w io.Writer) {
	buff := bufio.NewWriter(os.Stderr)
	fmt.Fprint(w, "FUNCTION MAP:\n")

	for i, _ := range runtime.Library.Functions {
		fname, err := runtime.Library.GetFunctionName(i)

		if err == nil {
			fmt.Fprintf(w, "%d %s\n", i, fname)
		} else {
			fmt.Fprintf(w, "%d UNKNOWN FUNCTION\n")
		}
	}

	fmt.Fprint(w, "BYTECODE:\n")
	runtime.Process.DumpRegisters(buff)
	runtime.Process.DumpBytecode(buff)
	buff.Flush()
}

const REGISTER_COUNT = 256
