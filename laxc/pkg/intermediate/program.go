package intermediate

import (
	"bytes"
	"fmt"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
	"text/tabwriter"
)

type Program struct {
	main              *Function
	functions         []Function
	instructions      []Instruction
	symRegAllocs      map[shared.SymReg]Allocation
	localSymVarAllocs map[shared.LocalSymVar]int32
}

func NewProg() *Program {
	program := Program{
		instructions: make([]Instruction, 0),
	}

	program.main = program.AddFunction()

	return &program
}

func (prog *Program) Main() *Function {
	return prog.main
}

func (prog *Program) AddFunction() *Function {
	function := Function{
		Id:      len(prog.functions),
		program: prog,
	}

	function.entry = function.AddBasicBlock("entry point")

	prog.functions = append(prog.functions, function)

	return &prog.functions[len(prog.functions)-1]
}

func (prog *Program) Append(instr Instruction) (symReg shared.SymReg) {
	prog.instructions = append(prog.instructions, instr)

	return shared.SymReg(len(prog.instructions)) - 1
}

func (prog Program) String() string {
	buf := bytes.NewBufferString("")
	writer := tabwriter.NewWriter(buf, 4, 2, 2, ' ', 0)

	for _, function := range prog.functions {
		function.WriteTo(writer)
	}

	writer.Flush()

	return buf.String()
}

func (prog Program) LascotFriendlyString() string {
	for _, function := range prog.functions {
		function.sortBasicBlocksTopologically()
	}

	buf := bytes.NewBufferString("")
	writer := tabwriter.NewWriter(buf, 0, 2, 2, ' ', 0)

	writer.Write([]byte("main:\n"))

	mapping := make(map[shared.SymReg]shared.SymReg)
	{
		line := shared.SymReg(0)
		for _, function := range prog.functions {
			for _, block := range function.blocks {
				for _, instruction := range block.instructions {
					switch instruction.(type) {
					case *comment:
						continue
					}

					write, ok := instruction.(WriteInstruction)
					if !ok {
						line++
						continue
					}

					mapping[write.Result()] = line
					line++
				}
			}
		}
	}

	{
		line := 0
		for _, function := range prog.functions {
			for _, block := range function.blocks {
				writer.Write([]byte(fmt.Sprintf("%s:\n", block.Label())))
				for _, instruction := range block.instructions {
					switch instruction.(type) {
					case *comment:
						continue
					}

					str := instruction.LascotFriendlyString(line, mapping)
					writer.Write([]byte(str))
					line++
				}
			}
		}
	}

	writer.Flush()

	return buf.String()
}

func (prog *Program) Mips32Program() mips32.Program {
	mips32Prog := mips32.NewProgram()

	mips32Prog.LABEL("main", "")

	mips32Prog.ADD(mips32.RegFp, mips32.RegZero, mips32.RegSp, "")
	mips32Prog.ADDI(mips32.RegSp, mips32.RegSp, -8, "")

	for _, function := range prog.functions {
		function.Mips32Program(mips32Prog)
	}

	return *mips32Prog
}

func (prog *Program) BytecodeProgram() bytecode.Program {
	result := bytecode.NewEmptyProgram()

	for _, instr := range prog.instructions {
		instr.Bytecode(prog.symRegAllocs, prog.localSymVarAllocs, &result)
	}

	return result
}

func (prog *Program) Optimize() {
	for _, function := range prog.functions {
		function.Optimize()
	}
}
