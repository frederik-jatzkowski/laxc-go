package attributed

import (
	"fmt"
	"laxc/internal/env"
	"laxc/internal/shared"
	"laxc/pkg/intermediate"
)

type CaseClause struct {
	Condition  Expression
	Cases      []Case
	Else       []Expression
	ResultType env.Type
}

type Case struct {
	Label Expression
	Stats []Expression
}

func (node CaseClause) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	exitBlock := ilFunc.AddBasicBlock("case exit")

	conditionSymReg := ilFunc.NextSymReg()
	node.Condition.IntermediateExpression(ilProg, ilFunc, ilBlock, conditionSymReg)

	nextTestBlock := ilBlock
	lastResultSymReg := shared.SymReg(-1)
	for _, c := range node.Cases {
		labelSymReg := ilFunc.NextSymReg()
		currTestBlock := c.Label.IntermediateExpression(ilProg, ilFunc, nextTestBlock, labelSymReg)

		testSymReg := ilFunc.NextSymReg()
		currTestBlock.IntEquals(testSymReg, conditionSymReg, labelSymReg)

		caseBlock := ilFunc.AddBasicBlock(fmt.Sprintf("case %d", c.Label.(IntegerLiteral).Value))
		nextTestBlock = ilFunc.AddBasicBlock("")

		currTestBlock.BranchIf(testSymReg, caseBlock)
		currTestBlock.Jump(nextTestBlock)

		caseResultSymReg := shared.SymReg(-1)
		for _, stat := range c.Stats {
			caseResultSymReg = ilFunc.NextSymReg()
			caseBlock = stat.IntermediateExpression(ilProg, ilFunc, caseBlock, caseResultSymReg)
		}

		currResultSymReg := caseResultSymReg
		if lastResultSymReg >= 0 {
			currResultSymReg = ilFunc.NextSymReg()
			caseBlock.Phi(currResultSymReg, lastResultSymReg, caseResultSymReg)
		}

		lastResultSymReg = currResultSymReg

		caseBlock.Jump(exitBlock)
	}

	elseBlock := nextTestBlock
	elseResultSymReg := shared.SymReg(-1)
	for _, stat := range node.Else {
		elseResultSymReg = ilFunc.NextSymReg()
		elseBlock = stat.IntermediateExpression(ilProg, ilFunc, elseBlock, elseResultSymReg)
	}

	elseBlock.Jump(exitBlock)

	exitBlock.Phi(result, lastResultSymReg, elseResultSymReg)

	return exitBlock
}

func (node CaseClause) Type() env.Type {
	return node.ResultType
}
