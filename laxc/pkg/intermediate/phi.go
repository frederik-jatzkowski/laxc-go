package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

var _ TwoArgWriteInstruction = &phi{}

type phi struct {
	result shared.SymReg
	arg1   shared.SymReg
	arg2   shared.SymReg
}

func (block *BasicBlock) Phi(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &phi{
		result: result,
		arg1:   arg1,
		arg2:   arg2,
	})
}

func (instr phi) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\t%s\tt[%d]\tt[%d]\n", line, "Phi", instr.arg1, instr.arg2)
}

func (instr phi) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\ts%d := Phi(s%d,s%d);", instr.result, instr.arg1, instr.arg2),
	))

	return int64(count), err
}

func (instr phi) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
}

func (instr phi) Bytecode(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, bytecodeProg *bytecode.Program) {
}

func (instr phi) Optimize(dependencies map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr phi) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{instr.arg1, instr.arg2}
}

func (instr phi) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr phi) Result() shared.SymReg {
	return instr.result
}

func (instr *phi) SetResult(result shared.SymReg) {
	instr.result = result
}

func (instr phi) Arg1() shared.SymReg {
	return instr.arg1
}

func (instr *phi) SetArg1(arg1 shared.SymReg) {
	instr.arg1 = arg1
}

func (instr phi) Arg2() shared.SymReg {
	return instr.arg2
}

func (instr *phi) SetArg2(arg2 shared.SymReg) {
	instr.arg2 = arg2
}

func (*phi) HasSideEffects() bool {
	return false
}
