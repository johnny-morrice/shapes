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
	"strings"

	"github.com/spf13/cobra"

	"github.com/johnny-morrice/shapes"
	"github.com/johnny-morrice/shapes/brainfuck"
)

// brainfuckCmd represents the brainfuck command
var brainfuckCmd = &cobra.Command{
	Use:     "brainfuck",
	Short:   "Brainfuck interpreter",
	Example: "shapes brainfuck prog." + __BRAINFUCK_EXTENSION,
	PreRun: func(cmd *cobra.Command, args []string) {
		const extension = "." + __BRAINFUCK_EXTENSION

		if len(args) == 0 {
			dieHelp(cmd)
		}

		if strings.HasSuffix(args[0], extension) {
			sourceFile = args[0]
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		if sourceFile == "" {
			dieHelp(cmd)
		}

		source, err := ioutil.ReadFile(sourceFile)

		if err != nil {
			die(err)
		}

		ast, err := brainfuck.Parse(source)

		err = shapes.InterpretProgramAST(ast, os.Stdin, os.Stdout)

		if err != nil {
			die(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(brainfuckCmd)
}

const __BRAINFUCK_EXTENSION = "bf"
