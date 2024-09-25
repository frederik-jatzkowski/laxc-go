package attributed

import (
	"laxc/internal/env"
	"laxc/internal/shared"
	"laxc/pkg/intermediate"
)

type ReadVariable struct {
	Var *env.Variable
}

func (expr ReadVariable) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	ilBlock.Assign(result, expr.Var.SymReg)

	return ilBlock
}

func (expr ReadVariable) Type() env.Type {
	return expr.Var.Type
}

type Assignment struct {
	Var            *env.Variable
	UnderlyingType env.Type
	Arg            Expression
}

func (expr Assignment) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	arg := ilFunc.NextSymReg()

	ilBlock = expr.Arg.IntermediateExpression(ilProg, ilFunc, ilBlock, arg)
	ilBlock.Store(expr.Var.SymReg, arg)
	ilBlock.Assign(result, expr.Var.SymReg)

	return ilBlock
}

func (expr Assignment) Type() env.Type {
	return expr.Var.Type
}
