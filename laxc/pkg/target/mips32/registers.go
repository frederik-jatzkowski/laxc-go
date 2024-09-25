package mips32

import "laxc/internal/shared"

var (
	RegFp   shared.Reg = "$fp"
	RegSp   shared.Reg = "$sp"
	RegZero shared.Reg = "$zero"
	RegV0   shared.Reg = "$v0"
	// argument registers
	RegA0 shared.Reg = "$a0"
	RegA1 shared.Reg = "$a1"
	RegA2 shared.Reg = "$a2"
	RegA3 shared.Reg = "$a3"
	// temporary registers
	RegT0 shared.Reg = "$t0"
	RegT1 shared.Reg = "$t1"
	RegT2 shared.Reg = "$t2"
	RegT3 shared.Reg = "$t3"
	RegT4 shared.Reg = "$t4"
	RegT5 shared.Reg = "$t5"
	RegT6 shared.Reg = "$t6"
	RegT7 shared.Reg = "$t7"
	// saved registers
	RegS0 shared.Reg = "$s0"
	RegS1 shared.Reg = "$s1"
	RegS2 shared.Reg = "$s2"
	RegS3 shared.Reg = "$s3"
	RegS4 shared.Reg = "$s4"
	RegS5 shared.Reg = "$s5"
	RegS6 shared.Reg = "$s6"
	RegS7 shared.Reg = "$s7"
	// working registers
	RegT8 shared.Reg = "$t8"
	RegT9 shared.Reg = "$t9"
	// floating point registers
	RegF0  shared.Reg = "$f0"
	RegF1  shared.Reg = "$f1"
	RegF12 shared.Reg = "$f12"
)

var SpillRegs = []shared.Reg{
	RegT0,
	RegT1,
	RegT2,
}

var GeneralPurposeRegs = []shared.Reg{
	RegT3,
	RegT4,
	RegT5,
	RegT6,
	RegT7,
	RegS0,
	RegS1,
	RegS2,
	RegS3,
	RegS4,
	RegS5,
	RegS6,
	RegS7,
}

var WorkRegs = []shared.Reg{
	RegT8,
	RegT9,
}

var FloatingPointRegs = []shared.Reg{
	RegF0,
	RegF1,
	RegF12,
}
