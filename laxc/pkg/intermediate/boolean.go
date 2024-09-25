package intermediate

import (
	"fmt"
	"laxc/internal/shared"
	"laxc/pkg/target/bytecode"
	"laxc/pkg/target/mips32"
)

func (block *BasicBlock) IntLessThan(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "IntLessThan",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.SLT(result, arg1, arg2, "")
		},
		bytecode: func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.INT_LT(result, arg1, arg2)
		},
		optimize: func(instr binOp, arg1, arg2 Instruction) (Instruction, bool) {
			{
				lit1, ok1 := arg1.(*literal)
				lit2, ok2 := arg2.(*literal)
				if ok1 && ok2 {
					substituted := literal{result: instr.result}
					if lit1.value < lit2.value {
						substituted.value = 1
					}

					return &substituted, true
				}
			}

			return &instr, false
		},
	})
}

func (block *BasicBlock) IntGreaterThan(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "IntGreaterThan",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.SLT(result, arg2, arg1, "")
		},
		bytecode: func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.INT_LT(result, arg2, arg1)
		},
		optimize: func(instr binOp, arg1, arg2 Instruction) (Instruction, bool) {
			{
				lit1, ok1 := arg1.(*literal)
				lit2, ok2 := arg2.(*literal)
				if ok1 && ok2 {
					substituted := literal{result: instr.result}
					if lit1.value > lit2.value {
						substituted.value = 1
					}

					return &substituted, true
				}
			}

			return &instr, false
		},
	})
}

func (block *BasicBlock) IntEquals(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "IntEquals",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			work1 := mips32.WorkRegs[0]
			work2 := mips32.WorkRegs[1]

			mips32Prog.SLT(work1, arg2, arg1, "")
			mips32Prog.SLT(work2, arg1, arg2, "")

			mips32Prog.NOR(work1, work1, work2, "")
			mips32Prog.ANDI(result, work1, 1, "")
		},
		bytecode: func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.INT_EQ(result, arg1, arg2)
		},
		optimize: func(instr binOp, arg1, arg2 Instruction) (Instruction, bool) {
			{
				lit1, ok1 := arg1.(*literal)
				lit2, ok2 := arg2.(*literal)
				if ok1 && ok2 {
					substituted := literal{result: instr.result}
					if lit1.value == lit2.value {
						substituted.value = 1
					}

					return &substituted, true
				}
			}

			return &instr, false
		},
	})
}

func (block *BasicBlock) FloatLessThan(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "FloatLessThan",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.MTC1(arg1, mips32.RegF0, "")
			mips32Prog.MTC1(arg2, mips32.RegF1, "")

			label := fmt.Sprintf("float_lt_t%d", mips32Prog.Line())

			mips32Prog.CLTS(mips32.RegF0, mips32.RegF1, "") // compare floats
			mips32Prog.ORI(result, mips32.RegZero, 1, "")   // set result to 1
			mips32Prog.BC1T(label, "")                      // if comparison was true, jump over next ORI
			mips32Prog.ORI(result, mips32.RegZero, 0, "")   // if comparison was false, set result to false
			mips32Prog.LABEL(label, "")                     // label for jump
		},
		bytecode: func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.FLT_LT(result, arg1, arg2)
		},
	})
}

func (block *BasicBlock) FloatEquals(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "FloatEquals",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.MTC1(arg1, mips32.RegF0, "")
			mips32Prog.MTC1(arg2, mips32.RegF1, "")

			label := fmt.Sprintf("float_eq_t%d", mips32Prog.Line())

			mips32Prog.CEQS(mips32.RegF0, mips32.RegF1, "") // compare floats
			mips32Prog.ORI(result, mips32.RegZero, 1, "")   // set result to 1
			mips32Prog.BC1T(label, "")                      // if comparison was true, jump over next ORI
			mips32Prog.ORI(result, mips32.RegZero, 0, "")   // if comparison was false, set result to false
			mips32Prog.LABEL(label, "")                     // label for jump
		},
		bytecode: func(instr binOp, arg1, arg2, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.FLT_EQ(result, arg1, arg2)
		},
	})
}

func (block *BasicBlock) BoolNeg(result, arg shared.SymReg) {
	block.instructions = append(block.instructions, &unOp{
		name:   "BoolNeg",
		result: result,
		arg:    arg,
		mips32: func(instr unOp, arg, result shared.Reg, mips32Prog *mips32.Program) {
			work1 := mips32.WorkRegs[0]

			mips32Prog.ORI(work1, mips32.RegZero, 1, "")
			mips32Prog.XOR(result, work1, arg, "")
		},
		bytecode: func(instr unOp, arg, result shared.Reg, bytecodeProg *bytecode.Program) {
			bytecodeProg.BOOL_NOT(result, arg)
		},
	})
}

func (block *BasicBlock) BitwiseAnd(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "BitwiseAnd",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.AND(result, arg1, arg2, "")
		},
	})
}

func (block *BasicBlock) BitwiseOr(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "BitwiseOr",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.OR(result, arg1, arg2, "")
		},
	})
}
