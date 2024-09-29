package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/mips32"
)

var _ WriteInstruction = &takeAddress{}

type takeAddress struct {
	result shared.SymReg
	arg    shared.LocalSymVar
}

func (block *BasicBlock) TakeAddress(result shared.SymReg, arg shared.LocalSymVar) {
	block.instructions = append(block.instructions, &takeAddress{
		result: result,
		arg:    arg,
	})
}

func (instr takeAddress) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\tTakeAddress\t%d\n", line, instr.arg)
}

func (instr takeAddress) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\ts%d := &v%d;", instr.result, instr.arg),
	))

	return int64(count), err
}

func (instr takeAddress) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	result := allocations[instr.result].Reg
	if alloc := allocations[instr.result]; alloc.IsSpilled {
		result = mips32.SpillRegs[0]
	}

	mips32Prog.ADDI(result, mips32.RegSp, localSymVarAllocs[instr.arg], "")

	if alloc := allocations[instr.result]; alloc.IsSpilled {
		mips32Prog.SW(result, mips32.RegSp, int16(alloc.MemLoc), "")
	}
}

func (instr takeAddress) Optimize(_ map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr takeAddress) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{}
}

func (instr takeAddress) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{instr.arg}
}

func (instr takeAddress) Result() shared.SymReg {
	return instr.result
}

func (instr *takeAddress) SetResult(result shared.SymReg) {
	instr.result = result
}

func (*takeAddress) HasSideEffects() bool {
	return false
}
