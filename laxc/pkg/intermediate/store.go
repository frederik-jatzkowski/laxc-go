package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/mips32"
)

type store struct {
	addr shared.SymReg
	arg  shared.SymReg
}

func (block *BasicBlock) Store(targetAddr, arg shared.SymReg) {
	block.instructions = append(block.instructions, &store{
		addr: targetAddr,
		arg:  arg,
	})
}

func (instr store) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\tStore\tt[%d]\tt[%d]\n", line, mapping[instr.arg], mapping[instr.addr])
}

func (instr store) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\t*s%d := s%d;", instr.addr, instr.arg),
	))

	return int64(count), err
}

func (instr store) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	arg := allocations[instr.arg].Reg
	target := allocations[instr.addr].Reg
	if alloc := allocations[instr.arg]; alloc.IsSpilled {
		arg = mips32.SpillRegs[0]
	}

	if alloc := allocations[instr.addr]; alloc.IsSpilled {
		target = mips32.SpillRegs[0]
	}

	mips32Prog.SW(arg, target, 0, "")
}

func (instr store) Optimize(_ map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr store) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{instr.addr, instr.arg}
}

func (instr store) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr store) Arg1() shared.SymReg {
	return instr.addr
}

func (instr *store) SetArg1(addr shared.SymReg) {
	instr.addr = addr
}

func (instr store) Arg2() shared.SymReg {
	return instr.arg
}

func (instr *store) SetArg2(arg shared.SymReg) {
	instr.arg = arg
}

func (*store) HasSideEffects() bool {
	return true
}
