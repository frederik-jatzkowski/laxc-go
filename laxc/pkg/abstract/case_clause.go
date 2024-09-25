package abstract

import (
	"fmt"
	"laxc/internal/env"
	"laxc/pkg/attributed"
)

type CaseClause struct {
	Condition Expression
	Cases     []Case
	Else      []Statement
}

type Case struct {
	Label Expression
	Stats []Statement
}

func (expr CaseClause) ResolveIdentifiers(scope env.Table) error {
	err := expr.Condition.ResolveIdentifiers(scope)
	if err != nil {
		return err
	}

	for _, c := range expr.Cases {
		err = c.Label.ResolveIdentifiers(scope)
		if err != nil {
			return err
		}

		for _, stat := range c.Stats {
			err = stat.Expr.ResolveIdentifiers(scope)
			if err != nil {
				return err
			}
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

func (expr CaseClause) AttributedExpression(scope env.Table) (attributed.Expression, error) {
	condition, err := expr.Condition.AttributedExpression(scope)
	if err != nil {
		return nil, err
	}

	intCondition, coercible, _ := coerceToOneOf(condition, env.IntegerType{})
	if !coercible {
		return nil, fmt.Errorf("cannot use expression of type %s as condition in case clause", condition.Type())
	}

	result := attributed.CaseClause{
		Condition: intCondition,
	}

	possibleResults := make([]attributed.Expression, 0)

	// build all statement lists
	for _, c := range expr.Cases {
		label, err := c.Label.AttributedExpression(scope)
		if err != nil {
			return nil, err
		}

		if !label.Type().Equals(env.IntegerType{}) {
			return nil, fmt.Errorf("only integers are currently supported as case labels but got %s", label.Type())
		}

		stats := make([]attributed.Expression, 0, len(c.Stats))
		for _, stat := range c.Stats {
			attributedStat, err := stat.Expr.AttributedExpression(scope)
			if err != nil {
				return nil, err
			}

			stats = append(stats, attributedStat)
		}

		result.Cases = append(result.Cases, attributed.Case{
			Label: label,
			Stats: stats,
		})

		possibleResults = append(possibleResults, stats[len(stats)-1])
	}

	for _, stat := range expr.Else {
		attributedStat, err := stat.Expr.AttributedExpression(scope)
		if err != nil {
			return nil, err
		}

		result.Else = append(result.Else, attributedStat)
	}
	possibleResults = append(possibleResults, result.Else[len(result.Else)-1])

	// determine gcd of all statement lists
	prev := possibleResults[0]
	result.ResultType = prev.Type()
	for _, possibleResult := range possibleResults {
		result.ResultType = env.GreatestCommonDenominator(
			prev.Type(),
			possibleResult.Type(),
		)

		prev = possibleResult
	}

	// coerce all as necessary
	for _, c := range result.Cases {
		c.Label, _, _ = coerceToOneOf(c.Label, result.ResultType)
		c.Stats[len(c.Stats)-1], _, _ = coerceToOneOf(c.Stats[len(c.Stats)-1], result.ResultType)
	}

	result.Else[len(result.Else)-1], _, _ = coerceToOneOf(result.Else[len(result.Else)-1], result.ResultType)

	return result, nil
}

func (expr CaseClause) Dependencies() (deps []string) {
	deps = append(deps, expr.Condition.Dependencies()...)

	for _, c := range expr.Cases {
		deps = append(deps, c.Label.Dependencies()...)
		for _, stat := range c.Stats {
			deps = append(deps, stat.Expr.Dependencies()...)
		}
	}

	for _, stat := range expr.Else {
		deps = append(deps, stat.Expr.Dependencies()...)
	}

	return deps
}
