package bytecode

import (
	"fmt"
	"laxc/internal/shared"
)

func (prog *Program) LIT(result shared.Reg, value int32) {
	instr := NewInstruction(OP_LIT, parse(result), 0, 0)
	data := Instruction(value)
	prog.instrs = append(prog.instrs, instr, data)
}

func (prog *Program) INT_PRINT(arg shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_INT_PRINT, 0, parse(arg), 0))
}

func (prog *Program) INT_ADD(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_INT_ADD, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) INT_SUB(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_INT_SUB, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) INT_MUL(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_INT_MUL, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) INT_DIV(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_INT_DIV, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) INT_MOD(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_INT_MOD, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) INT_NEG(result, arg shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_INT_NEG, parse(result), parse(arg), 0))
}

func (prog *Program) INT_LT(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_INT_LT, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) INT_EQ(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_INT_EQ, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) OR(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_OR, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) BOOL_NOT(result, arg shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_BOOL_NOT, parse(result), parse(arg), 0))
}

func (prog *Program) INT_TO_FLOAT(arg, result shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_INT_TO_FLOAT, parse(result), parse(arg), 0))
}

func (prog *Program) FLT_PRINT(arg shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_FLT_PRINT, 0, parse(arg), 0))
}

func (prog *Program) FLT_ADD(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_FLT_ADD, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) FLT_SUB(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_FLT_SUB, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) FLT_MUL(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_FLT_MUL, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) FLT_DIV(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_FLT_DIV, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) FLT_NEG(result, arg shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_FLT_NEG, parse(result), parse(arg), 0))
}

func (prog *Program) FLT_LT(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_FLT_LT, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) FLT_EQ(result, arg1, arg2 shared.Reg) {
	prog.instrs = append(prog.instrs, NewInstruction(OP_FLT_EQ, parse(result), parse(arg1), parse(arg2)))
}

func (prog *Program) LABEL(name string) {
	_, ok := prog.labelToAddress[name]
	if ok {
		panic(fmt.Errorf("label '%s' was already defined", name))
	}

	prog.labelToAddress[name] = int32(len(prog.instrs))
}

func (prog *Program) JUMP(label string) {
	instr := NewInstruction(OP_JUMP, 0, 0, 0)
	prog.instrs = append(prog.instrs, instr, Instruction(len(prog.labels)))
	prog.labels = append(prog.labels, label)
}

func (prog *Program) JUMP_IF_NOT(arg shared.Reg, label string) {
	instr := NewInstruction(OP_JUMP_IF_NOT, 0, parse(arg), 0)
	prog.instrs = append(prog.instrs, instr, Instruction(len(prog.labels)))
	prog.labels = append(prog.labels, label)
}
