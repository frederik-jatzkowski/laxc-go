package bytecode

type Instruction uint32

func NewInstruction(op OpCode, result, arg1, arg2 Register) Instruction {
	opBits := (uint32(op) << 24)
	resultBits := (uint32(result) << 16)
	arg1Bits := (uint32(arg1) << 8)
	arg2Bits := (uint32(arg2))

	return Instruction(opBits | resultBits | arg1Bits | arg2Bits)
}

func (instr Instruction) OpCode() OpCode {
	return OpCode((instr >> 24) & 0xff)
}

func (instr Instruction) ResultReg() Register {
	return Register(((instr) >> 16) & 0xff)
}

func (instr Instruction) Arg1Reg() Register {
	return Register(((instr) >> 8) & 0xff)
}

func (instr Instruction) Arg2Reg() Register {
	return Register((instr) & 0xff)
}
