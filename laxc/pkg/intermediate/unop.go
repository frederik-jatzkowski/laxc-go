package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

var _ OneArgWriteInstruction = &unOp{}

type unOp struct {
	name     string
	result   shared.SymReg
	arg      shared.SymReg
	mips32   func(instr unOp, arg, result shared.Reg, mips32Prog *mips32.Program)
	bytecode func(instr unOp, arg, result shared.Reg, bytecodeProg *bytecode.Program)
	optimize func(instr unOp, arg Instruction) (Instruction, bool)
}

func (instr unOp) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\t%s\tt[%d]\n", line, instr.name, mapping[instr.arg])
}

func (instr unOp) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\ts%d := %s(s%d);", instr.result, instr.name, instr.arg),
	))

	return int64(count), err
}

func (instr unOp) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	arg := allocations[instr.arg].Reg
	result := allocations[instr.result].Reg

	if alloc := allocations[instr.arg]; alloc.IsSpilled {
		arg = mips32.SpillRegs[0]
		mips32Prog.LW(arg, mips32.RegSp, int16(alloc.MemLoc), "")
	}

	if alloc := allocations[instr.result]; alloc.IsSpilled {
		result = mips32.SpillRegs[0]
	}

	instr.mips32(instr, arg, result, mips32Prog)

	if alloc := allocations[instr.result]; alloc.IsSpilled {
		mips32Prog.SW(result, mips32.RegSp, int16(alloc.MemLoc), "")
	}
}

func (instr unOp) Bytecode(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, bytecodeProg *bytecode.Program) {
	instr.bytecode(instr, allocations[instr.arg].Reg, allocations[instr.result].Reg, bytecodeProg)
}

func (instr unOp) Optimize(dependencies map[shared.SymReg]Instruction) (Instruction, bool) {
	if instr.optimize != nil {
		arg := dependencies[instr.arg]
		return instr.optimize(instr, arg)
	}

	return &instr, false
}

func (instr unOp) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{instr.arg}
}

func (instr unOp) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr unOp) Result() shared.SymReg {
	return instr.result
}

func (instr *unOp) SetResult(result shared.SymReg) {
	instr.result = result
}

func (instr unOp) Arg1() shared.SymReg {
	return instr.arg
}

func (instr *unOp) SetArg1(arg shared.SymReg) {
	instr.arg = arg
}

func (*unOp) HasSideEffects() bool {
	return false
}
