package concrete

import (
	"fmt"
	"laxc/pkg/abstract"
)

type Prog struct {
	Block Block `@@`
}

func (node Prog) AbstractExpression() (result abstract.Program, err error) {
	result.Decls = abstract.Declarations{
		abstract.IdentityDeclaration{
			Identifier: "true",
			Expr: abstract.BinaryOperation{
				Left: abstract.IntegerLiteral{
					Int: "1",
				},
				OpSym: "<",
				Right: abstract.IntegerLiteral{
					Int: "2",
				},
			},
			Type: "boolean",
		},
		abstract.IdentityDeclaration{
			Identifier: "false",
			Expr: abstract.UnaryOperation{
				OpSym: "not",
				Arg: abstract.UseVariable{
					Identifier: "true",
				},
			},
			Type: "boolean",
		},
	}

	for _, decl := range node.Block.Decls {
		abstractDecl, err := decl.GenerateAstNode()
		if err != nil {
			return result, err
		}

		result.Block.Decls = append(result.Block.Decls, abstractDecl)
	}

	for _, stat := range node.Block.StatementList.Stats {
		abstractStat, err := stat.GenerateAstNode()
		if err != nil {
			return result, err
		}

		result.Block.Stats = append(result.Block.Stats, abstractStat)
	}

	return result, nil
}

type Block struct {
	Declare       string        `"declare":Keyword`
	Decls         []Declaration `(@@(";":Special@@)*)`
	Begin         string        `"begin":Keyword`
	StatementList StatementList `@@`
	End           string        `"end":Keyword`
}

func (node Block) AbstractExpression() (abstract.Expression, error) {
	result := abstract.Block{}

	for _, decl := range node.Decls {
		abstractDecl, err := decl.GenerateAstNode()
		if err != nil {
			return result, err
		}

		result.Decls = append(result.Decls, abstractDecl)
	}

	for _, stat := range node.StatementList.Stats {
		abstractStat, err := stat.GenerateAstNode()
		if err != nil {
			return result, err
		}

		result.Stats = append(result.Stats, abstractStat)
	}

	return result, nil
}

type StatementList struct {
	Stats []Statement `(@@(";":Special@@)*)`
}

type Declaration struct {
	Ident string      `@Ident`
	Expr  *Expression `("is":Keyword @@)?`
	Type  string      `":":Special @Ident`
}

func (node Declaration) GenerateAstNode() (result abstract.Declaration, err error) {
	switch node.Ident {
	case "true", "false", "integer", "boolean", "real", "nil":
		return result, fmt.Errorf("%s is a predefined identifier", node.Ident)
	}

	if node.Expr != nil {
		abstractExpr, err := node.Expr.AbstractExpression()
		if err != nil {
			return result, err
		}

		return abstract.IdentityDeclaration{
			Identifier: node.Ident,
			Expr:       abstractExpr,
			Type:       node.Type,
		}, nil
	}

	return abstract.VariableDeclaration{
		Identifier: node.Ident,
		Type:       node.Type,
	}, nil
}

type Statement struct {
	Expr Expression `@@`
}

func (node Statement) GenerateAstNode() (result abstract.Statement, err error) {
	abstractExpr, err := node.Expr.AbstractExpression()
	if err != nil {
		return result, err
	}

	return abstract.Statement{Expr: abstractExpr}, nil
}
