package abstract

import (
	"fmt"
	"laxc/internal/env"
	"laxc/internal/graph"
	"laxc/pkg/attributed"
	"slices"
)

type Declarations []Declaration

type Declaration interface {
	Ident() string
	Dependencies() []string
	ResolveIdentifiers(scope env.Table) error
	AttributedDeclaration(env.Table) (attributed.Declaration, error)
}

func (nodes Declarations) ResolveIdentifiers(scope env.Table) error {
	declaredVariables := make([]string, len(nodes))
	for _, node := range nodes {
		declaredVariables = append(declaredVariables, node.Ident())
	}

	defUseGraph := graph.NewDiGraph[string](declaredVariables...)
	for _, node := range nodes {
		dependencies := node.Dependencies()

		for _, dependency := range dependencies {
			defUseGraph.AddEdge(node.Ident(), dependency)
		}
	}

	topologicalOrder, err := defUseGraph.TopologicalSort()
	if err != nil {
		return err
	}

	slices.Reverse(topologicalOrder)

	for _, ident := range topologicalOrder {
		for _, node := range nodes {
			if ident == node.Ident() {
				err := node.ResolveIdentifiers(scope)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (nodes Declarations) AttributedDeclarations(scope env.Table) (result []attributed.Declaration, err error) {
	declaredVariables := make([]string, len(nodes))
	for _, node := range nodes {
		declaredVariables = append(declaredVariables, node.Ident())
	}

	defUseGraph := graph.NewDiGraph[string](declaredVariables...)
	for _, node := range nodes {
		dependencies := node.Dependencies()

		for _, dependency := range dependencies {
			defUseGraph.AddEdge(node.Ident(), dependency)
		}
	}

	topologicalOrder, err := defUseGraph.TopologicalSort()
	if err != nil {
		return result, err
	}

	slices.Reverse(topologicalOrder)

	for _, ident := range topologicalOrder {
		for _, node := range nodes {
			if ident == node.Ident() {
				resultDecl, err := node.AttributedDeclaration(scope)
				if err != nil {
					return result, err
				}

				result = append(result, resultDecl)
			}
		}
	}

	return result, nil
}

type IdentityDeclaration struct {
	Identifier string
	Expr       Expression
	Type       string
}

func (node IdentityDeclaration) Ident() string {
	return node.Identifier
}

func (node IdentityDeclaration) Dependencies() []string {
	return node.Expr.Dependencies()
}

func (node IdentityDeclaration) ResolveIdentifiers(scope env.Table) error {
	var resultType env.Type
	switch node.Type {
	case "integer":
		resultType = env.IntegerType{}
	case "boolean":
		resultType = env.BooleanType{}
	case "real":
		resultType = env.RealType{}
	default:
		return fmt.Errorf("unknown type: %s", node.Type)
	}

	err := node.Expr.ResolveIdentifiers(scope)
	if err != nil {
		return err
	}

	err = scope.DeclareVariable(node.Identifier, &env.Variable{Name: node.Identifier, IsConstant: false, Type: resultType})
	if err != nil {
		return err
	}

	return nil
}

func (node IdentityDeclaration) AttributedDeclaration(scope env.Table) (result attributed.Declaration, err error) {
	var resultType env.Type
	switch node.Type {
	case "integer":
		resultType = env.IntegerType{}
	case "boolean":
		resultType = env.BooleanType{}
	case "real":
		resultType = env.RealType{}
	default:
		return result, fmt.Errorf("unknown type: %s", node.Type)
	}

	expr, err := node.Expr.AttributedExpression(scope)
	if err != nil {
		return result, err
	}

	resultExpr, coercible, _ := coerceToOneOf(expr, resultType)
	if !coercible {
		return result, fmt.Errorf("cannot use expression of type %s in identity declaration of type", expr.Type(), resultType)
	}

	err = scope.DeclareVariable(node.Identifier, &env.Variable{Name: node.Identifier, IsConstant: true, Type: resultType})
	if err != nil {
		return result, err
	}

	return attributed.IdentityDeclaration{
		Ident: node.Identifier,
		Expr:  resultExpr,
	}, nil
}

type VariableDeclaration struct {
	Identifier string
	Type       string
}

func (node VariableDeclaration) Ident() string {
	return node.Identifier
}

func (node VariableDeclaration) Dependencies() []string {
	return []string{}
}

func (node VariableDeclaration) ResolveIdentifiers(scope env.Table) (err error) {
	var resultType env.Type
	switch node.Type {
	case "integer":
		resultType = env.RefType{Underlying: env.IntegerType{}}
	case "boolean":
		resultType = env.RefType{Underlying: env.BooleanType{}}
	case "real":
		resultType = env.RefType{Underlying: env.RealType{}}
	default:
		return fmt.Errorf("unknown type: %s", node.Type)
	}

	err = scope.DeclareVariable(node.Identifier, &env.Variable{Name: node.Identifier, IsConstant: false, Type: resultType})
	if err != nil {
		return err
	}

	return nil
}

func (node VariableDeclaration) AttributedDeclaration(scope env.Table) (result attributed.Declaration, err error) {
	var resultType env.Type
	switch node.Type {
	case "integer":
		resultType = env.RefType{Underlying: env.IntegerType{}}
	case "boolean":
		resultType = env.RefType{Underlying: env.BooleanType{}}
	case "real":
		resultType = env.RefType{Underlying: env.RealType{}}
	default:
		return result, fmt.Errorf("unknown type: %s", node.Type)
	}

	err = scope.DeclareVariable(node.Identifier, &env.Variable{Name: node.Identifier, IsConstant: false, Type: resultType})
	if err != nil {
		return result, err
	}

	return attributed.VariableDeclaration{
		Ident:        node.Identifier,
		ExpectedType: resultType,
	}, nil
}
