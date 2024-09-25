package intermediate

import (
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

type Instruction interface {
	io.WriterTo
	LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string
	Mips32(
		symRegAllocs map[shared.SymReg]Allocation,
		localSymVarAllocs map[shared.LocalSymVar]int32,
		mips32Prog *mips32.Program,
	)
	Bytecode(
		symRegAllocs map[shared.SymReg]Allocation,
		localSymVarAllocs map[shared.LocalSymVar]int32,
		bytecodeProg *bytecode.Program,
	)
	Optimize(dependencies map[shared.SymReg]Instruction) (optimized Instruction, propagate bool)
	UsedSymRegs() []shared.SymReg
	UsedLocalSymVars() []shared.LocalSymVar
	HasSideEffects() bool
}

type OneArgInstruction interface {
	Instruction
	Arg1() shared.SymReg
	SetArg1(shared.SymReg)
}

type TwoArgInstruction interface {
	Instruction
	OneArgInstruction
	Arg2() shared.SymReg
	SetArg2(shared.SymReg)
}

type WriteInstruction interface {
	Instruction
	Result() shared.SymReg
	SetResult(shared.SymReg)
}

type OneArgWriteInstruction interface {
	Instruction
	WriteInstruction
	OneArgInstruction
}

type TwoArgWriteInstruction interface {
	Instruction
	WriteInstruction
	TwoArgInstruction
}
