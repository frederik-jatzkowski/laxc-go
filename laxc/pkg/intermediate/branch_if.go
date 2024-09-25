package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

type branchIf struct {
	arg      shared.SymReg
	target   *BasicBlock
	mips32   func(instr branchIf, condition shared.Reg, mips32Prog *mips32.Program)
	bytecode func(instr branchIf, condition shared.Reg, bytecodeProg *bytecode.Program)
}

func (block *BasicBlock) BranchIf(condition shared.SymReg, target *BasicBlock) {
	block.instructions = append(block.instructions, &branchIf{
		arg:    condition,
		target: target,
		mips32: func(instr branchIf, condition shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.BNE(condition, mips32.RegZero, instr.target.Label(), "")
		},
		bytecode: func(instr branchIf, condition shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.JUMP_IF_NOT(condition, instr.target.Label())
		},
	})
}

func (instr branchIf) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\tBranchIf\tt[%d]\t%s\t\n", line, mapping[instr.arg], instr.target.Label())
}

func (instr branchIf) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\tif s%d goto %s;", instr.arg, instr.target.Label()),
	))

	return int64(count), err
}

func (instr branchIf) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	condition := allocations[instr.arg].Reg
	if alloc := allocations[instr.arg]; alloc.IsSpilled {
		condition = mips32.SpillRegs[0]
		mips32Prog.LW(condition, mips32.RegSp, int16(alloc.MemLoc), "")
	}

	instr.mips32(instr, condition, mips32Prog)
}

func (instr branchIf) Bytecode(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, bytecodeProg *bytecode.Program) {
	instr.bytecode(instr, allocations[instr.arg].Reg, bytecodeProg)
}

func (instr branchIf) Optimize(_ map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr branchIf) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{instr.arg}
}

func (instr branchIf) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr branchIf) Arg1() shared.SymReg {
	return instr.arg
}

func (instr *branchIf) SetArg1(arg shared.SymReg) {
	instr.arg = arg
}

func (*branchIf) HasSideEffects() bool {
	return true
}
