package shapes

import (
	"testing"
)

func TestInfiniteTapeMoveHead(t *testing.T) {
	const expectA = 97
	const expectB = 82
	const expectC = 31
	tape := &InfiniteTape{}
	tape.WriteHead(expectA)
	tape.MoveHead(100)
	tape.WriteHead(expectB)
	tape.MoveHead(-300)
	tape.WriteHead(expectC)

	tape.MoveHead(200)
	readExpect(t, expectA, tape)
	tape.MoveHead(100)
	readExpect(t, expectB, tape)
	tape.MoveHead(-300)
	readExpect(t, expectC, tape)
	tape.MoveHead(200)
	readExpect(t, expectA, tape)
}

func TestInfiniteTapeReadWriteHead(t *testing.T) {
	const expect = 64
	tape := &InfiniteTape{}
	tape.WriteHead(expect)
	readExpect(t, expect, tape)
}

func readExpect(t *testing.T, expect uint64, tape *InfiniteTape) {
	t.Helper()
	actual := tape.ReadHead()

	if expect != actual {
		t.Errorf("Expected %d but received %d", expect, actual)
	}
}
