// Copyright © 2017 Johnny Morrice <john@functorama.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/johnny-morrice/shapes"
	"github.com/johnny-morrice/shapes/brainfuck"
)

// brainfuckCmd represents the brainfuck command
var brainfuckCmd = &cobra.Command{
	Use:     "brainfuck",
	Short:   "Brainfuck interpreter",
	Example: "shapes brainfuck --file prog." + __BRAINFUCK_EXTENSION,
	Run:     runBrainfuck,
}

func runBrainfuck(cmd *cobra.Command, args []string) {
	source := getBrainfuckSource(cmd)

	ast, err := brainfuck.Parse(source)

	err = shapes.InterpretProgramAST(ast, os.Stdin, os.Stdout)

	if err != nil {
		die(err)
	}
}

var sourceFile string
var expression string

func init() {
	RootCmd.AddCommand(brainfuckCmd)

	brainfuckCmd.Flags().StringVar(&sourceFile, __BRAINFUCK_FILE_PARAM, __BRAINFUCK_FILE_DEFAULT, __BRAINFUCK_FILE_USAGE)
	brainfuckCmd.Flags().StringVar(&expression, __BRAINFUCK_EXPRESSION_PARAM, __BRAINFUCK_EXPRESSION_DEFAULT, __BRAINFUCK_EXPRESSION_USAGE)
}

func getBrainfuckSource(cmd *cobra.Command) []byte {
	if expression != "" {
		return []byte(expression)
	}

	if sourceFile == "" {
		dieHelp(cmd)
	}

	source, err := ioutil.ReadFile(sourceFile)

	if err != nil {
		die(err)
	}

	return source
}

const __BRAINFUCK_EXTENSION = "bf"
const __BRAINFUCK_FILE_PARAM = "file"
const __BRAINFUCK_FILE_USAGE = "Brainfuck source code file"
const __BRAINFUCK_FILE_DEFAULT = ""
const __BRAINFUCK_EXPRESSION_PARAM = "expression"
const __BRAINFUCK_EXPRESSION_USAGE = "Brainfuck source code"
const __BRAINFUCK_EXPRESSION_DEFAULT = ""
