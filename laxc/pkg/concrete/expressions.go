package concrete

import (
	"fmt"
	"laxc/pkg/abstract"
)

type Expression struct {
	Assignment  *Assignment  `  @@`
	Disjunction *Disjunction `| @@`
}

func (expr Expression) AbstractExpression() (abstract.Expression, error) {
	if expr.Assignment == nil {
		return expr.Disjunction.AbstractExpression()
	}

	return expr.Assignment.AbstractExpression()
}

type Assignment struct {
	Name       Name       `@@ ":="`
	Expression Expression `@@`
}

func (expr Assignment) AbstractExpression() (abstract.Expression, error) {
	arg, err := expr.Expression.AbstractExpression()
	if err != nil {
		return nil, err
	}

	return abstract.Assignment{
		Identifier: *expr.Name.Identifier,
		Arg:        arg,
	}, nil
}

type Disjunction struct {
	Conjunction Conjunction   `@@`
	More        []Conjunction `("or" @@)*`
}

func (expr Disjunction) AbstractExpression() (abstract.Expression, error) {
	result, err := expr.Conjunction.AbstractExpression()
	if err != nil {
		return result, err
	}

	for _, more := range expr.More {
		wrapper := abstract.BinaryOperation{
			Left:  result,
			OpSym: "or",
		}

		wrapper.Right, err = more.AbstractExpression()
		if err != nil {
			return result, err
		}

		result = wrapper
	}

	return result, nil
}

type Conjunction struct {
	Comparison Comparison   `@@`
	More       []Comparison `("and" @@)*`
}

func (expr Conjunction) AbstractExpression() (abstract.Expression, error) {
	result, err := expr.Comparison.AbstractExpression()
	if err != nil {
		return result, err
	}

	for _, more := range expr.More {
		wrapper := abstract.BinaryOperation{
			Left:  result,
			OpSym: "and",
		}

		wrapper.Right, err = more.AbstractExpression()
		if err != nil {
			return result, err
		}

		result = wrapper
	}

	return result, nil
}

type Comparison struct {
	Left  Relation  `@@`
	OpSym *string   `((@"=")`
	Right *Relation `@@)?`
}

func (expr Comparison) AbstractExpression() (abstract.Expression, error) {
	left, err := expr.Left.AbstractExpression()
	if err != nil {
		return nil, err
	}

	if expr.OpSym == nil {
		return left, nil
	}

	right, err := expr.Right.AbstractExpression()
	if err != nil {
		return nil, err
	}

	return abstract.BinaryOperation{
		Left:  left,
		OpSym: *expr.OpSym,
		Right: right,
	}, nil
}

type Relation struct {
	Left  Sum     `@@`
	OpSym *string `((@"<"|@">")`
	Right *Sum    `@@)?`
}

func (expr Relation) AbstractExpression() (abstract.Expression, error) {
	left, err := expr.Left.AbstractExpression()
	if err != nil {
		return nil, err
	}

	if expr.OpSym == nil {
		return left, nil
	}

	right, err := expr.Right.AbstractExpression()
	if err != nil {
		return nil, err
	}

	return abstract.BinaryOperation{
		Left:  left,
		OpSym: *expr.OpSym,
		Right: right,
	}, nil
}

type Sum struct {
	Term       Term        `@@`
	AddedTerms []AddedTerm `@@*`
}

func (expr Sum) AbstractExpression() (abstract.Expression, error) {
	result, err := expr.Term.AbstractExpression()
	if err != nil {
		return result, err
	}

	// ensure left associativity (we use a top-down parser, so this must be done here)
	for _, addedTerm := range expr.AddedTerms {
		wrapper := abstract.BinaryOperation{
			Left:  result,
			OpSym: addedTerm.OpSym,
		}

		wrapper.Right, err = addedTerm.Term.AbstractExpression()
		if err != nil {
			return result, err
		}

		result = wrapper
	}

	return result, nil
}

type AddedTerm struct {
	OpSym string `(@"+"|@"-")`
	Term  Term   `@@`
}

type Term struct {
	Factor       Factor        `@@`
	AddedFactors []AddedFactor `@@*`
}

func (expr Term) AbstractExpression() (abstract.Expression, error) {
	result, err := expr.Factor.AbstractExpression()
	if err != nil {
		return result, err
	}

	// ensure left associativity (we use a top-down parser, so this must be done here)
	for _, addedFactor := range expr.AddedFactors {
		wrapper := abstract.BinaryOperation{
			Left:  result,
			OpSym: addedFactor.OpSym,
		}

		wrapper.Right, err = addedFactor.Factor.AbstractExpression()
		if err != nil {
			return result, err
		}

		result = wrapper
	}

	return result, nil
}

type AddedFactor struct {
	OpSym  string `(@"*"|@"div"|@"mod"|@"/")`
	Factor Factor `@@`
}

type Factor struct {
	Primary    *Primary    `  @@`
	UnopFactor *UnopFactor `| @@`
}

type UnopFactor struct {
	OpSym  string `(@"+"|@"-"|@"not")`
	Factor Factor `@@`
}

