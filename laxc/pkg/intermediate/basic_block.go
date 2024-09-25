package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/mips32"
)

type BasicBlock struct {
	id           int
	comment      string
	function     *Function
	instructions []Instruction
}

func (block BasicBlock) Id() int {
	return block.id
}

func (block BasicBlock) Label() string {
	return fmt.Sprintf("f%d_b%d", block.function.Id, block.id)
}

func (block *BasicBlock) WriteTo(writer io.Writer) {
	if block.comment == "" {
		writer.Write([]byte(fmt.Sprintf("\n\t%s:\n", block.Label())))
	} else {
		writer.Write([]byte(fmt.Sprintf("\n\t%s: // %s\n", block.Label(), block.comment)))
	}

	for _, instruction := range block.instructions {
		instruction.WriteTo(writer)
		writer.Write([]byte("\n"))
	}
}

func (block *BasicBlock) Mips32Program(
	localSymVarAllocs map[shared.LocalSymVar]int32,
	symRegAllocs map[shared.SymReg]Allocation,
	mips32Prog *mips32.Program,
) {
	mips32Prog.LABEL(block.Label(), "")
	for _, instr := range block.instructions {
		instr.Mips32(symRegAllocs, localSymVarAllocs, mips32Prog)
	}
}
