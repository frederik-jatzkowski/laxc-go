package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

type jump struct {
	Name       string
	Target     *BasicBlock
	Mips32Func func(instr jump, mips32Prog *mips32.Program)
}

func (block *BasicBlock) Jump(target *BasicBlock) {
	block.instructions = append(block.instructions, &jump{
		Name:   "Jump",
		Target: target,
		Mips32Func: func(instr jump, mips32Prog *mips32.Program) {
			mips32Prog.J(instr.Target.Label(), "")
		},
	})
}

func (instr jump) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\t%s\t\t%s\t\n", line, instr.Name, instr.Target.Label())
}

func (instr jump) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\tgoto %s;", instr.Target.Label()),
	))

	return int64(count), err
}

func (instr jump) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	instr.Mips32Func(instr, mips32Prog)
}

func (instr jump) Bytecode(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, bytecodeProg *bytecode.Program) {
	bytecodeProg.JUMP(instr.Target.Label())
}

func (instr jump) Optimize(_ map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr jump) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{}
}

func (instr jump) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr jump) Result() shared.Option[shared.SymReg] {
	return shared.None[shared.SymReg]()
}

func (*jump) HasSideEffects() bool {
	return true
}
