package aot

import "wasm.go/binary"

func isBrTarget(block binary.Expr) bool {
	n := len(block)
	return n > 0 && block[n-1].Opcode == 0xFF
}

func analyzeBr(code binary.Code) []binary.Instruction {
	expr := code.Expr
	brTargets := analyzeExpr(0, expr)
	for _, target := range brTargets {
		if target == 0 {
			return append(expr, binary.Instruction{Opcode: 0xFF})
		}
	}
	return expr
}

func analyzeExpr(depth uint32, expr binary.Expr) (allTargets []uint32) {
	for i, instr := range expr {
		switch instr.Opcode {
		case binary.Block, binary.Loop:
			args := instr.Args.(binary.BlockArgs)
			targets := analyzeExpr(depth+1, args.Instrs)
			for _, target := range targets {
				if target == depth+1 {
					args.Instrs = append(args.Instrs, binary.Instruction{Opcode: 0xFF})
					expr[i].Args = args // hack!
					break
				}
			}
			allTargets = append(allTargets, targets...)
		case binary.If:
			args := instr.Args.(binary.IfArgs)
			targets := analyzeExpr(depth+1, args.Instrs1)
			targets2 := analyzeExpr(depth+1, args.Instrs2)
			targets = append(targets, targets2...)
			for _, target := range targets {
				if target == depth+1 {
					args.Instrs1 = append(args.Instrs1, binary.Instruction{Opcode: 0xFF})
					expr[i].Args = args // hack!
					break
				}
			}
			allTargets = append(allTargets, targets...)
		case binary.Br:
			allTargets = append(allTargets, depth-instr.Args.(uint32))
		case binary.BrIf:
			allTargets = append(allTargets, depth-instr.Args.(uint32))
		case binary.BrTable:
			args := instr.Args.(binary.BrTableArgs)
			for _, label := range args.Labels {
				allTargets = append(allTargets, depth-label)
			}
			allTargets = append(allTargets, depth-args.Default)
		case binary.Return:
			//allTargets = append(allTargets, 0)
		}
	}
	return
}
