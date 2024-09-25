package cmd

import (
	"fmt"
	"io"
	"laxc/pkg/target/bytecode"
	"os"

	"github.com/spf13/cobra"
)

var disassembleCmd = &cobra.Command{
	Args: cobra.MaximumNArgs(1),
	Use:  "disassemble <srcfile>",
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

		targetProgram, err := bytecode.NewProgramFromBytes(src)
		if err != nil {
			fmt.Println(fmt.Errorf("could not read bytecode program from binary data: %w", err))
			os.Exit(1)
		}

		fmt.Println(targetProgram.Disassemble())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(disassembleCmd)
}
