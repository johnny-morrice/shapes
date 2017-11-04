package integration

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/johnny-morrice/shapes"
	"github.com/johnny-morrice/shapes/asm"
	"github.com/johnny-morrice/shapes/brainfuck"
)

type parseFunc func(source []byte) (*asm.AST, error)

type integrationTest struct {
	parseFunc      parseFunc
	source         []byte
	parseOk        bool
	input          []byte
	expectedOutput []byte
}

func TestBrainfuck(t *testing.T) {
	testCases := []integrationTest{
		integrationTest{
			parseFunc:      brainfuck.Parse,
			source:         []byte("++++++++++."),
			parseOk:        true,
			expectedOutput: []byte{10},
		},
		integrationTest{
			parseFunc:      brainfuck.Parse,
			source:         []byte("++++++++++-----."),
			parseOk:        true,
			expectedOutput: []byte{5},
		},
		integrationTest{
			parseFunc:      brainfuck.Parse,
			source:         []byte(">++++>++++++>+++++++<<<.>.>."),
			parseOk:        true,
			expectedOutput: []byte{4, 6, 7},
		},
		integrationTest{
			parseFunc:      brainfuck.Parse,
			source:         []byte("++++++++++[[.--]]."),
			parseOk:        true,
			expectedOutput: []byte{10, 8, 6, 4, 2, 0},
		},
		integrationTest{
			parseFunc: brainfuck.Parse,
			source:    []byte("[[[]]"),
			parseOk:   false,
		},
		integrationTest{
			parseFunc: brainfuck.Parse,
			source:    []byte("[]]"),
			parseOk:   false,
		},
	}

	for i, test := range testCases {
		passed := integrationTestHelper(t, test)

		if !passed {
			t.Errorf("Test case %d failed", i)
		}
	}
}

func integrationTestHelper(t *testing.T, test integrationTest) bool {
	t.Helper()

	ast, err := test.parseFunc(test.source)

	if err == nil != test.parseOk {
		if err == nil {
			t.Error("Expected error")
		} else {
			t.Error(err.Error())
		}
		return false
	}

	inputBuff := &bytes.Buffer{}
	_, err = inputBuff.Write(test.input)

	if err != nil {
		t.Error("Failed to write test input buffer")
		t.FailNow()
		return false
	}

	outputBuff := &bytes.Buffer{}
	shapes.InterpretProgramAST(ast, inputBuff, outputBuff)

	actualOutput := outputBuff.Bytes()

	if !reflect.DeepEqual(test.expectedOutput, actualOutput) {
		t.Errorf("Expected %v but received %v", test.expectedOutput, actualOutput)
		return false
	}

	return true
}
