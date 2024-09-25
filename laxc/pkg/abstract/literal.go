package abstract

import (
	"fmt"
	"laxc/internal/env"
	"laxc/pkg/attributed"
	"math/big"
	"strconv"
)

type IntegerLiteral struct {
	Int string
}

func (expr IntegerLiteral) ResolveIdentifiers(scope env.Table) error {
	return nil
}

func (expr IntegerLiteral) AttributedExpression(scope env.Table) (attributed.Expression, error) {
	parsedInt, err := strconv.ParseInt(expr.Int, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Integer does not fit in range: %w", err)
	}
	return attributed.IntegerLiteral{Value: int32(parsedInt)}, nil
}

func (expr IntegerLiteral) Dependencies() []string {
	return []string{}
}

type RealLiteral struct {
	Real string
}

func (expr RealLiteral) ResolveIdentifiers(scope env.Table) error {
	return nil
}

func (expr RealLiteral) AttributedExpression(scope env.Table) (attributed.Expression, error) {
	value, _, err := big.ParseFloat(expr.Real, 10, 32, big.ToNearestAway)
	if err != nil {
		return nil, err
	}

	return attributed.RealLiteral{Value: *value}, nil
}

func (expr RealLiteral) Dependencies() []string {
	return []string{}
}
