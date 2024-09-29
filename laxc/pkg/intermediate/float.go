package intermediate

import (
	"laxc/internal/shared"
	"laxc/pkg/target/mips32"
)

func (block *BasicBlock) IntToFloat(result, arg shared.SymReg) {
	block.instructions = append(block.instructions, &unOp{
		name:   "IntToFloat",
		result: result,
		arg:    arg,
		mips32: func(instr unOp, arg, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.MTC1(arg, mips32.RegF0, "")
			mips32Prog.CVTSW(mips32.RegF0, mips32.RegF0, "")
			mips32Prog.MFC1(result, mips32.RegF0, "")
		},
	})
}

func (block *BasicBlock) FloatAdd(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "FloatAdd",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.MTC1(arg1, mips32.RegF0, "")
			mips32Prog.MTC1(arg2, mips32.RegF1, "")

			mips32Prog.ADDS(mips32.RegF0, mips32.RegF0, mips32.RegF1, "")

			mips32Prog.MFC1(result, mips32.RegF0, "")
		},
	})
}

func (block *BasicBlock) FloatSub(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "FloatSub",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.MTC1(arg1, mips32.RegF0, "")
			mips32Prog.MTC1(arg2, mips32.RegF1, "")

			mips32Prog.SUBS(mips32.RegF0, mips32.RegF0, mips32.RegF1, "")

			mips32Prog.MFC1(result, mips32.RegF0, "")
		},
	})
}

func (block *BasicBlock) FloatMul(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "FloatMul",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.MTC1(arg1, mips32.RegF0, "")
			mips32Prog.MTC1(arg2, mips32.RegF1, "")

			mips32Prog.MULS(mips32.RegF0, mips32.RegF0, mips32.RegF1, "")

			mips32Prog.MFC1(result, mips32.RegF0, "")
		},
	})
}

func (block *BasicBlock) FloatDiv(result, arg1, arg2 shared.SymReg) {
	block.instructions = append(block.instructions, &binOp{
		name:   "FloatDiv",
		result: result,
		arg1:   arg1,
		arg2:   arg2,
		mips32: func(instr binOp, arg1, arg2, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.MTC1(arg1, mips32.RegF0, "")
			mips32Prog.MTC1(arg2, mips32.RegF1, "")

			mips32Prog.DIVS(mips32.RegF0, mips32.RegF0, mips32.RegF1, "")

			mips32Prog.MFC1(result, mips32.RegF0, "")
		},
	})
}

func (block *BasicBlock) FloatNeg(result, arg shared.SymReg) {
	block.instructions = append(block.instructions, &unOp{
		name:   "FloatNeg",
		result: result,
		arg:    arg,
		mips32: func(instr unOp, arg, result shared.Reg, mips32Prog *mips32.Program) {
			mips32Prog.MTC1(arg, mips32.RegF0, "")

			mips32Prog.NEGS(mips32.RegF0, mips32.RegF0, "")

			mips32Prog.MFC1(result, mips32.RegF0, "")
		},
	})
}
