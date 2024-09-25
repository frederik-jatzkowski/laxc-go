package abstract

import (
	"fmt"
	"laxc/internal/env"
	"laxc/pkg/attributed"
)

type OneSidedConditionalClause struct {
	Condition Expression
	Then      []Statement
}

func (expr OneSidedConditionalClause) ResolveIdentifiers(scope env.Table) error {
	err := expr.Condition.ResolveIdentifiers(scope)
	if err != nil {
		return err
	}

	for _, stat := range expr.Then {
		err = stat.Expr.ResolveIdentifiers(scope)
		if err != nil {
			return err
		}
	}

	return nil
}

func (expr OneSidedConditionalClause) AttributedExpression(scope env.Table) (attributed.Expression, error) {
	condition, err := expr.Condition.AttributedExpression(scope)
	if err != nil {
		return nil, err
	}

	boolCondition, coercible, _ := coerceToOneOf(condition, env.BooleanType{})
	if !coercible {
		return nil, fmt.Errorf("cannot use expression of type %s as condition in conditional clause", condition.Type())
	}

	result := attributed.OneSidedConditionalClause{
		Condition: boolCondition,
	}

	for _, stat := range expr.Then {
		attributedStat, err := stat.AttributedStatement(scope)
		if err != nil {
			return nil, err
		}

		result.Then = append(result.Then, attributedStat)
	}

	return result, nil

}

func (expr OneSidedConditionalClause) Dependencies() (deps []string) {
	deps = append(deps, expr.Condition.Dependencies()...)

	for _, stat := range expr.Then {
		deps = append(deps, stat.Expr.Dependencies()...)
	}

	return deps
}

type TwoSidedConditionalClause struct {
	Condition Expression
	Then      []Statement
	Else      []Statement
}

func (expr TwoSidedConditionalClause) ResolveIdentifiers(scope env.Table) error {
	err := expr.Condition.ResolveIdentifiers(scope)
	if err != nil {
		return err
	}

	for _, stat := range expr.Then {
		err = stat.Expr.ResolveIdentifiers(scope)
		if err != nil {
			return err
		}
	}

	for _, stat := range expr.Else {
		err = stat.Expr.ResolveIdentifiers(scope)
		if err != nil {
			return err
		}
	}

	return nil
}

func (expr TwoSidedConditionalClause) AttributedExpression(scope env.Table) (attributed.Expression, error) {
	condition, err := expr.Condition.AttributedExpression(scope)
	if err != nil {
		return nil, err
	}

	boolCondition, coercible, _ := coerceToOneOf(condition, env.BooleanType{})
	if !coercible {
		return nil, fmt.Errorf("cannot use expression of type %s as condition in conditional clause", condition.Type())
	}

	result := attributed.TwoSidedConditionalClause{
		Condition:  boolCondition,
		ResultType: env.VoidType{},
	}

	for _, stat := range expr.Then {
		attributedStat, err := stat.AttributedStatement(scope)
		if err != nil {
			return nil, err
		}

		result.Then = append(result.Then, attributedStat)
	}

	for _, stat := range expr.Else {
		attributedStat, err := stat.AttributedStatement(scope)
		if err != nil {
			return nil, err
		}

		result.Else = append(result.Else, attributedStat)
	}

	thenResult := result.Then[len(result.Then)-1]
	elseResult := result.Else[len(result.Else)-1]

	gcd := env.GreatestCommonDenominator(thenResult.Type(), elseResult.Type())

	result.Then[len(result.Then)-1], result.Else[len(result.Else)-1], err = coerceBothToOneOf(
		thenResult,
		elseResult,
		gcd,
	)
	if err != nil {
		return nil, err
	}

	result.ResultType = result.Then[len(result.Then)-1].Type()

	return result, nil
}

func (expr TwoSidedConditionalClause) Dependencies() (deps []string) {
	deps = append(deps, expr.Condition.Dependencies()...)

	for _, stat := range expr.Then {
		deps = append(deps, stat.Expr.Dependencies()...)
	}

	for _, stat := range expr.Else {
		deps = append(deps, stat.Expr.Dependencies()...)
	}

	return deps
}
