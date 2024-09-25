package mips32

import (
	"fmt"
	"laxc/internal/shared"
)

func (prog *Program) ADD(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"add", arg1, arg2, arg3, comment})
}

func (prog *Program) OR(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"or", arg1, arg2, arg3, comment})
}

func (prog *Program) AND(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"and", arg1, arg2, arg3, comment})
}

func (prog *Program) SUB(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"sub", arg1, arg2, arg3, comment})
}

func (prog *Program) XOR(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"xor", arg1, arg2, arg3, comment})
}

func (prog *Program) ADDI(arg1, arg2 shared.Reg, immed int32, comment string) {
	prog.instrs = append(prog.instrs, instr3I{"addi", arg1, arg2, immed, comment})
}

func (prog *Program) MULT(arg1, arg2 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr2{"mult", arg1, arg2, comment})
}

func (prog *Program) DIV(arg1, arg2 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr2{"div", arg1, arg2, comment})
}

func (prog *Program) MFHI(arg1 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr1{"mfhi", arg1, comment})
}

func (prog *Program) MFLO(arg1 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr1{"mflo", arg1, comment})
}

func (prog *Program) SRA(arg1, arg2 shared.Reg, immed int32, comment string) {
	prog.instrs = append(prog.instrs, instr3I{"sra", arg1, arg2, immed, comment})
}

func (prog *Program) TNE(arg1, arg2 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr2{"tne", arg1, arg2, comment})
}

func (prog *Program) TEQ(arg1, arg2 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr2{"teq", arg1, arg2, comment})
}

func (prog *Program) ORI(arg1, arg2 shared.Reg, immed int32, comment string) {
	prog.instrs = append(prog.instrs, instr3I{"ori", arg1, arg2, immed, comment})
}

func (prog *Program) LUI(arg1 shared.Reg, immed int32, comment string) {
	prog.instrs = append(prog.instrs, instr2I{"lui", arg1, immed, comment})
}

func (prog *Program) SYSCALL(comment string) {
	prog.instrs = append(prog.instrs, instr0{"syscall", comment})
}

func (prog *Program) SW(arg1, arg2 shared.Reg, immed int16, comment string) {
	prog.instrs = append(prog.instrs, instr2{"sw", arg1, shared.Reg(fmt.Sprintf("%d(%s)", immed, arg2)), comment})
}

func (prog *Program) LW(arg1, arg2 shared.Reg, immed int16, comment string) {
	prog.instrs = append(prog.instrs, instr2{"lw", arg1, shared.Reg(fmt.Sprintf("%d(%s)", immed, arg2)), comment})
}

func (prog *Program) SLT(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"slt", arg1, arg2, arg3, comment})
}

func (prog *Program) NOR(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"nor", arg1, arg2, arg3, comment})
}

func (prog *Program) ANDI(arg1, arg2 shared.Reg, immed int32, comment string) {
	prog.instrs = append(prog.instrs, instr3I{"andi", arg1, arg2, immed, comment})
}

func (prog *Program) MTC1(arg1, arg2 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr2{"mtc1", arg1, arg2, comment})
}

func (prog *Program) MFC1(arg1, arg2 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr2{"mfc1", arg1, arg2, comment})
}

func (prog *Program) CVTSW(arg1, arg2 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr2{"cvt.s.w", arg1, arg2, comment})
}

func (prog *Program) ADDS(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"add.s", arg1, arg2, arg3, comment})
}

func (prog *Program) SUBS(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"sub.s", arg1, arg2, arg3, comment})
}

func (prog *Program) MULS(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"mul.s", arg1, arg2, arg3, comment})
}

func (prog *Program) DIVS(arg1, arg2, arg3 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr3{"div.s", arg1, arg2, arg3, comment})
}

func (prog *Program) LABEL(name string, comment string) {
	prog.instrs = append(prog.instrs, label{name, comment})
}

func (prog *Program) CLTS(arg1, arg2 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr2{"c.lt.s", arg1, arg2, comment})
}

func (prog *Program) CEQS(arg1, arg2 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr2{"c.eq.s", arg1, arg2, comment})
}

func (prog *Program) BC1T(label string, comment string) {
	prog.instrs = append(prog.instrs, instr1{"bc1t", shared.Reg(label), comment})
}

func (prog *Program) NEGS(arg1, arg2 shared.Reg, comment string) {
	prog.instrs = append(prog.instrs, instr2{"neg.s", arg1, arg2, comment})
}

func (prog *Program) BEQ(arg1, arg2 shared.Reg, label string, comment string) {
	prog.instrs = append(prog.instrs, instr3{"beq", arg1, arg2, shared.Reg(label), comment})
}

func (prog *Program) BNE(arg1, arg2 shared.Reg, label string, comment string) {
	prog.instrs = append(prog.instrs, instr3{"bne", arg1, arg2, shared.Reg(label), comment})
}

func (prog *Program) J(label string, comment string) {
	prog.instrs = append(prog.instrs, instr1{"j", shared.Reg(label), comment})
}
