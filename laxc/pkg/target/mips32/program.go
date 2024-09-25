package mips32

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

type Program struct {
	instrs []Instr
}

func NewProgram() *Program {
	return &Program{}
}

func (prog Program) String() string {
	buf := bytes.NewBufferString("")
	writer := tabwriter.NewWriter(buf, 0, 2, 4, ' ', 0)

	for _, instr := range prog.instrs {
		writer.Write([]byte(fmt.Sprintf("%s\t%s\t%s\t# %s\n", instr.Label(), instr.Name(), instr.Args(), instr.Comment())))
	}

	writer.Flush()

	return buf.String()
}

func (prog Program) Line() int {
	return len(prog.instrs) + 1
}
