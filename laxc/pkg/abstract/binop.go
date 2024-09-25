package abstract

import (
	"fmt"
	"laxc/internal/env"
	"laxc/pkg/attributed"
)

type BinaryOperation struct {
	Left  Expression
	OpSym string
	Right Expression
}

func (expr BinaryOperation) ResolveIdentifiers(scope env.Table) error {
	err := expr.Left.ResolveIdentifiers(scope)
	if err != nil {
		return err
	}

	return expr.Right.ResolveIdentifiers(scope)
}

func (expr BinaryOperation) AttributedExpression(scope env.Table) (attributed.Expression, error) {
	left, err := expr.Left.AttributedExpression(scope)
	if err != nil {
		return nil, err
	}

	right, err := expr.Right.AttributedExpression(scope)
	if err != nil {
		return nil, err
	}

	switch expr.OpSym {
	case "+":
		left, right, err = coerceBothToOneOf(left, right, env.IntegerType{}, env.RealType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in addition operation: %w", err)
		}

		if left.Type().Equals(env.IntegerType{}) {
			return attributed.IntAdd{
				Left:  left,
				OpSym: expr.OpSym,
				Right: right,
			}, nil
		}

		return attributed.RealAdd{
			Left:  left,
			OpSym: expr.OpSym,
			Right: right,
		}, nil
	case "-":
		left, right, err = coerceBothToOneOf(left, right, env.IntegerType{}, env.RealType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in subtraction operation: %w", err)
		}

		if left.Type().Equals(env.IntegerType{}) {
			return attributed.IntSub{
				Left:  left,
				OpSym: expr.OpSym,
				Right: right,
			}, nil
		}

		return attributed.RealSub{
			Left:  left,
			OpSym: expr.OpSym,
			Right: right,
		}, nil
	case "*":
		left, right, err = coerceBothToOneOf(left, right, env.IntegerType{}, env.RealType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in multiplication operation: %w", err)
		}

		if left.Type().Equals(env.IntegerType{}) {
			return attributed.IntMul{
				Left:  left,
				OpSym: expr.OpSym,
				Right: right,
			}, nil
		}

		return attributed.RealMul{
			Left:  left,
			OpSym: expr.OpSym,
			Right: right,
		}, nil
	case "div":
		left, right, err = coerceBothToOneOf(left, right, env.IntegerType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in integer division: %w", err)
		}

		return attributed.IntDiv{
			Left:  left,
			OpSym: expr.OpSym,
			Right: right,
		}, nil
	case "/":
		left, right, err = coerceBothToOneOf(left, right, env.RealType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in floating point division: %w", err)
		}

		return attributed.RealDiv{
			Left:  left,
			OpSym: expr.OpSym,
			Right: right,
		}, nil
	case "mod":
		left, right, err = coerceBothToOneOf(left, right, env.IntegerType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in integer division operation: %w", err)
		}

		return attributed.IntMod{
			Left:  left,
			OpSym: expr.OpSym,
			Right: right,
		}, nil
	case "<":
		left, right, err = coerceBothToOneOf(left, right, env.IntegerType{}, env.RealType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in less than comparison: %w", err)
		}

		if left.Type().Equals(env.IntegerType{}) {
			return attributed.IntLessThan{
				Left:  left,
				OpSym: expr.OpSym,
				Right: right,
			}, nil
		}

		return attributed.RealLessThan{
			Left:  left,
			OpSym: expr.OpSym,
			Right: right,
		}, nil
	case ">":
		left, right, err = coerceBothToOneOf(left, right, env.IntegerType{}, env.RealType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in greater than comparison: %w", err)
		}

		if left.Type().Equals(env.IntegerType{}) {
			return attributed.IntGreaterThan{
				Left:  left,
				OpSym: expr.OpSym,
				Right: right,
			}, nil
		}

		return attributed.RealGreaterThan{
			Left:  left,
			OpSym: expr.OpSym,
			Right: right,
		}, nil
	case "=":
		left, right, err = coerceBothToOneOf(left, right, env.IntegerType{}, env.RealType{}, env.BooleanType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in equality: %w", err)
		}

		if left.Type().Equals(env.RealType{}) {
			return attributed.RealEquals{
				Left:  left,
				OpSym: expr.OpSym,
				Right: right,
			}, nil
		}

		return attributed.Equals{
			Left:  left,
			OpSym: expr.OpSym,
			Right: right,
		}, nil
	case "and":
		left, right, err = coerceBothToOneOf(left, right, env.BooleanType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in boolean conjunction: %w", err)
		}

		return attributed.BoolAnd{
			Left:  left,
			Right: right,
		}, nil
	case "or":
		left, right, err = coerceBothToOneOf(left, right, env.BooleanType{})
		if err != nil {
			return nil, fmt.Errorf("invalid type in boolean disjunction: %w", err)
		}

		return attributed.BoolOr{
			Left:  left,
			Right: right,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected operation symbol: %s", expr.OpSym)
	}

}
func (expr BinaryOperation) Dependencies() []string {
	return append(expr.Left.Dependencies(), expr.Right.Dependencies()...)
}
