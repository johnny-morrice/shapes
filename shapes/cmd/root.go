// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/johnny-morrice/shapes"
)

var cfgFile string
var sourceFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "shapes",
	Short:   "Programming language based on common shapes",
	Example: "shapes file." + __SHAPE_EXTENSION + " OPTIONS",
	PreRun: func(cmd *cobra.Command, args []string) {
		const extension = "." + __SHAPE_EXTENSION

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

		bs, err := ioutil.ReadFile(sourceFile)

		if err != nil {
			die(err)
		}

		source := string(bs)
		err = shapes.ExecuteProgramCode(source, os.Stdin, os.Stdout)

		if err != nil {
			die(err)
		}
	},
}

func dieHelp(cmd *cobra.Command) {
	err := cmd.Usage()

	if err != nil {
		die(err)
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(__EXIT_FAILURE)
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		die(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shapes.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			die(err)
		}

		// Search config in home directory with name ".shapes" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".shapes")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.ReadInConfig()
}

const __EXIT_FAILURE = 1
const __SHAPE_EXTENSION = "shape"
