package attributed

import (
	"laxc/internal/env"
	"laxc/internal/shared"
	"laxc/pkg/intermediate"
)

type OneSidedConditionalClause struct {
	Condition Expression
	Then      []Expression
}

func (node OneSidedConditionalClause) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	ilBlock.Comment("if")

	conditionSymReg := ilFunc.NextSymReg()
	ilBlock = node.Condition.IntermediateExpression(ilProg, ilFunc, ilBlock, conditionSymReg)

	thenBlock := ilFunc.AddBasicBlock("then")
	exitBlock := ilFunc.AddBasicBlock("end")

	ilBlock.BranchIf(conditionSymReg, thenBlock)
	ilBlock.Jump(exitBlock)

	for _, stat := range node.Then {
		thenBlock = stat.IntermediateExpression(ilProg, ilFunc, thenBlock, result)
	}

	thenBlock.Jump(exitBlock)

	return exitBlock
}

func (node OneSidedConditionalClause) Type() env.Type {
	return env.VoidType{}
}

type TwoSidedConditionalClause struct {
	Condition  Expression
	Then       []Expression
	Else       []Expression
	ResultType env.Type
}

func (node TwoSidedConditionalClause) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	ilBlock.Comment("if")

	conditionSymReg := ilFunc.NextSymReg()
	ilBlock = node.Condition.IntermediateExpression(ilProg, ilFunc, ilBlock, conditionSymReg)

	thenBlock := ilFunc.AddBasicBlock("then")
	elseBlock := ilFunc.AddBasicBlock("else")
	exitBlock := ilFunc.AddBasicBlock("end")

	ilBlock.BranchIf(conditionSymReg, thenBlock)
	ilBlock.Jump(elseBlock)

	thenResult := ilFunc.NextSymReg()

	for i, stat := range node.Then {
		if i != len(node.Then)-1 {
			thenBlock = stat.IntermediateExpression(ilProg, ilFunc, thenBlock, ilFunc.NextSymReg())
		} else {
			thenBlock = stat.IntermediateExpression(ilProg, ilFunc, thenBlock, thenResult)
		}
	}

	thenBlock.Jump(exitBlock)

	elseResult := ilFunc.NextSymReg()

	for i, stat := range node.Else {
		if i != len(node.Then)-1 {
			elseBlock = stat.IntermediateExpression(ilProg, ilFunc, elseBlock, ilFunc.NextSymReg())
		} else {
			elseBlock = stat.IntermediateExpression(ilProg, ilFunc, elseBlock, elseResult)
		}
	}

	elseBlock.Jump(exitBlock)
	exitBlock.Phi(result, thenResult, elseResult)

	return exitBlock
}

func (node TwoSidedConditionalClause) Type() env.Type {
	return node.ResultType
}
