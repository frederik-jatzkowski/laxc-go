package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

type branchIfNot struct {
	arg      shared.SymReg
	target   *BasicBlock
	mips32   func(instr branchIfNot, condition shared.Reg, mips32Prog *mips32.Program)
	bytecode func(instr branchIfNot, condition shared.Reg, bytecodeProg *bytecode.Program)
}

func (block *BasicBlock) BranchIfNot(condition shared.SymReg, target *BasicBlock) {
	block.instructions = append(block.instructions, &branchIfNot{
		arg:    condition,
		target: target,
		mips32: func(instr branchIfNot, condition shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.BEQ(condition, mips32.RegZero, instr.target.Label(), "")
		},
		bytecode: func(instr branchIfNot, condition shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.JUMP_IF_NOT(condition, instr.target.Label())
		},
	})
}

func (instr branchIfNot) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\tBranchIfNot\tt[%d]\t%s\t\n", line, mapping[instr.arg], instr.target.Label())
}

func (instr branchIfNot) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\tif not s%d goto %s;", instr.arg, instr.target.Label()),
	))

	return int64(count), err
}

func (instr branchIfNot) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	condition := allocations[instr.arg].Reg
	if alloc := allocations[instr.arg]; alloc.IsSpilled {
		condition = mips32.SpillRegs[0]
		mips32Prog.LW(condition, mips32.RegSp, int16(alloc.MemLoc), "")
	}

	instr.mips32(instr, condition, mips32Prog)
}

func (instr branchIfNot) Bytecode(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, bytecodeProg *bytecode.Program) {
	instr.bytecode(instr, allocations[instr.arg].Reg, bytecodeProg)
}

func (instr branchIfNot) Optimize(_ map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr branchIfNot) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{instr.arg}
}

func (instr branchIfNot) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr branchIfNot) Arg1() shared.SymReg {
	return instr.arg
}

func (instr *branchIfNot) SetArg1(arg shared.SymReg) {
	instr.arg = arg
}

func (*branchIfNot) HasSideEffects() bool {
	return true
}
