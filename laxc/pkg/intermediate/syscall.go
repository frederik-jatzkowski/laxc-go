package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/mips32"
)

var _ OneArgInstruction = &syscall{}

type syscall struct {
	name   string
	arg    shared.SymReg
	mips32 func(instr syscall, arg shared.Reg, mips32Prog *mips32.Program)
}

func (block *BasicBlock) PrintInt(arg shared.SymReg) {
	block.instructions = append(block.instructions, &syscall{
		name: "PrintInt",
		arg:  arg,
		mips32: func(instr syscall, arg shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.OR(mips32.RegA0, mips32.RegZero, arg, "")
			mips32Prog.ORI(mips32.RegV0, mips32.RegZero, 1, "")
			mips32Prog.SYSCALL("")
		},
	})
}

func (block *BasicBlock) PrintFloat(arg shared.SymReg) {
	block.instructions = append(block.instructions, &syscall{
		name: "PrintFloat",
		arg:  arg,
		mips32: func(instr syscall, arg shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.MTC1(arg, mips32.RegF12, "")
			mips32Prog.ORI(mips32.RegV0, mips32.RegZero, 2, "")
			mips32Prog.SYSCALL("")
		},
	})
}

func (block *BasicBlock) Exit(arg shared.SymReg) {
	block.instructions = append(block.instructions, &syscall{
		name: "Exit",
		arg:  arg,
		mips32: func(instr syscall, arg shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.OR(mips32.RegV0, mips32.RegZero, arg, "exit status 0")
			mips32Prog.SYSCALL("")
		},
	})
}

func (instr syscall) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\t%s\tt[%d]\n", line, instr.name, mapping[instr.arg])
}

func (instr syscall) WriteTo(writer io.Writer) (int64, error) {

	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\t%s(s%d);", instr.name, instr.arg),
	))

	return int64(count), err
}

func (instr syscall) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	arg := allocations[instr.arg].Reg
	if alloc := allocations[instr.arg]; alloc.IsSpilled {
		arg = mips32.SpillRegs[0]
		mips32Prog.LW(arg, mips32.RegSp, int16(alloc.MemLoc), "")
	}

	instr.mips32(instr, arg, mips32Prog)
}

func (instr syscall) Optimize(_ map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr syscall) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{instr.arg}
}

func (instr syscall) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr syscall) Arg1() shared.SymReg {
	return instr.arg
}

func (instr *syscall) SetArg1(arg shared.SymReg) {
	instr.arg = arg
}

func (*syscall) HasSideEffects() bool {
	return true
}
