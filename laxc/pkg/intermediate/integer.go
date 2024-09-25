package intermediate

import (
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

func (block *BasicBlock) IntAdd(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "IntAdd",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.ADD(result, arg1, arg2, "")
		},
		bytecode: func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.INT_ADD(result, arg1, arg2)
		},
		optimize: func(instr binOp, arg1, arg2 Instruction) (Instruction, bool) {
			{
				lit1, ok1 := arg1.(*literal)
				lit2, ok2 := arg2.(*literal)
				if ok1 && ok2 {
					return &literal{result: instr.result, value: lit1.value + lit2.value}, true
				}
			}

			return &instr, false
		},
	})
}

func (block *BasicBlock) IntSub(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "IntSub",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.SUB(result, arg1, arg2, "")
		},
		bytecode: func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.INT_SUB(result, arg1, arg2)
		},
		optimize: func(instr binOp, arg1, arg2 Instruction) (Instruction, bool) {
			{
				lit1, ok1 := arg1.(*literal)
				lit2, ok2 := arg2.(*literal)
				if ok1 && ok2 {
					return &literal{result: instr.result, value: lit1.value - lit2.value}, true
				}
			}

			return &instr, false
		},
	})
}

func (block *BasicBlock) IntMul(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "IntMul",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			// do multiplication, result is in $HI and $LO
			mips32Prog.MULT(arg1, arg2, "")

			// get $HI
			work1 := mips32.WorkRegs[0]
			mips32Prog.MFHI(work1, "")

			// get $LO
			work2 := mips32.WorkRegs[1]
			mips32Prog.MFLO(work2, "")

			// shift $LO, such that it should always equal $HI, if no overflow occurred
			mips32Prog.SRA(work2, work2, 31, "")

			// trap, if overflow occurred
			mips32Prog.TNE(work1, work2, "")

			// write result to result register
			mips32Prog.MFLO(result, "")
		},
		bytecode: func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.INT_MUL(result, arg1, arg2)
		},
		optimize: func(instr binOp, arg1, arg2 Instruction) (Instruction, bool) {
			{
				lit1, ok1 := arg1.(*literal)
				lit2, ok2 := arg2.(*literal)
				if ok1 && ok2 {
					return &literal{result: instr.result, value: lit1.value * lit2.value}, true
				}
			}

			return &instr, false
		},
	})
}

func (block *BasicBlock) IntDiv(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "IntDiv",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			// trap on division by 0
			mips32Prog.TEQ(arg2, mips32.RegZero, "")

			// do division, result is in $HI and $LO
			mips32Prog.DIV(arg1, arg2, "")

			// write result from $LO to result register
			mips32Prog.MFLO(result, "")
		},
		bytecode: func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.INT_DIV(result, arg1, arg2)
		},
		optimize: func(instr binOp, arg1, arg2 Instruction) (Instruction, bool) {
			{
				lit1, ok1 := arg1.(*literal)
				lit2, ok2 := arg2.(*literal)
				if ok1 && ok2 {
					return &literal{result: instr.result, value: lit1.value / lit2.value}, true
				}
			}

			return &instr, false
		},
	})
}

func (block *BasicBlock) IntMod(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "IntMod",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			// trap on division by 0
			mips32Prog.TEQ(arg2, mips32.RegZero, "")

			// do division, result is in $HI and $LO
			mips32Prog.DIV(arg1, arg2, "")

			// write result from $HI to result register
			mips32Prog.MFHI(result, "")
		},
		bytecode: func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.INT_MOD(result, arg1, arg2)
		},
		optimize: func(instr binOp, arg1, arg2 Instruction) (Instruction, bool) {
			{
				lit1, ok1 := arg1.(*literal)
				lit2, ok2 := arg2.(*literal)
				if ok1 && ok2 {
					return &literal{result: instr.result, value: lit1.value % lit2.value}, true
				}
			}

			return &instr, false
		},
	})
}

func (block *BasicBlock) IntNeg(result, arg shared.SymReg) {
	block.instructions = append(block.instructions, &unOp{
		name:   "IntNeg",
		result: result,
		arg:    arg,
		mips32: func(instr unOp, arg, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.SUB(result, mips32.RegZero, arg, "")
		},
		bytecode: func(instr unOp, arg, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.INT_NEG(result, arg)
		},
		optimize: func(instr unOp, arg Instruction) (Instruction, bool) {
			{
				lit, ok := arg.(*literal)
				if ok {
					return &literal{result: instr.result, value: 0 - lit.value}, true
				}
			}

			return &instr, false
		},
	})
}
