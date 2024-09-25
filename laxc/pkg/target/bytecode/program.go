package bytecode

import (
	"bytes"
	"fmt"
	"math"
	"text/tabwriter"
)

type Program struct {
	instrs         []Instruction
	labels         []string
	labelToAddress map[string]int32
}

func NewProgramFromBytes(buf []byte) (prog Program, err error) {
	if len(buf)%4 != 0 {
		return prog, fmt.Errorf("the number of bytes given is not divisible by 4: %d", len(buf))
	}

	for i := 0; i < len(buf); i += 4 {
		prog.instrs = append(
			prog.instrs,
			NewInstruction(OpCode(buf[i]), Register(buf[i+1]), Register(buf[i+2]), Register(buf[i+3])),
		)
	}

	return prog, nil
}

func NewEmptyProgram() Program {
	return Program{
		labelToAddress: make(map[string]int32),
	}
}

func (prog *Program) Finalize() {
	if len(prog.labels) > 0 {
		for pc := 0; pc < len(prog.instrs); pc++ {
			instr := prog.instrs[pc]
			switch instr.OpCode() {
			case OP_JUMP, OP_JUMP_IF_NOT:
				labelIndex := prog.instrs[pc+1]
				address := prog.labelToAddress[prog.labels[labelIndex]]
				prog.instrs[pc+1] = Instruction(address)
				pc++
			}
		}

		prog.labels = nil
	}
}

func (prog *Program) Run() {
	prog.Finalize()

	var regs [8]int32

	for pc := 0; pc < len(prog.instrs); {
		instr := prog.instrs[pc]
		opCode := instr.OpCode()
		// resultReg := instr.ResultReg()
		// arg1Reg := instr.Arg1Reg()
		// arg2Reg := instr.Arg2Reg()

		// fmt.Println(opCode, resultReg, arg1Reg, arg2Reg)
		switch opCode {
		case OP_LIT:
			regs[instr.ResultReg()] = int32(prog.instrs[pc+1])
			pc += 2
		case OP_OR:
			regs[instr.ResultReg()] = regs[instr.Arg1Reg()] | regs[instr.Arg2Reg()]
			pc++
		case OP_INT_PRINT:
			println(regs[instr.Arg1Reg()])
			pc++
		case OP_INT_ADD:
			regs[instr.ResultReg()] = regs[instr.Arg1Reg()] + regs[instr.Arg2Reg()]
			pc++
		case OP_INT_SUB:
			regs[instr.ResultReg()] = regs[instr.Arg1Reg()] - regs[instr.Arg2Reg()]
			pc++
		case OP_INT_MUL:
			regs[instr.ResultReg()] = regs[instr.Arg1Reg()] * regs[instr.Arg2Reg()]
			pc++
		case OP_INT_DIV:
			regs[instr.ResultReg()] = regs[instr.Arg1Reg()] / regs[instr.Arg2Reg()]
			pc++
		case OP_INT_MOD:
			regs[instr.ResultReg()] = regs[instr.Arg1Reg()] % regs[instr.Arg2Reg()]
			pc++
		case OP_INT_NEG:
			regs[instr.ResultReg()] = -regs[instr.Arg1Reg()]
			pc++
		case OP_INT_LT:
			if regs[instr.Arg1Reg()] < regs[instr.Arg2Reg()] {
				regs[instr.ResultReg()] = 1
			} else {
				regs[instr.ResultReg()] = 0
			}
			pc++
		case OP_INT_EQ:
			if regs[instr.Arg1Reg()] == regs[instr.Arg2Reg()] {
				regs[instr.ResultReg()] = 1
			} else {
				regs[instr.ResultReg()] = 0
			}
			pc++
		case OP_INT_TO_FLOAT:
			regs[instr.ResultReg()] = int32(math.Float32bits(float32(regs[instr.Arg1Reg()])))
			pc++
		case OP_FLT_PRINT:
			fmt.Printf("%f\n", math.Float32frombits(uint32(regs[instr.Arg1Reg()])))
			pc++
		case OP_FLT_ADD:
			regs[instr.ResultReg()] = int32(
				math.Float32bits(math.Float32frombits(uint32(regs[instr.Arg1Reg()])) +
					math.Float32frombits(uint32(regs[instr.Arg2Reg()]))),
			)
			pc++
		case OP_FLT_SUB:
			regs[instr.ResultReg()] = int32(
				math.Float32bits(math.Float32frombits(uint32(regs[instr.Arg1Reg()])) -
					math.Float32frombits(uint32(regs[instr.Arg2Reg()]))),
			)
			pc++
		case OP_FLT_MUL:
			regs[instr.ResultReg()] = int32(
				math.Float32bits(math.Float32frombits(uint32(regs[instr.Arg1Reg()])) *
					math.Float32frombits(uint32(regs[instr.Arg2Reg()]))),
			)
			pc++
		case OP_FLT_DIV:
			regs[instr.ResultReg()] = int32(
				math.Float32bits(math.Float32frombits(uint32(regs[instr.Arg1Reg()])) /
					math.Float32frombits(uint32(regs[instr.Arg2Reg()]))),
			)
			pc++
		case OP_FLT_NEG:
			regs[instr.ResultReg()] = int32(-math.Float32bits(math.Float32frombits(uint32(regs[instr.Arg1Reg()]))))
			pc++
		case OP_FLT_LT:
			if math.Float32frombits(uint32(regs[instr.Arg1Reg()])) <
				math.Float32frombits(uint32(regs[instr.Arg2Reg()])) {
				regs[instr.ResultReg()] = 1
			} else {
				regs[instr.ResultReg()] = 0
			}
			pc++
		case OP_FLT_EQ:
			arg1 := math.Float32frombits(uint32(regs[instr.Arg1Reg()]))
			arg2 := math.Float32frombits(uint32(regs[instr.Arg2Reg()]))
			if arg1 == arg2 {
				regs[instr.ResultReg()] = 1
			} else {
				regs[instr.ResultReg()] = 0
			}
			pc++
		case OP_JUMP:
			pc = int(prog.instrs[pc+1])
		case OP_JUMP_IF_NOT:
			if regs[instr.Arg1Reg()] == 0 {
				pc = int(prog.instrs[pc+1])
			} else {
				pc += 2
			}
		case OP_BOOL_NOT:
			regs[instr.ResultReg()] = regs[instr.Arg1Reg()] ^ 1
			pc++
		default:
			panic(fmt.Sprintf("invalid instruction: %x", instr))
		}
	}
}

