package shapes

import (
	"github.com/johnny-morrice/shapes/asm"
)

type InfiniteTape struct {
	left  []uint64
	right []uint64
	// Negative index is on left tape.
	index int
}

func (tape *InfiniteTape) MoveHead(offset int) {
	tape.index += offset
}

func (tape *InfiniteTape) ReadHead() uint64 {
	return *tape.getCell()
}

func (tape *InfiniteTape) WriteHead(val uint64) {
	cell := tape.getCell()
	*cell = val
}

func (tape *InfiniteTape) getCell() *uint64 {
	if tape.index >= 0 {
		extendTape(&tape.right, tape.index)
		return &tape.right[tape.index]
	} else {
		leftIndex := (-tape.index) - 1
		extendTape(&tape.left, leftIndex)
		return &tape.left[leftIndex]
	}
}

func extendTape(tape *[]uint64, index int) {
	t := *tape
	grow := (index - len(t)) + 1

	if grow > 0 {
		*tape = append(t, make([]uint64, grow)...)
	}
}

type infiniteTapeList struct {
	list []*InfiniteTape
}

func (tape *infiniteTapeList) NewTape() int {
	index := len(tape.list)
	tape.list = append(tape.list, &InfiniteTape{})
	return index
}

func (tape *infiniteTapeList) MoveHead(index int, offset int) {
	tape.list[index].MoveHead(offset)
}

func (tape *infiniteTapeList) ReadHead(index int) uint64 {
	return tape.list[index].ReadHead()
}

func (tape *infiniteTapeList) WriteHead(index int, val uint64) {
	tape.list[index].WriteHead(val)
}

type InfiniteTapeVmWrapper struct {
	list *infiniteTapeList
}

func (tape *InfiniteTapeVmWrapper) NewTape(runtime *Runtime, stackAddr Address) {
	if tape.list == nil {
		tape.list = &infiniteTapeList{}
	}

	runtime.Process.Pop(stackAddr)
	index := tape.list.NewTape()
	runtime.Process.Push(stackAddr, uint64(index))
	runtime.Process.IncrementPC()
}

func (tape *InfiniteTapeVmWrapper) MoveHead(runtime *Runtime, stackAddr Address) {
	runtime.Process.Pop(stackAddr)
	index := runtime.Process.Pop(stackAddr)
	offset := runtime.Process.Pop(stackAddr)
	tape.list.MoveHead(int(index), int(offset))
	runtime.Process.IncrementPC()
}

func (tape *InfiniteTapeVmWrapper) ReadHead(runtime *Runtime, stackAddr Address) {
	runtime.Process.Pop(stackAddr)
	index := runtime.Process.Pop(stackAddr)
	val := tape.list.ReadHead(int(index))
	runtime.Process.Push(stackAddr, val)
	runtime.Process.IncrementPC()
}

func (tape *InfiniteTapeVmWrapper) WriteHead(runtime *Runtime, stackAddr Address) {
	runtime.Process.Pop(stackAddr)
	index := runtime.Process.Pop(stackAddr)
	val := runtime.Process.Pop(stackAddr)
	tape.list.WriteHead(int(index), val)
	runtime.Process.IncrementPC()
}

func init() {
	lib := &Library{}
	tape := &InfiniteTapeVmWrapper{}
	lib.AddFunction(asm.TAPE_NEW, tape.NewTape)
	lib.AddFunction(asm.TAPE_MOVE_HEAD, tape.MoveHead)
	lib.AddFunction(asm.TAPE_READ_HEAD, tape.ReadHead)
	lib.AddFunction(asm.TAPE_WRITE_HEAD, tape.WriteHead)

	StdLib().AddLibrary(lib)
}
