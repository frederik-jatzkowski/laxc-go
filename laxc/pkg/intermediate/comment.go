package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/mips32"
)

type comment struct {
	text string
}

func (block *BasicBlock) Nop() {
	block.instructions = append(block.instructions, &comment{
		text: "Nop",
	})
}

func (block *BasicBlock) Comment(text string) {
	block.instructions = append(block.instructions, &comment{
		text: text,
	})
}

func (instr comment) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return ""
}

func (instr comment) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\n\t\t// %s", instr.text),
	))

	return int64(count), err
}

func (instr comment) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
}

func (instr comment) Optimize(_ map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr comment) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{}
}

func (instr comment) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (*comment) HasSideEffects() bool {
	return true
}
