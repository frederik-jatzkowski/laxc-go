package attributed

import (
	"laxc/internal/env"
	"laxc/internal/shared"
	"laxc/pkg/intermediate"
	"math"
	"math/big"
)

type RealLiteral struct {
	Value big.Float
}

func (expr RealLiteral) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	float, _ := expr.Value.Float32()
	bits := math.Float32bits(float)
	ilBlock.Literal(result, int32(bits))

	return ilBlock
}

func (expr RealLiteral) Type() env.Type {
	return env.RealType{}
}

type RealAdd struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr RealAdd) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.FloatAdd(result, left, right)

	return ilBlock
}

func (expr RealAdd) Type() env.Type {
	return env.RealType{}
}

type RealSub struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr RealSub) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.FloatSub(result, left, right)

	return ilBlock
}

func (expr RealSub) Type() env.Type {
	return env.RealType{}
}

type RealMul struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr RealMul) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.FloatMul(result, left, right)

	return ilBlock
}

func (expr RealMul) Type() env.Type {
	return env.RealType{}
}

type RealDiv struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr RealDiv) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	left := ilFunc.NextSymReg()
	right := ilFunc.NextSymReg()

	ilBlock = expr.Left.IntermediateExpression(ilProg, ilFunc, ilBlock, left)
	ilBlock = expr.Right.IntermediateExpression(ilProg, ilFunc, ilBlock, right)

	ilBlock.FloatDiv(result, left, right)

	return ilBlock
}

func (expr RealDiv) Type() env.Type {
	return env.RealType{}
}

type RealNeg struct {
	OpSym string
	Arg   Expression
}

func (expr RealNeg) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	arg := ilFunc.NextSymReg()
	ilBlock = expr.Arg.IntermediateExpression(ilProg, ilFunc, ilBlock, arg)
	ilBlock.FloatNeg(result, arg)

	return ilBlock
}

func (expr RealNeg) Type() env.Type {
	return env.RealType{}
}
