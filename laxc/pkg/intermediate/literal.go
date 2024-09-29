package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/shared"
	"laxc/pkg/target/mips32"
)

var _ WriteInstruction = &literal{}

type literal struct {
	result shared.SymReg
	value  int32
	mips32 func(instr literal, result shared.Reg, mips32Prog *mips32.Program)
}

func (block *BasicBlock) Literal(result shared.SymReg, value int32) {
	block.instructions = append(block.instructions, &literal{
		result: result,
		value:  value,
		mips32: func(instr literal, result shared.Reg, mips32Prog *mips32.Program) {
			upper := int32(uint32(instr.value) >> 16)
			lower := 0b00000000_00000000_11111111_11111111 & instr.value

			if upper != 0 {
				mips32Prog.LUI(result, upper, "")
				mips32Prog.ORI(result, result, lower, "")
			} else {
				mips32Prog.ORI(result, mips32.RegZero, instr.value, "")
			}
		},
	})
}

func (instr literal) LascotFriendlyString(line int, mapping map[shared.SymReg]shared.SymReg) string {
	return fmt.Sprintf("\tt[%d]:\tLit\t%d\n", line, instr.value)
}

func (instr literal) WriteTo(writer io.Writer) (int64, error) {
	count, err := writer.Write([]byte(
		fmt.Sprintf("\t\ts%d := %d;", instr.result, instr.value),
	))

	return int64(count), err
}

func (instr literal) Mips32(allocations map[shared.SymReg]Allocation, localSymVarAllocs map[shared.LocalSymVar]int32, mips32Prog *mips32.Program) {
	result := allocations[instr.result].Reg
	if alloc := allocations[instr.result]; alloc.IsSpilled {
		result = mips32.SpillRegs[0]
	}

	instr.mips32(instr, result, mips32Prog)

	if alloc := allocations[instr.result]; alloc.IsSpilled {
		mips32Prog.SW(result, mips32.RegSp, int16(alloc.MemLoc), "")
	}
}

func (instr literal) Optimize(_ map[shared.SymReg]Instruction) (Instruction, bool) {
	return &instr, false
}

func (instr literal) UsedSymRegs() []shared.SymReg {
	return []shared.SymReg{instr.result}
}

func (instr literal) UsedLocalSymVars() []shared.LocalSymVar {
	return []shared.LocalSymVar{}
}

func (instr literal) Result() shared.SymReg {
	return instr.result
}

func (instr *literal) SetResult(result shared.SymReg) {
	instr.result = result
}

func (*literal) HasSideEffects() bool {
	return false
}
