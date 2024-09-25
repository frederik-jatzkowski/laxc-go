package abstract

import (
	"fmt"
	"laxc/internal/env"
	"laxc/pkg/attributed"
	"math"
)

func dereferenceOnce(expr attributed.Expression) (attributed.Expression, bool) {
	_, err := expr.Type().Dereference()
	if err != nil {
		return expr, false
	}

	return attributed.Dereference{Arg: expr}, true
}

func fullyDereference(expr attributed.Expression) attributed.Expression {
	var didDereference bool
	for {
		expr, didDereference = dereferenceOnce(expr)
		if !didDereference {
			return expr
		}
	}
}

func coerceBothToOneOf(arg1, arg2 attributed.Expression, permitted ...env.Type) (arg1Out, arg2Out attributed.Expression, err error) {
	shortest := math.MaxInt
	for _, t := range permitted {
		out1, coercible1, length1 := coerceToOneOf(arg1, t)
		out2, coercible2, length2 := coerceToOneOf(arg2, t)
		if !coercible1 || !coercible2 {
			continue
		}

		shorter := min(length1, length2)
		if shorter < shortest {
			shortest = shorter
			arg1Out = out1
			arg2Out = out2
		}
	}

	if shortest == math.MaxInt {
		return arg1, arg2, fmt.Errorf("there are no coercion sequences such that %s and %s can be coerced to the same permitted type", arg1.Type(), arg2.Type())
	}

	return arg1Out, arg2Out, nil
}

func coerceToOneOf(arg attributed.Expression, permitted ...env.Type) (argOut attributed.Expression, coercible bool, length int) {
	argOut = arg
	for !isPermitted(argOut, permitted...) {
		_, err := argOut.Type().Dereference()
		if err == nil {
			argOut = attributed.Dereference{Arg: argOut}
			length++

			continue
		}

		if argOut.Type().Equals(env.IntegerType{}) {
			argOut = attributed.Widen{Arg: argOut}
			length++

			continue
		}

		return arg, false, length
	}

	return argOut, true, length
}

func isPermitted(arg1 attributed.Expression, permitted ...env.Type) bool {
	for _, typ := range permitted {
		if arg1.Type().Equals(typ) {
			return true
		}
	}

	return false
}
