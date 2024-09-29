package cmd

import (
	"fmt"
	"io"
	"laxc/pkg/concrete"
	"laxc/pkg/lex"
	"os"
	"regexp"
	"slices"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var dumpConfig = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
}

var stages = []string{"lexer", "concrete", "abstract", "attributed", "intermediate", "target"}

func dump(something any) {

	raw := dumpConfig.Sdump(something)

	s0 := regexp.MustCompile(
		`(` +
			`concrete\.` +
			`|abstract\.` +
			`|attributed\.` +
			`|laxc/pkg/` +
			`|laxc/internal/` +
			`|shared\.` +
			`|env\.` +
			`|RefType` +
			`|IntegerType` +
			`|BooleanType` +
			`|RealType` +
			`)[ ]*`,
	)

	s1 := regexp.MustCompile(
		`(` +
			`\(` +
			`|\)` +
			`)`,
	)

	s2 := regexp.MustCompile(
		`(` +
			`len=\d+` +
			`|string ` +
			`|bool ` +
			`|int32 ` +
			`)[ ]*`,
	)

	fmt.Println(
		string(
			// s3.ReplaceAll(
			s2.ReplaceAll(
				s1.ReplaceAll(
					s0.ReplaceAll(
						[]byte(raw),
						nil,
					),
					nil,
				),
				nil,
			),
			// nil,
			// ),
		),
	)
}

var rootCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "laxc <srcfile>",
	Short: "Compile LAX source files into MIPS32 assembly.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no input file given")
		}

		timing, err := cmd.Flags().GetBool("timing")
		cobra.CheckErr(err)

		stage, err := cmd.Flags().GetString("stage")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		optimize, err := cmd.Flags().GetBool("optimize")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if !slices.Contains(stages, stage) && stage != "" {
			cobra.CheckErr(fmt.Errorf("unknown compiler stage %s, permitted are %v", stage, stages))
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

		if stage == "lexer" {
			lexed, err := lex.Lexer.LexString(inputFileName, string(src))
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			for {
				token, err := lexed.Next()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				if token.EOF() {
					break
				}

				fmt.Println(token.GoString())
			}

			return nil
		}

		start := time.Now()

		concreteProg, err := concrete.Parse(inputFileName, string(src))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if timing {
			fmt.Printf("parsing took %s\n", time.Since(start))
		}

		if stage == "concrete" {
			dump(concreteProg)

			return nil
		}

		start = time.Now()

		abstractProg, err := concreteProg.AbstractExpression()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if timing {
			fmt.Printf("concrete => abstract took %s\n", time.Since(start))
		}

		if stage == "abstract" {
			dump(abstractProg)

			return nil
		}

		start = time.Now()

		attributedProg, err := abstractProg.AttributedProgram()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if timing {
			fmt.Printf("abstract => attributed took %s\n", time.Since(start))
		}

		if stage == "attributed" {
			dump(attributedProg)

			return nil
		}

		start = time.Now()

		intermediateProg := attributedProg.IntermediateProgram()

		if timing {
			fmt.Printf("il generation took %s\n", time.Since(start))
		}

		lascotFriendlyIlString := intermediateProg.LascotFriendlyString()

		if optimize {
			intermediateProg.Optimize()
		}

		if stage == "intermediate" {
			fmt.Println(intermediateProg)

			return nil
		}

		targetProgram := intermediateProg.Mips32Program()
		targetString := targetProgram.String()
		targetFileExtension := ".s"

		if stage == "target" {
			fmt.Println(targetString)

			return nil
		}

		ilOutFile, err := os.Create(inputFileName + ".il")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		_, err = ilOutFile.WriteString(lascotFriendlyIlString)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		mips32OutFile, err := os.Create(inputFileName + targetFileExtension)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		_, err = mips32OutFile.WriteString(targetString)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("stage", "s", "", "if set, the compiler will stop with the specified stage and print the result to stdout. "+
		"Possible values are "+fmt.Sprint(stages))

	rootCmd.Flags().Bool("timing", false, "outputs timing information about the different compilation steps")

	rootCmd.Flags().BoolP("optimize", "o", false, "if set, the compiler will perform optimizations "+
		"(constant folding, dead code elimination) using a single static assignment form")
}
