package attributed

import (
	"laxc/internal/env"
	"laxc/pkg/intermediate"
)

type Program struct {
	Decls []Declaration
	Block Expression
	Scope env.Table
}

func (node Program) IntermediateProgram() intermediate.Program {
	ilProg := intermediate.NewProg()
	ilFunc := ilProg.Main()
	ilBlock := ilFunc.Entry()

	for _, decl := range node.Decls {
		ilBlock = decl.IntermediateDeclaration(node.Scope, ilProg, ilFunc, ilBlock)
	}

	resultSymReg := ilFunc.NextSymReg()
	ilBlock = node.Block.IntermediateExpression(ilProg, ilFunc, ilBlock, resultSymReg)

	if node.Block.Type().Equals(env.RealType{}) {
		ilBlock.PrintFloat(resultSymReg)
	} else {
		ilBlock.PrintInt(resultSymReg)
	}

	exitCode := ilFunc.NextSymReg()
	ilBlock.Literal(exitCode, 10)
	ilBlock.Exit(exitCode)

	return *ilProg
}
