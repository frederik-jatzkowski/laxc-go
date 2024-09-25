package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

var _ OneArgWriteInstruction = &assign{}

type assign struct {
	result shared.SymReg
	arg    shared.SymReg
}

func (block *BasicBlock) Assign(result, arg shared.SymReg) {
	block.instructions = append(block.instructions, &assign{
		result: result,
		arg:    arg,
	})
}

func (instr assign) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\tAssign\tt[%d]\n", line, mapping[instr.arg])
}

func (instr assign) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\ts%d := s%d;", instr.result, instr.arg),
	))

	return int64(count), err
}

func (instr assign) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	arg := allocations[instr.arg].Reg
	result := allocations[instr.result].Reg

	if alloc := allocations[instr.arg]; alloc.IsSpilled {
		arg = mips32.SpillRegs[0]
		mips32Prog.LW(arg, mips32.RegSp, int16(alloc.MemLoc), "")
	}

	if alloc := allocations[instr.result]; alloc.IsSpilled {
		result = mips32.SpillRegs[0]
	}

	mips32Prog.OR(result, mips32.RegZero, arg, "")

	if alloc := allocations[instr.result]; alloc.IsSpilled {
		mips32Prog.SW(result, mips32.RegSp, int16(alloc.MemLoc), "")
	}
}

func (instr assign) Bytecode(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, bytecodeProg *bytecode.Program) {
	result := allocations[instr.result].Reg
	arg := allocations[instr.arg].Reg

	bytecodeProg.OR(result, arg, arg)
}

func (instr assign) Optimize(_ map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr assign) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{instr.arg}
}

func (instr assign) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr assign) Result() shared.SymReg {
	return instr.result
}

func (instr *assign) SetResult(result shared.SymReg) {
	instr.result = result
}

func (instr assign) Arg1() shared.SymReg {
	return instr.arg
}

func (instr *assign) SetArg1(arg shared.SymReg) {
	instr.arg = arg
}

func (*assign) HasSideEffects() bool {
	return false
}
