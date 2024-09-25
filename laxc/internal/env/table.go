package env

import (
	"fmt"
)

type Table struct {
	outer     *Table
	Variables map[string]*Variable
}

func NewTable(outer *Table) Table {
	return Table{
		outer:     outer,
		Variables: make(map[string]*Variable),
	}
}

func (table Table) IsDefined(ident string) bool {
	return table.IsDefinedHere(ident) || table.outer != nil && table.outer.IsDefined(ident)
}

func (table Table) IsDefinedHere(ident string) bool {
	switch ident {
	case "boolean", "integer":
		return true
	}

	_, exists := table.Variables[ident]

	return exists
}

func (table Table) DeclareVariable(ident string, variable *Variable) error {
	if table.IsDefinedHere(ident) {
		return fmt.Errorf("identifier already defined here: %s", ident)
	}

	table.Variables[ident] = variable

	return nil
}

func (table Table) GetVariable(ident string) (*Variable, error) {
	variable, exists := table.Variables[ident]
	if exists {
		return variable, nil
	}

	if table.outer != nil {
		return table.outer.GetVariable(ident)
	}

	return variable, fmt.Errorf("undeclared identifier: %s", ident)
}
