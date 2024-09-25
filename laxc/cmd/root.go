package cmd

import (
	"fmt"
	"io"
	"laxc/pkg/concrete"
	"laxc/pkg/lex"
	"os"
	"regexp"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var dumpConfig = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
}

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
			// `|\*\w+` +
			`)`,
		// "",
	)

	s2 := regexp.MustCompile(
		`(` +
			// `(\[\])[A-Za-z\.]+` +
			// `|[A-Za-z\.]*<nil>` +
			`len=\d+` +
			`|string ` +
			`|bool ` +
			`|int32 ` +
			// `|\s*Int: ,` +
			// `|\s*Real:  ` +
			// `\(\*` +
			`)[ ]*`,
		// "",
	)

	// s3 := regexp.MustCompile(
	// 	`(` +
	// 		`\s*Int: ,` +
	// 		// `|\s*Real: ` +
	// 		`)[ ]*`,
	// 	// "",
	// )

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

		switch stage {
		case "lexer", "concrete", "abstract", "attributed", "intermediate", "target", "":
			break
		default:
			cobra.CheckErr(fmt.Errorf("unknown compiler stage: %s", stage))
		}

		inputFileName := args[0]

		go func() {
			time.Sleep(time.Second * 3)
			fmt.Println("i was stuck somewhere")
			os.Exit(1)
		}()

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
			// data, _ := json.MarshalIndent(concreteProg, "", "  ")
			// fmt.Println(string(data))

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
			// data, _ := json.MarshalIndent(abstractProg, "", "  ")
			// fmt.Println(string(data))

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
			// data, _ := json.MarshalIndent(attributedProg, "", "  ")
			// fmt.Println(string(data))

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

		arch, err := cmd.Flags().GetString("arch")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var (
			targetData          string
			targetString        string
			targetFileExtension string
		)
		switch arch {
		case "mips32":
			if stage == "intermediate" {
				intermediateProg.Main().ResolvePhiFunctions()
				fmt.Println(intermediateProg)

				return nil
			}

			targetProgram := intermediateProg.Mips32Program()
			targetData = targetProgram.String()
			targetString = targetData
			targetFileExtension = ".s"
		case "bytecode":
			if stage == "intermediate" {
				fmt.Println(intermediateProg)

				return nil
			}

			targetProgram := intermediateProg.BytecodeProgram()
			targetProgram.Finalize()
			targetData = string(targetProgram.Bytes())
			targetString = targetProgram.Disassemble()
			targetFileExtension = ".lxb"
		default:
			fmt.Printf("unknown target architecture: %s\n", arch)
			os.Exit(1)
		}

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

		_, err = mips32OutFile.WriteString(targetData)
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
		"Possible values are (lexer|concrete|abstract|attributed|intermediate|target)")

	rootCmd.Flags().StringP("arch", "a", "mips32", "sets the target architecture of the compilation (mips32|wasm)")

	rootCmd.Flags().Bool("timing", false, "outputs timing information about the different compilation steps")

	rootCmd.Flags().BoolP("optimize", "o", false, "if set, the compiler will perform optimizations "+
		"(constant folding, dead code elimination) using a single static assignment form")
}
