package integration

import (
	"testing"

	"github.com/johnny-morrice/shapes/brainfuck"
)

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
			source:         []byte(">++++>++++++>+++++++<<<.>.>.>."),
			parseOk:        true,
			expectedOutput: []byte{0, 4, 6, 7},
		},
		integrationTest{
			parseFunc:      brainfuck.Parse,
			source:         []byte("++++++++++[[.--]]."),
			parseOk:        true,
			expectedOutput: []byte{10, 8, 6, 4, 2, 0},
		},
		integrationTest{
			parseFunc:      brainfuck.Parse,
			source:         []byte("+++++[>,.<-]"),
			parseOk:        true,
			input:          []byte("hello"),
			expectedOutput: []byte("hello"),
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
		t.Logf("Running test case %d", i)
		passed := integrationTestHelper(t, test)

		if !passed {
			t.Errorf("Test case %d failed", i)
		}
	}
}
