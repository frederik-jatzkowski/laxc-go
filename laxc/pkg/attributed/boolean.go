package attributed

import (
	"laxc/internal/env"
	"laxc/internal/shared"
	"laxc/pkg/intermediate"
)

type IntLessThan struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr IntLessThan) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.IntLessThan(result, left, right)

	return ilBlock
}

func (expr IntLessThan) Type() env.Type {
	return env.BooleanType{}
}

type IntGreaterThan struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr IntGreaterThan) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.IntGreaterThan(result, left, right)

	return ilBlock
}

func (expr IntGreaterThan) Type() env.Type {
	return env.BooleanType{}
}

type Equals struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr Equals) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.IntEquals(result, left, right)

	return ilBlock
}

func (expr Equals) Type() env.Type {
	return env.BooleanType{}
}

type RealLessThan struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr RealLessThan) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.FloatLessThan(result, left, right)

	return ilBlock
}

func (expr RealLessThan) Type() env.Type {
	return env.BooleanType{}
}

type RealGreaterThan struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr RealGreaterThan) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.FloatLessThan(result, right, left)

	return ilBlock
}

func (expr RealGreaterThan) Type() env.Type {
	return env.BooleanType{}
}

type RealEquals struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr RealEquals) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.FloatEquals(result, left, right)

	return ilBlock
}

func (expr RealEquals) Type() env.Type {
	return env.BooleanType{}
}

type BoolNeg struct {
	OpSym string
	Arg   Expression
}

func (expr BoolNeg) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	arg := ilFunc.NextSymReg()
	ilBlock = expr.Arg.IntermediateExpression(ilProg, ilFunc, ilBlock, arg)

	ilBlock.BoolNeg(result, arg)

	return ilBlock
}

func (expr BoolNeg) Type() env.Type {
	return env.BooleanType{}
}

type BoolAnd struct {
	Left  Expression
	Right Expression
}

func (expr BoolAnd) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	lhsSymReg := ilFunc.NextSymReg()
	lhsBlock := expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, lhsSymReg)

	lhsSymRegForPhi := ilFunc.NextSymReg()
	lhsBlock.Assign(lhsSymRegForPhi, lhsSymReg)

	rhsBlock := ilFunc.AddBasicBlock("rhs of conjunction")
	exitBlock := ilFunc.AddBasicBlock("exit of conjunction")

	lhsBlock.BranchIfNot(lhsSymReg, exitBlock)
	lhsBlock.Jump(rhsBlock)

	rhsSymReg := ilFunc.NextSymReg()
	rhsBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, rhsBlock, rhsSymReg)

	rhsSymRegForPhi := ilFunc.NextSymReg()
	rhsBlock.BitwiseAnd(rhsSymRegForPhi, lhsSymReg, rhsSymReg)
	rhsBlock.Jump(exitBlock)

	exitBlock.Phi(result, lhsSymRegForPhi, rhsSymRegForPhi)

	return exitBlock
}

func (expr BoolAnd) Type() env.Type {
	return env.BooleanType{}
}

type BoolOr struct {
	Left  Expression
	Right Expression
}

func (expr BoolOr) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	lhsSymReg := ilFunc.NextSymReg()
	lhsBlock := expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, lhsSymReg)

	lhsSymRegForPhi := ilFunc.NextSymReg()
	lhsBlock.Assign(lhsSymRegForPhi, lhsSymReg)

	rhsBlock := ilFunc.AddBasicBlock("rhs of disjunction")
	exitBlock := ilFunc.AddBasicBlock("exit of disjunction")

	lhsBlock.BranchIf(lhsSymReg, exitBlock)
	lhsBlock.Jump(rhsBlock)

	rhsSymReg := ilFunc.NextSymReg()
	rhsBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, rhsBlock, rhsSymReg)

	rhsSymRegForPhi := ilFunc.NextSymReg()
	rhsBlock.BitwiseOr(rhsSymRegForPhi, lhsSymReg, rhsSymReg)
	rhsBlock.Jump(exitBlock)

	exitBlock.Phi(result, lhsSymRegForPhi, rhsSymRegForPhi)

	return exitBlock
}

func (expr BoolOr) Type() env.Type {
	return env.BooleanType{}
}
