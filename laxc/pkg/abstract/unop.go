package abstract

import (
	"fmt"
	"laxc/internal/env"
	"laxc/pkg/attributed"
)

type UnaryOperation struct {
	OpSym string
	Arg   Expression
}

func (expr UnaryOperation) ResolveIdentifiers(scope env.Table) error {
	return expr.Arg.ResolveIdentifiers(scope)
}

func (expr UnaryOperation) AttributedExpression(scope env.Table) (attributed.Expression, error) {
	arg, err := expr.Arg.AttributedExpression(scope)
	if err != nil {
		return nil, err
	}

	arg = fullyDereference(arg)

	switch expr.OpSym {
	case "+":
		arg, coercible, _ := coerceToOneOf(arg, env.IntegerType{}, env.RealType{})
		if !coercible {
			return nil, fmt.Errorf("invalid type in unary plus operation: %s", arg.Type())
		}

		return arg, nil
	case "-":
		arg, coercible, _ := coerceToOneOf(arg, env.IntegerType{}, env.RealType{})
		if !coercible {
			return nil, fmt.Errorf("invalid type in unary numeric negation operation: %s", arg.Type())
		}

		if arg.Type().Equals(env.IntegerType{}) {
			return attributed.IntNeg{
				OpSym: expr.OpSym,
				Arg:   arg,
			}, nil
		}

		return attributed.RealNeg{
			OpSym: expr.OpSym,
			Arg:   arg,
		}, nil
	case "not":
		arg, coercible, _ := coerceToOneOf(arg, env.BooleanType{})
		if !coercible {
			return nil, fmt.Errorf("invalid type in unary not operation: %s", arg.Type())
		}

		return attributed.BoolNeg{
			OpSym: expr.OpSym,
			Arg:   arg,
		}, nil
	default:
		return nil, fmt.Errorf("unexpected operation symbol: %s", expr.OpSym)
	}
}

func (expr UnaryOperation) Dependencies() []string {
	return expr.Arg.Dependencies()
}
