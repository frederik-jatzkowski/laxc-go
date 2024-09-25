package mips32

import (
	"fmt"
	"laxc/internal/shared"
)

type Instr interface {
	Label() string
	Name() string
	Args() string
	Comment() string
}

type instr0 struct {
	name    string
	comment string
}

func (instr instr0) Label() string {
	return ""
}

func (instr instr0) Name() string {
	return instr.name
}

func (instr instr0) Args() string {
	return ""
}

func (instr instr0) Comment() string {
	return instr.comment
}

type instr1 struct {
	name    string
	arg1    shared.Reg
	comment string
}

func (instr instr1) Label() string {
	return ""
}

func (instr instr1) Name() string {
	return instr.name
}

func (instr instr1) Args() string {
	return string(instr.arg1)
}

func (instr instr1) Comment() string {
	return instr.comment
}

type instr2 struct {
	name    string
	arg1    shared.Reg
	arg2    shared.Reg
	comment string
}

func (instr instr2) Label() string {
	return ""
}

func (instr instr2) Name() string {
	return instr.name
}

func (instr instr2) Args() string {
	return fmt.Sprintf("%s,%s", instr.arg1, instr.arg2)
}

func (instr instr2) Comment() string {
	return instr.comment
}

type instr3 struct {
	name    string
	arg1    shared.Reg
	arg2    shared.Reg
	arg3    shared.Reg
	comment string
}

func (instr instr3) Label() string {
	return ""
}

func (instr instr3) Name() string {
	return instr.name
}

func (instr instr3) Args() string {
	return fmt.Sprintf("%s,%s,%s", instr.arg1, instr.arg2, instr.arg3)
}

func (instr instr3) Comment() string {
	return instr.comment
}

type instr2I struct {
	name    string
	arg1    shared.Reg
	immed   int32
	comment string
}

func (instr instr2I) Label() string {
	return ""
}

func (instr instr2I) Name() string {
	return instr.name
}

func (instr instr2I) Args() string {
	return fmt.Sprintf("%s,%d", instr.arg1, instr.immed)
}

func (instr instr2I) Comment() string {
	return instr.comment
}

type instr3I struct {
	name    string
	arg1    shared.Reg
	arg2    shared.Reg
	immed   int32
	comment string
}

func (instr instr3I) Label() string {
	return ""
}

func (instr instr3I) Name() string {
	return instr.name
}

func (instr instr3I) Args() string {
	return fmt.Sprintf("%s,%s,%d", instr.arg1, instr.arg2, instr.immed)
}

func (instr instr3I) Comment() string {
	return instr.comment
}

type label struct {
	label   string
	comment string
}

func (instr label) Label() string {
	return instr.label + ":"
}

func (instr label) Name() string {
	return ""
}

func (instr label) Args() string {
	return ""
}

func (instr label) Comment() string {
	return instr.comment
}
