package integration

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/johnny-morrice/shapes"
	"github.com/johnny-morrice/shapes/asm"
)

type parseFunc func(source []byte) (*asm.AST, error)

type integrationTest struct {
	parseFunc      parseFunc
	source         []byte
	parseOk        bool
	input          []byte
	expectedOutput []byte
}

func integrationTestHelper(t *testing.T, test integrationTest) bool {
	t.Helper()
	t.Logf("Testing source: %s", string(test.source))

	ast, err := test.parseFunc(test.source)

	if err == nil != test.parseOk {
		t.Log("Unexpected parse error state")
		if err == nil {
			t.Error("Expected error")
		} else {
			t.Error(err.Error())
		}
		return false
	} else if err != nil {
		return true
	}

	inputBuff := &bytes.Buffer{}
	_, err = inputBuff.Write(test.input)

	if err != nil {
		t.Error("Failed to write test input buffer")
		t.FailNow()
		return false
	}

	outputBuff := &bytes.Buffer{}
	err = shapes.InterpretProgramAST(ast, inputBuff, outputBuff)

	if err != nil {
		t.Errorf("Runtime crash: %s", err.Error())
		return false
	}

	actualOutput := outputBuff.Bytes()

	if !reflect.DeepEqual(test.expectedOutput, actualOutput) {
		t.Errorf("Expected %v but received %v", test.expectedOutput, actualOutput)
		return false
	}

	return true
}
