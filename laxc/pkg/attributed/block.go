package attributed

import (
	"laxc/internal/env"
	"laxc/internal/shared"
	"laxc/pkg/intermediate"
)

type Block struct {
	Decls []Declaration
	Stats []Expression
	Scope env.Table
}

func (node Block) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	ilBlock.Comment("declare")

	for _, decl := range node.Decls {
		ilBlock = decl.IntermediateDeclaration(node.Scope, ilProg, ilFunc, ilBlock)
	}

	ilBlock.Comment("begin")

	for i, stat := range node.Stats {
		if i != len(node.Stats)-1 {
			ilBlock = stat.IntermediateExpression(ilProg, ilFunc, ilBlock, ilFunc.NextSymReg())
		} else {
			ilBlock = stat.IntermediateExpression(ilProg, ilFunc, ilBlock, result)
		}
	}

	ilBlock.Comment("end")

	return ilBlock
}

func (node Block) Type() env.Type {
	return node.Stats[len(node.Stats)-1].Type()
}

type ExpressionStatement struct {
	Expr Expression
}

func (stat ExpressionStatement) IntermediateExpression(
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
	result shared.SymReg,
) *intermediate.BasicBlock {
	return stat.Expr.IntermediateExpression(ilProg, ilFunc, ilBlock, result)
}

func (stat ExpressionStatement) Type() env.Type {
	return stat.Expr.Type()
}