func (expr Factor) AbstractExpression() (abstract.Expression, error) {
	if expr.Primary != nil {
		return expr.Primary.AbstractExpression()
	} else if expr.UnopFactor != nil {
		arg, err := expr.UnopFactor.Factor.AbstractExpression()

		return abstract.UnaryOperation{
			OpSym: expr.UnopFactor.OpSym,
			Arg:   arg,
		}, err
	}

	return nil, fmt.Errorf("no option was matched in factor: %+v", expr)
}

type Primary struct {
	Denotation *Denotation `  @@`
	Name       *Name       `| @@`
	Block      *Block      `| @@`
	Clause     *Clause     `| @@`
	Expression *Expression `| "(":Special@@")":Special`
}

func (expr Primary) AbstractExpression() (abstract.Expression, error) {
	if expr.Denotation != nil {
		return expr.Denotation.AbstractExpression()
	} else if expr.Name != nil {
		return expr.Name.AbstractExpression()
	} else if expr.Block != nil {
		return expr.Block.AbstractExpression()
	} else if expr.Clause != nil {
		return expr.Clause.AbstractExpression()
	} else if expr.Expression != nil {
		return expr.Expression.AbstractExpression()
	}

	return nil, fmt.Errorf("no option was matched in factor: %+v", expr)
}

type Denotation struct {
	Int  *string `  @IntLit`
	Real *string `| @RealLit`
}

func (expr Denotation) AbstractExpression() (abstract.Expression, error) {
	if expr.Int != nil {
		return abstract.IntegerLiteral{
			Int: *expr.Int,
		}, nil
	} else {
		return abstract.RealLiteral{
			Real: *expr.Real,
		}, nil
	}
}

type Name struct {
	Identifier *string `@Ident`
}

func (expr Name) AbstractExpression() (abstract.Expression, error) {
	return abstract.UseVariable{
		Identifier: *expr.Identifier,
	}, nil
}

type Clause struct {
	ConditionalClause *ConditionalClause `  @@`
	CaseClause        *CaseClause        `| @@`
}

func (expr Clause) AbstractExpression() (abstract.Expression, error) {
	if expr.ConditionalClause != nil {
		return expr.ConditionalClause.AbstractExpression()
	}

	return expr.CaseClause.AbstractExpression()
}

type ConditionalClause struct {
	If              string         `"if":Keyword`
	Expr            Expression     `@@`
	Then            string         `"then":Keyword`
	TrueStatements  StatementList  `@@`
	Else            *string        `("else":Keyword`
	FalseStatements *StatementList `@@)?`
	End             string         `"end":Keyword`
}

func (expr ConditionalClause) AbstractExpression() (abstract.Expression, error) {
	var (
		err error
	)

	condition, err := expr.Expr.AbstractExpression()
	if err != nil {
		return nil, err
	}

	if expr.FalseStatements != nil {
		result := abstract.TwoSidedConditionalClause{
			Condition: condition,
		}

		for _, stat := range expr.TrueStatements.Stats {
			abstractStat, err := stat.GenerateAstNode()
			if err != nil {
				return result, nil
			}

			result.Then = append(result.Then, abstractStat)
		}

		for _, stat := range expr.FalseStatements.Stats {
			abstractStat, err := stat.GenerateAstNode()
			if err != nil {
				return result, nil
			}

			result.Else = append(result.Else, abstractStat)
		}

		return result, nil
	} else {
		result := abstract.OneSidedConditionalClause{
			Condition: condition,
		}

		for _, stat := range expr.TrueStatements.Stats {
			abstractStat, err := stat.GenerateAstNode()
			if err != nil {
				return result, nil
			}

			result.Then = append(result.Then, abstractStat)
		}

		return result, nil
	}
}

type Case struct {
	Label      Denotation    `@@ ":"`
	Statements StatementList `@@`
}

type CaseClause struct {
	Condition Expression    `"case" @@ "of"`
	Case      Case          `@@`
	MoreCases []Case        `("//":Keyword @@)*`
	Else      StatementList `"else" @@ "end"`
}

func (expr CaseClause) AbstractExpression() (abstract.Expression, error) {
	condition, err := expr.Condition.AbstractExpression()
	if err != nil {
		return nil, err
	}

	result := abstract.CaseClause{
		Condition: condition,
	}

	for _, c := range append([]Case{expr.Case}, expr.MoreCases...) {
		label, err := c.Label.AbstractExpression()
		if err != nil {
			return nil, err
		}

		stats := make([]abstract.Statement, 0, len(c.Statements.Stats))
		for _, stat := range c.Statements.Stats {
			abstractStat, err := stat.GenerateAstNode()
			if err != nil {
				return nil, err
			}

			stats = append(stats, abstractStat)
		}

		result.Cases = append(result.Cases, abstract.Case{
			Label: label,
			Stats: stats,
		})
	}

	for _, stat := range expr.Else.Stats {
		abstractStat, err := stat.GenerateAstNode()
		if err != nil {
			return nil, err
		}

		result.Else = append(result.Else, abstractStat)
	}

	return result, nil
}
