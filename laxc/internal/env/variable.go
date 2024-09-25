package env

import "laxc/internal/shared"

type Variable struct {
	Name        string
	IsConstant  bool
	SymReg      shared.SymReg
	LocalSymVar shared.LocalSymVar
	Type        Type
}
