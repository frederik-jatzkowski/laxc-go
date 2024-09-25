package attributed

import (
	"laxc/internal/env"
	"laxc/internal/shared"
	"laxc/pkg/intermediate"
)

type Expression interface {
	IntermediateExpression(
		ilProg *intermediate.Program,
		ilFunc *intermediate.Function,
		ilBlock *intermediate.BasicBlock,
		result shared.SymReg,
	) (exitBlock *intermediate.BasicBlock)
	Type() env.Type
}
