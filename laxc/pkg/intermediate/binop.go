package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

var _ TwoArgWriteInstruction = &binOp{}

type binOp struct {
	name     string
	result   shared.SymReg
	arg1     shared.SymReg
	arg2     shared.SymReg
	mips32   func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program)
	bytecode func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program)
	optimize func(instr binOp, arg1, arg2 Instruction) (Instruction, bool)
}

func (instr binOp) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\t%s\tt[%d]\tt[%d]\n", line, instr.name, mapping[instr.arg1], mapping[instr.arg2])
}

func (instr binOp) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\ts%d := %s(s%d,s%d);", instr.result, instr.name, instr.arg1, instr.arg2),
	))

	return int64(count), err
}

func (instr binOp) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	arg1 := allocations[instr.arg1].Reg
	arg2 := allocations[instr.arg2].Reg
	result := allocations[instr.result].Reg

	if alloc := allocations[instr.arg1]; alloc.IsSpilled {
		arg1 = mips32.SpillRegs[0]
		mips32Prog.LW(arg1, mips32.RegSp, int16(alloc.MemLoc), "")
	}

	if alloc := allocations[instr.arg2]; alloc.IsSpilled {
		arg2 = mips32.SpillRegs[1]
		mips32Prog.LW(arg2, mips32.RegSp, int16(alloc.MemLoc), "")
	}

	if alloc := allocations[instr.result]; alloc.IsSpilled {
		result = mips32.SpillRegs[0]
	}

	instr.mips32(instr, arg1, arg2, result, mips32Prog)

	if alloc := allocations[instr.result]; alloc.IsSpilled {
		mips32Prog.SW(result, mips32.RegSp, int16(alloc.MemLoc), "")
	}
}

func (instr binOp) Bytecode(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, bytecodeProg *bytecode.Program) {
	instr.bytecode(
		instr,
		allocations[instr.arg1].Reg,
		allocations[instr.arg2].Reg,
		allocations[instr.result].Reg,
		bytecodeProg,
	)
}

func (instr binOp) Optimize(dependencies map[shared.SymReg]Instruction) (Instruction, bool) {
	if instr.optimize != nil {
		arg1 := dependencies[instr.arg1]
		arg2 := dependencies[instr.arg2]

		return instr.optimize(instr, arg1, arg2)
	}

	return &instr, false
}

func (instr binOp) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{instr.arg1, instr.arg2}
}

func (instr binOp) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr binOp) Result() shared.SymReg {
	return instr.result
}

func (instr *binOp) SetResult(result shared.SymReg) {
	instr.result = result
}

func (instr binOp) Arg1() shared.SymReg {
	return instr.arg1
}

func (instr *binOp) SetArg1(arg1 shared.SymReg) {
	instr.arg1 = arg1
}

func (instr binOp) Arg2() shared.SymReg {
	return instr.arg2
}

func (instr *binOp) SetArg2(arg2 shared.SymReg) {
	instr.arg2 = arg2
}

func (*binOp) HasSideEffects() bool {
	return false
}
