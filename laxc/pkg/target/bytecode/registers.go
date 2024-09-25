package bytecode

import (
	"fmt"
	"laxc/internal/shared"
	"strconv"
)

type Register byte

func parse(reg shared.Reg) (parsed Register) {
	_, err := fmt.Sscanf(string(reg), "r%d", &parsed)
	if err != nil {
		panic(err)
	}

	return parsed
}

var GeneralPurposeRegs = []shared.Reg{}

func init() {
	for i := 0; i < 255; i++ {
		GeneralPurposeRegs = append(GeneralPurposeRegs, shared.Reg("r"+strconv.Itoa(i)))
	}
}