func (prog *Program) Bytes() []byte {
	data := make([]byte, 0, len(prog.instrs)*4)

	for _, instr := range prog.instrs {
		opCode := instr.OpCode()
		result := instr.ResultReg()
		arg1 := instr.Arg1Reg()
		arg2 := instr.Arg2Reg()

		data = append(
			data,
			byte(opCode),
			byte(result),
			byte(arg1),
			byte(arg2),
		)
	}

	return data
}

func (prog *Program) Disassemble() string {
	buf := bytes.NewBufferString("")
	writer := tabwriter.NewWriter(buf, 0, 2, 4, ' ', 0)

	for pc := 0; pc < len(prog.instrs); pc++ {
		instr := prog.instrs[pc]
		asm := ""
		switch instr.OpCode() {
		case OP_LIT:
			asm = fmt.Sprintf("%d\tOP_LIT\tr%d\n", pc, instr.ResultReg())
			pc++
			asm += fmt.Sprintf("%d\tDATA\t%d\n", pc, prog.instrs[pc])
		case OP_OR:
			asm = fmt.Sprintf("%d\tOP_OR\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_INT_PRINT:
			asm = fmt.Sprintf("%d\tOP_INT_PRINT\tr%d\n", pc, instr.Arg1Reg())
		case OP_INT_ADD:
			asm = fmt.Sprintf("%d\tOP_INT_ADD\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_INT_SUB:
			asm = fmt.Sprintf("%d\tOP_INT_SUB\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_INT_MUL:
			asm = fmt.Sprintf("%d\tOP_INT_MUL\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_INT_DIV:
			asm = fmt.Sprintf("%d\tOP_INT_DIV\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_INT_MOD:
			asm = fmt.Sprintf("%d\tOP_INT_MOD\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_INT_NEG:
			asm = fmt.Sprintf("%d\tOP_INT_NEG\tr%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg())
		case OP_INT_LT:
			asm = fmt.Sprintf("%d\tOP_INT_LT\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_INT_EQ:
			asm = fmt.Sprintf("%d\tOP_INT_EQ\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_INT_TO_FLOAT:
			asm = fmt.Sprintf("%d\tOP_INT_TO_FLOAT\tr%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg())
			pc++
		case OP_FLT_PRINT:
			asm = fmt.Sprintf("%d\tOP_FLT_PRINT\tr%d\n", pc, instr.Arg1Reg())
		case OP_FLT_ADD:
			asm = fmt.Sprintf("%d\tOP_FLT_ADD\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_FLT_SUB:
			asm = fmt.Sprintf("%d\tOP_FLT_SUB\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_FLT_MUL:
			asm = fmt.Sprintf("%d\tOP_FLT_MUL\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_FLT_DIV:
			asm = fmt.Sprintf("%d\tOP_FLT_DIV\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_FLT_NEG:
			asm = fmt.Sprintf("%d\tOP_FLT_NEG\tr%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg())
		case OP_FLT_LT:
			asm = fmt.Sprintf("%d\tOP_FLT_LT\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_FLT_EQ:
			asm = fmt.Sprintf("%d\tOP_FLT_EQ\tr%d,r%d,r%d\n", pc, instr.ResultReg(), instr.Arg1Reg(), instr.Arg2Reg())
		case OP_JUMP:
			asm = fmt.Sprintf("%d\tOP_JUMP\t\n", pc)
			pc++
			asm += fmt.Sprintf("%d\tTARGET\t%d\n", pc, prog.instrs[pc])
		case OP_JUMP_IF_NOT:
			asm = fmt.Sprintf("%d\tOP_JUMP_IF_NOT\tr%d\n", pc, instr.Arg1Reg())
			pc++
			asm += fmt.Sprintf("%d\tTARGET\t%d\n", pc, prog.instrs[pc])
		default:
			asm = fmt.Sprintf("%d\tUNKNOWN\t%x\n", pc, instr)
		}

		writer.Write([]byte(asm))
	}

	writer.Flush()

	return buf.String()
}
