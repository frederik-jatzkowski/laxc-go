package intermediate

import (
	"fmt"
	"io"
	"laxc/internal/graph"
	"laxc/internal/shared"
	"laxc/pkg/target/mips32"
	"slices"
)

type Function struct {
	Id               int
	program          *Program
	entry            *BasicBlock
	blocks           []*BasicBlock
	symRegCount      shared.SymReg
	localSymVarCount shared.LocalSymVar
	labelCount       int
}

func (function *Function) AddBasicBlock(comment string) *BasicBlock {
	block := &BasicBlock{
		id:       len(function.blocks),
		comment:  comment,
		function: function,
	}

	function.blocks = append(function.blocks, block)

	return function.blocks[len(function.blocks)-1]
}

func (function *Function) Entry() *BasicBlock {
	return function.entry
}

func (function *Function) NextSymReg() shared.SymReg {
	result := function.symRegCount
	function.symRegCount++

	return result
}

func (function *Function) NextLocalSymVar() shared.LocalSymVar {
	result := function.localSymVarCount
	function.localSymVarCount++

	return result
}

func (function *Function) NextLabel() int {
	label := function.labelCount
	function.labelCount++

	return label
}

func (function *Function) WriteTo(writer io.Writer) {
	writer.Write([]byte(fmt.Sprintf("func f%d() {\n", function.Id)))
	for _, block := range function.blocks {
		block.WriteTo(writer)
	}
	writer.Write([]byte("}\n"))
}

func (function *Function) localSymVars() (result []shared.LocalSymVar) {
	for _, block := range function.blocks {
		for _, instruction := range block.instructions {
			result = append(result, instruction.UsedLocalSymVars()...)
		}
	}

	return result
}

func (function *Function) Mips32Program(mips32Prog *mips32.Program) {
	localSymVarAllocs, symRegAllocs := function.AllocateGreedyWithRecycling(mips32.GeneralPurposeRegs)
	for _, block := range function.blocks {
		block.Mips32Program(localSymVarAllocs, symRegAllocs, mips32Prog)
	}
}

func (function *Function) Optimize() {
	function.substituteInstructions()
	function.removeUnusedInstructions()
}

func (function *Function) substituteInstructions() {
	assignments := make(map[shared.SymReg]Instruction)
	for _, block := range function.blocks {
		for _, instruction := range block.instructions {
			write, ok := instruction.(WriteInstruction)
			if ok {
				assignments[write.Result()] = instruction
			}
		}
	}

	for _, block := range function.blocks {
		for i, instruction := range block.instructions {
			substituted, _ := instruction.Optimize(assignments)
			block.instructions[i] = substituted

			write, ok := substituted.(WriteInstruction)
			if ok {
				assignments[write.Result()] = substituted
			}
		}
	}
}

func (function *Function) sortBasicBlocksTopologically() {
	symRegDefs := make(map[shared.SymReg]*BasicBlock)
	for _, block := range function.blocks {
		for _, instruction := range block.instructions {
			write, ok := instruction.(WriteInstruction)
			if ok {
				symRegDefs[write.Result()] = block
			}
		}
	}

	cfg := graph.NewDiGraph[int]()
	for _, block := range function.blocks {
		for _, instruction := range block.instructions {
			usedSymRegs := instruction.UsedSymRegs()
			for _, symReg := range usedSymRegs {
				defBlock := symRegDefs[symReg]
				if defBlock.id != block.id {
					cfg.AddEdge(defBlock.id, block.id)
				}
			}

			switch branch := instruction.(type) {
			case *jump:
				cfg.AddEdge(block.id, branch.Target.id)
			case *branchIf:
				cfg.AddEdge(block.id, branch.target.id)
			case *branchIfNot:
				cfg.AddEdge(block.id, branch.target.id)
			}
		}
	}

	order, err := cfg.TopologicalSort()
	if err != nil {
		panic(err)
	}

	orderMap := make(map[int]int, len(order))
	for i, id := range order {
		orderMap[id] = i
	}

	slices.SortFunc(function.blocks, func(a, b *BasicBlock) int {
		return orderMap[a.id] - orderMap[b.id]
	})
}

func (function *Function) removeUnusedInstructions() {
	dataflowGraph := graph.NewDiGraph[shared.SymReg]()
	roots := make([]shared.SymReg, 0)
	for _, block := range function.blocks {
		for _, instruction := range block.instructions {
			for _, used := range instruction.UsedSymRegs() {
				write, ok := instruction.(WriteInstruction)
				if ok {
					dataflowGraph.AddEdge(write.Result(), used)
				}

				if instruction.HasSideEffects() {
					roots = append(roots, used)
				}
			}
		}
	}

	reachabilityInfo := make(map[shared.SymReg]bool)
	for _, root := range roots {
		dataflowGraph.Dfs(root, reachabilityInfo)
	}

	for _, block := range function.blocks {
		cleaned := make([]Instruction, 0, len(block.instructions))

		for _, instruction := range block.instructions {

			if instruction.HasSideEffects() {
				cleaned = append(cleaned, instruction)
				continue
			}

			write, ok := instruction.(WriteInstruction)
			if !ok {
				break
			}

			_, isReachable := reachabilityInfo[write.Result()]
			if !isReachable {
				break
			}

			cleaned = append(cleaned, instruction)

		}

		block.instructions = cleaned
	}
}
