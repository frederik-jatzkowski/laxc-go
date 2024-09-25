package abstract

import (
	"fmt"
	"laxc/internal/env"
	"laxc/pkg/attributed"
)

type Program struct {
	Decls Declarations
	Block Block
}

func (node Program) AttributedProgram() (result attributed.Program, err error) {
	{
		testScope := env.NewTable(nil)

		err := node.Decls.ResolveIdentifiers(testScope)
		if err != nil {
			return result, err
		}

		err = node.Block.ResolveIdentifiers(testScope)
		if err != nil {
			return result, err
		}
	}

	result.Scope = env.NewTable(nil)

	result.Decls, err = node.Decls.AttributedDeclarations(result.Scope)
	if err != nil {
		return result, err
	}

	result.Block, err = node.Block.AttributedExpression(result.Scope)
	if err != nil {
		return result, err
	}

	coercedBlock, coercible, _ := coerceToOneOf(result.Block, env.IntegerType{}, env.BooleanType{}, env.RealType{})
	if !coercible {
		return result, fmt.Errorf("block of type %s cannot be used as a program result", result.Block.Type())
	}

	result.Block = coercedBlock

	return result, nil
}

type Statement struct {
	Expr Expression
}

func (stat Statement) AttributedStatement(scope env.Table) (result attributed.Expression, err error) {
	expr, err := stat.Expr.AttributedExpression(scope)
	if err != nil {
		return result, err
	}

	return attributed.ExpressionStatement{
		Expr: expr,
	}, nil
}
