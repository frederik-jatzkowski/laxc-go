package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

type load struct {
	addr   shared.SymReg
	result shared.SymReg
}

func (block *BasicBlock) Load(addr, result shared.SymReg) {
	block.instructions = append(block.instructions, &load{
		addr:   addr,
		result: result,
	})
}

func (instr load) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\tLoad\tt[%d]\n", line, mapping[instr.addr])
}

func (instr load) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\ts%d := *s%d;", instr.result, instr.addr),
	))

	return int64(count), err
}

func (instr load) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	addr := allocations[instr.addr].Reg
	result := allocations[instr.result].Reg
	if alloc := allocations[instr.addr]; alloc.IsSpilled {
		addr = mips32.SpillRegs[0]
	}

	if alloc := allocations[instr.result]; alloc.IsSpilled {
		result = mips32.SpillRegs[0]
	}

	mips32Prog.LW(result, addr, 0, "")

	if alloc := allocations[instr.result]; alloc.IsSpilled {
		mips32Prog.SW(result, mips32.RegSp, int16(alloc.MemLoc), "")
	}
}

func (instr load) Bytecode(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, bytecodeProg *bytecode.Program) {
	panic("not implemented")
}

func (instr load) Optimize(dependencies map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr load) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{instr.addr}
}

func (instr load) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr load) Result() shared.SymReg {
	return instr.result
}

func (instr *load) SetResult(result shared.SymReg) {
	instr.result = result
}

func (instr load) Arg1() shared.SymReg {
	return instr.addr
}

func (instr *load) SetArg1(addr shared.SymReg) {
	instr.addr = addr
}

func (*load) HasSideEffects() bool {
	return false
}
