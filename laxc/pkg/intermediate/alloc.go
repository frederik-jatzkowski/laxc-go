package intermediate

import (
	"fmt"
	"laxc/internal/graph"
	"laxc/internal/shared"
	"laxc/pkg/target/mips32"
)

type Allocation struct {
	Reg       shared.Reg
	MemLoc    int
	IsSpilled bool
}

func (alloc Allocation) String() string {
	if alloc.IsSpilled {
		return fmt.Sprintf("spill\t%d", alloc.MemLoc)
	}

	return fmt.Sprintf("reg\t%s", alloc.Reg)
}

func (function *Function) AllocateGreedyWithRecycling(registers []shared.Reg) (
	localSymVarAllocs map[shared.LocalSymVar]int32,
	symRegAllocs map[shared.SymReg]Allocation,
) {
	function.ResolvePhiFunctions()

	localSymVarAllocs = make(map[shared.LocalSymVar]int32)
	for _, localSymVar := range function.localSymVars() {
		_, alreadyAllocated := localSymVarAllocs[localSymVar]
		if alreadyAllocated {
			continue
		}

		localSymVarAllocs[localSymVar] = int32(len(localSymVarAllocs)) * 4
	}

	lastUsages := make(map[shared.SymReg]int)
	line := 0
	for _, block := range function.blocks {
		for _, instruction := range block.instructions {
			for _, used := range instruction.UsedSymRegs() {
				lastUsages[used] = max(lastUsages[used], line)
			}
			line++
		}
	}

	symRegAllocs = make(map[shared.SymReg]Allocation)
	usedRegisters := make(map[shared.Reg]bool, len(mips32.GeneralPurposeRegs))
	usedSpillAddresses := make(map[int]bool)
	nextSpillAddress := len(localSymVarAllocs) * 4

	line = 0
	for _, block := range function.blocks {
		for _, instruction := range block.instructions {
			write, ok := instruction.(WriteInstruction)
			if ok {
				result := write.Result()

				_, isAllocated := symRegAllocs[result]
				if isAllocated {
					continue
				}

				var freeRegister shared.Option[shared.Reg]
				for _, reg := range registers {
					if !usedRegisters[reg] {
						usedRegisters[reg] = true
						freeRegister = shared.Some(reg)

						break
					}
				}

				if freeRegister.IsSet {
					symRegAllocs[result] = Allocation{
						Reg:       freeRegister.Value,
						IsSpilled: false,
					}
				} else {
					var freeSpillAddress shared.Option[int]
					for address, used := range usedSpillAddresses {
						if !used {
							usedSpillAddresses[address] = true
							freeSpillAddress = shared.Some(address)

							break
						}
					}

					if freeSpillAddress.IsSet {
						symRegAllocs[result] = Allocation{
							MemLoc:    freeSpillAddress.Value,
							IsSpilled: true,
						}
					} else {
						symRegAllocs[result] = Allocation{
							MemLoc:    nextSpillAddress,
							IsSpilled: true,
						}
						usedSpillAddresses[nextSpillAddress] = true

						nextSpillAddress += 4
					}
				}
			}

			for reg, lastUsage := range lastUsages {
				if lastUsage == line {
					alloc := symRegAllocs[reg]
					if !alloc.IsSpilled {
						delete(usedRegisters, alloc.Reg)
					} else {
						usedSpillAddresses[alloc.MemLoc] = false
					}
				}
			}

			line++
		}
	}

	// resolve phi functions
	for _, block := range function.blocks {
		for _, instruction := range block.instructions {
			switch op := instruction.(type) {
			case *binOp:
				if op.name == "Phi" {
					symRegAllocs[op.arg1] = symRegAllocs[op.result]
					symRegAllocs[op.arg2] = symRegAllocs[op.result]
				}
			}
		}
	}

	return localSymVarAllocs, symRegAllocs
}
func (function *Function) ResolvePhiFunctions() {
	// build graph of phi value usages
	phiGraph := graph.NewDiGraph[shared.SymReg]()
	for _, block := range function.blocks {
		for _, instruction := range block.instructions {
			switch op := instruction.(type) {
			case *phi:
				phiGraph.AddEdge(op.arg1, op.result)
				phiGraph.AddEdge(op.arg2, op.result)
			}
		}
	}

	// use outermost symreg within blocks
	for _, block := range function.blocks {
		for _, instruction := range block.instructions {
			write, ok := instruction.(WriteInstruction)
			if !ok {
				continue
			}

			result := write.Result()
			if phiGraph.Has(result) {
				last := phiGraph.Dfs(result, make(map[shared.SymReg]bool))
				write.SetResult(last)
			}
		}
	}
}
