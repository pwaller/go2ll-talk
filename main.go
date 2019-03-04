package main

import (
	"fmt"
	"os"
	
	goconstant "go/constant"

	ir "github.com/llir/llvm/ir"
	irconstant "github.com/llir/llvm/ir/constant"
	irtypes "github.com/llir/llvm/ir/types"
	irvalue "github.com/llir/llvm/ir/value"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	cfg := packages.Config{Mode: packages.LoadAllSyntax}
	pkgs, err := packages.Load(&cfg, "./hello-world")
	if err != nil || packages.PrintErrors(pkgs) > 0 {
		panic(err)
	}

	prog, ssaPkgs := ssautil.AllPackages(pkgs, 0)
	prog.Build()

	t := translator{
		goToIR: map[ssa.Value]irvalue.Value{},
	}

	t.printf = t.m.NewFunc("printf", irtypes.Void)
	t.printf.Sig.Variadic = true

	for _, ssaPkg := range ssaPkgs {
		t.translatePkg(ssaPkg)
	}

	os.Stdout.WriteString(t.m.String())
}

type translator struct {
	m      ir.Module
	printf *ir.Func

	goToIR map[ssa.Value]irvalue.Value
}

func (t *translator) translatePkg(ssaPkg *ssa.Package) {
	for _, ssaM := range ssaPkg.Members {
		switch ssaM := ssaM.(type) {
		case *ssa.Function:
			t.translateFunction(ssaM)
		}
	}
}

func (t *translator) translateFunction(goFn *ssa.Function) {
	if goFn.Name() != "main" {
		return
	}

	irFn := t.m.NewFunc(goFn.Name(), irtypes.Void)

	for _, goBlock := range goFn.Blocks {
		irBlock := irFn.NewBlock(goBlock.String())
		for _, goInst := range goBlock.Instrs {
			t.translateInst(irBlock, goInst)
		}
	}
}

func (t *translator) translateValue(goValue ssa.Value) irvalue.Value {
	irValue, ok := t.goToIR[goValue]
	if ok {
		return irValue
	}

	switch goValue := goValue.(type) {
	case *ssa.Const:
		switch goValue.Value.Kind() {
		case goconstant.Int:
			return irconstant.NewInt(irtypes.I64, goValue.Int64())
		case goconstant.String:
			strVal := goconstant.StringVal(goValue.Value)
			return t.m.NewGlobalDef("$const_str", irconstant.NewCharArrayFromString(strVal))
		default:
			panic("unimplemented constant")
		}
	case *ssa.Builtin:
		switch goValue.Name() {
		case "println":
			return t.printf
		default:
			panic("unimplemented builtin")
		}
	default:
		panic(fmt.Errorf("unimplemented translateValue: %T: %v", goValue, goValue))
	}
}

func (t *translator) translateInst(irBlock *ir.Block, goInst ssa.Instruction) {
	switch goInst := goInst.(type) {
	case *ssa.Call:
		irCallee := t.translateValue(goInst.Call.Value)
		var irArgs []irvalue.Value
		for _, goArg := range goInst.Call.Args {
			irArgs = append(irArgs, t.translateValue(goArg))
		}
		t.goToIR[goInst] = irBlock.NewCall(irCallee, irArgs...)

	case *ssa.BinOp:
		irX := t.translateValue(goInst.X)
		irY := t.translateValue(goInst.Y)
		t.goToIR[goInst] = irBlock.NewAdd(irX, irY)

	case *ssa.Return:
		irBlock.NewRet(nil) // TODO(pwaller): Implement returning values.

	default:
		panic("unimplemented instruction")
	}
}
