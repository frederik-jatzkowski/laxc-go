package attributed

import (
	"laxc/internal/env"
	"laxc/internal/shared"
	"laxc/pkg/intermediate"
)

type IntAdd struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr IntAdd) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.IntAdd(result, left, right)

	return ilBlock
}

func (expr IntAdd) Type() env.Type {
	return env.IntegerType{}
}

type IntSub struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr IntSub) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.IntSub(result, left, right)

	return ilBlock
}

func (expr IntSub) Type() env.Type {
	return env.IntegerType{}
}

type IntMul struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr IntMul) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.IntMul(result, left, right)

	return ilBlock
}

func (expr IntMul) Type() env.Type {
	return env.IntegerType{}
}

type IntDiv struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr IntDiv) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.IntDiv(result, left, right)

	return ilBlock
}

func (expr IntDiv) Type() env.Type {
	return env.IntegerType{}
}

type IntMod struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr IntMod) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.IntMod(result, left, right)

	return ilBlock
}

func (expr IntMod) Type() env.Type {
	return env.IntegerType{}
}

type IntegerLiteral struct {
	Value int32
}

func (expr IntegerLiteral) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	ilBlock.Literal(result, expr.Value)

	return ilBlock
}

func (expr IntegerLiteral) Type() env.Type {
	return env.IntegerType{}
}

type IntNeg struct {
	OpSym string
	Arg   Expression
}

func (expr IntNeg) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	arg := ilFunc.NextSymReg()
	ilBlock = expr.Arg.IntermediateExpression(ilProg, ilFunc, ilBlock, arg)
	ilBlock.IntNeg(result, arg)

	return ilBlock
}

func (expr IntNeg) Type() env.Type {
	return env.IntegerType{}
}
