package shapes

import (
	"bytes"
	"testing"
)

type cannedProcess struct {
	process *Process
	input   *bytes.Buffer
	output  *bytes.Buffer
}

func (canned cannedProcess) makeRuntime() *Runtime {
	return MakeRuntime(canned.process, canned.input, canned.output)
}

func (canned cannedProcess) sameRegisters(t *testing.T, other cannedProcess) bool {
	t.Helper()

	for i, myReg := range canned.process.Register {
		theirReg := other.process.Register[i]
		if myReg != theirReg {
			t.Errorf("Mismatched register %d: expected %d but was %d", i, myReg, theirReg)
			return false
		}
	}

	return true
}
func (canned cannedProcess) sameStack(t *testing.T, other cannedProcess) bool {
	t.Helper()

	for i, myStack := range canned.process.Stack {
		theirStack := other.process.Stack[i]
		myStackLen := len(canned.process.Stack)
		theirStackLen := len(other.process.Stack)
		if myStackLen != theirStackLen {
			t.Errorf("Expected stack %d length %d but was %d", i, myStackLen, theirStackLen)
			return false
		}

		for j, myEntry := range myStack {
			theirEntry := theirStack[j]

			if myEntry != theirEntry {
				t.Errorf("Mismatched stack %d entry %d: expected %d but was %d", i, j, myEntry, theirEntry)
				return false
			}
		}

	}

	return true
}
func (canned cannedProcess) sameOutput(t *testing.T, other cannedProcess) bool {
	t.Helper()

	myOutput := canned.output.Bytes()
	theirOutput := other.output.Bytes()

	myOutputLen := len(myOutput)
	theirOutputLen := len(theirOutput)
	if myOutputLen != theirOutputLen {
		t.Errorf("Expected stack length %d but was %d", myOutputLen, theirOutputLen)
		return false
	}

	for i, myOut := range myOutput {
		theirOut := theirOutput[i]

		if myOut != theirOut {
			t.Errorf("Mismatched stack entry %d: expected %d but was %d", i, myOut, theirOut)
			return false
		}
	}

	return true
}

func executeHelper(t *testing.T, input cannedProcess, expect cannedProcess) {
	t.Helper()

	runtime := input.makeRuntime()
	err := runtime.Execute()

	ok := (err == nil) == (input.process.Error == nil) == (expect.process.Error == nil)

	if !ok {
		t.Error("Unexpected error state")
	}

	ok = ok && input.sameRegisters(t, expect)
	ok = ok && input.sameStack(t, expect)
	ok = ok && input.sameOutput(t, expect)

	if !ok {
		t.Fail()
	}
}
