package cmd

import (
	"fmt"
	"io"
	"laxc/pkg/concrete"
	"laxc/pkg/target/bytecode"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "run <srcfile>",
	Short: "Compile LAX source files into bytecode and execute it.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no input file given")
		}

		inputFileName := args[0]

		inputFile, err := os.Open(inputFileName)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		src, err := io.ReadAll(inputFile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if strings.HasSuffix(inputFileName, ".lxb") {
			targetProgram, err := bytecode.NewProgramFromBytes(src)
			if err != nil {
				fmt.Println(fmt.Errorf("could not read bytecode program from binary data: %w", err))
				os.Exit(1)
			}

			targetProgram.Run()

			return nil
		}

		concreteProg, err := concrete.Parse(inputFileName, string(src))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		abstractProg, err := concreteProg.AbstractExpression()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		attributedProg, err := abstractProg.AttributedProgram()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		intermediateProg := attributedProg.IntermediateProgram()
		targetProgram := intermediateProg.BytecodeProgram()
		targetProgram.Run()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
