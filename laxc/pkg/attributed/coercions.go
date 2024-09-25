package attributed

import (
	"laxc/internal/env"
	"laxc/internal/shared"
	"laxc/pkg/intermediate"
)

type Dereference struct {
	Arg Expression
}

func (expr Dereference) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	addr := ilFunc.NextSymReg()
	ilBlock = expr.Arg.IntermediateExpression(ilProg, ilFunc, ilBlock, addr)
	ilBlock.Load(addr, result)

	return ilBlock
}

func (expr Dereference) Type() env.Type {
	argType, err := expr.Arg.Type().Dereference()
	if err != nil {
		panic(err)
	}

	return argType
}

type Widen struct {
	Arg Expression
}

func (expr Widen) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	arg := ilFunc.NextSymReg()
	ilBlock = expr.Arg.IntermediateExpression(ilProg, ilFunc, ilBlock, arg)
	ilBlock.IntToFloat(result, arg)

	return ilBlock
}

func (expr Widen) Type() env.Type {
	return env.RealType{}
}
