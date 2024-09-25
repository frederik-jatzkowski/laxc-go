package abstract

import (
	"laxc/internal/env"
	"laxc/pkg/attributed"
	"slices"
)

type Block struct {
	Decls Declarations
	Stats []Statement
}

func (expr Block) ResolveIdentifiers(scope env.Table) error {
	innerScope := env.NewTable(&scope)
	err := expr.Decls.ResolveIdentifiers(innerScope)
	if err != nil {
		return err
	}

	for _, stat := range expr.Stats {
		err = stat.Expr.ResolveIdentifiers(innerScope)
		if err != nil {
			return err
		}
	}

	return nil
}

func (expr Block) AttributedExpression(scope env.Table) (attributed.Expression, error) {
	var err error

	result := attributed.Block{}
	result.Scope = env.NewTable(&scope)

	result.Decls, err = expr.Decls.AttributedDeclarations(result.Scope)
	if err != nil {
		return result, err
	}

	for _, stat := range expr.Stats {
		resultStat, err := stat.AttributedStatement(result.Scope)
		if err != nil {
			return result, err
		}

		result.Stats = append(result.Stats, resultStat)
	}

	return result, nil
}

func (expr Block) Dependencies() (dependencies []string) {
	for _, stat := range expr.Stats {
		dependencies = append(dependencies, stat.Expr.Dependencies()...)
	}

	for _, decl := range expr.Decls {
		dependencies = append(dependencies, decl.Dependencies()...)
	}

	for _, decl := range expr.Decls {
		dependencies = slices.DeleteFunc(dependencies, func(dependency string) bool {
			return dependency == decl.Ident()
		})
	}

	return dependencies
}
