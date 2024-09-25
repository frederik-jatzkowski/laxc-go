package abstract

import (
	"fmt"
	"laxc/internal/env"
	"laxc/pkg/attributed"
)

type Expression interface {
	ResolveIdentifiers(scope env.Table) error
	AttributedExpression(scope env.Table) (attributed.Expression, error)
	Dependencies() []string
}

type UseVariable struct {
	Identifier string
}

func (expr UseVariable) ResolveIdentifiers(scope env.Table) error {
	_, err := scope.GetVariable(expr.Identifier)

	return err
}

func (expr UseVariable) AttributedExpression(scope env.Table) (attributed.Expression, error) {
	variable, err := scope.GetVariable(expr.Identifier)
	if err != nil {
		return nil, err
	}

	return attributed.ReadVariable{Var: variable}, nil
}

func (expr UseVariable) Dependencies() []string {
	return []string{expr.Identifier}
}

type Assignment struct {
	Identifier string
	Arg        Expression
}

func (expr Assignment) ResolveIdentifiers(scope env.Table) error {
	_, err := scope.GetVariable(expr.Identifier)
	if err != nil {
		return err
	}

	return expr.Arg.ResolveIdentifiers(scope)
}

func (expr Assignment) AttributedExpression(scope env.Table) (attributed.Expression, error) {
	arg, err := expr.Arg.AttributedExpression(scope)
	if err != nil {
		return nil, err
	}

	variable, err := scope.GetVariable(expr.Identifier)
	if err != nil {
		return nil, err
	}

	varType, err := variable.Type.Dereference()
	if err != nil {
		return nil, err
	}

	arg, coercible, _ := coerceToOneOf(arg, varType)
	if !coercible {
		return nil, fmt.Errorf("cannot assign expression of type %s to variable of type %s", arg.Type(), variable.Type)
	}

	return attributed.Assignment{
		Var:            variable,
		UnderlyingType: arg.Type(),
		Arg:            arg,
	}, nil
}

func (expr Assignment) Dependencies() []string {
	return append(expr.Arg.Dependencies(), expr.Identifier)
}
