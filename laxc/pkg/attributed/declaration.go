package attributed

import (
	"laxc/internal/env"
	"laxc/pkg/intermediate"
)

type Declaration interface {
	IntermediateDeclaration(
		scope env.Table,
		ilProg *intermediate.Program,
		ilFunc *intermediate.Function,
		ilBlock *intermediate.BasicBlock,
	) (exitBlock *intermediate.BasicBlock)
}

type IdentityDeclaration struct {
	Ident string
	Expr  Expression
}

func (decl IdentityDeclaration) IntermediateDeclaration(
	scope env.Table,
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
) *intermediate.BasicBlock {
	variable, err := scope.GetVariable(decl.Ident)
	if err != nil {
		panic(err)
	}

	variable.SymReg = ilFunc.NextSymReg()
	ilBlock = decl.Expr.IntermediateExpression(ilProg, ilFunc, ilBlock, variable.SymReg)

	return ilBlock
}

func (expr IdentityDeclaration) Type() env.Type {
	return expr.Expr.Type()
}

type VariableDeclaration struct {
	Ident        string
	ExpectedType env.Type
}

func (node VariableDeclaration) IntermediateDeclaration(
	scope env.Table,
	ilProg *intermediate.Program,
	ilFunc *intermediate.Function,
	ilBlock *intermediate.BasicBlock,
) *intermediate.BasicBlock {
	variable, err := scope.GetVariable(node.Ident)
	if err != nil {
		panic(err)
	}

	variable.LocalSymVar = ilFunc.NextLocalSymVar()
	variable.SymReg = ilFunc.NextSymReg()

	ilBlock.TakeAddress(variable.SymReg, variable.LocalSymVar)

	return ilBlock
}

func (node VariableDeclaration) Type() env.Type {
	return node.ExpectedType
}
