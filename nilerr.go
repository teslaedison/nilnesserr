// some code was copy from https://github.com/gostaticanalysis/nilerr/blob/master/nilerr.go

package nilnesserr

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ssa"
)

var errType = types.Universe.Lookup("error").Type().Underlying().(*types.Interface)

func isErrType(res ssa.Value) bool {
	return types.Implements(res.Type(), errType)
}

func isConstNil(res ssa.Value) bool {
	v, ok := res.(*ssa.Const)
	if ok && v.IsNil() {
		return true
	}
	return false
}

func getCheckedErr(binOp *ssa.BinOp) ssa.Value {
	if isErrType(binOp.X) && isConstNil(binOp.Y) {
		return binOp.X
	}
	if isErrType(binOp.Y) && isConstNil(binOp.X) {
		return binOp.Y
	}
	return nil
}

func checkIsNilnesserr(pass *analysis.Pass, b *ssa.BasicBlock, isLast bool, isNilnees func(value ssa.Value) bool, checkedErrors []ssa.Value) {
	for i := range b.Instrs {
		instr, ok := b.Instrs[i].(*ssa.Return)
		if !ok {
			continue
		}
		// skip for last Block return
		if isLast && i == len(b.Instrs)-1 {
			return
		}

		for _, res := range instr.Results {
			if !isErrType(res) || isConstNil(res) || !isNilnees(res) {
				continue
			}

			// skip for res is the last checked error
			if len(checkedErrors) > 0 && checkedErrors[len(checkedErrors)-1] == res {
				continue
			}
			// report
			pos := instr.Pos()
			if pos.IsValid() {
				pass.Report(analysis.Diagnostic{
					Pos:      pos,
					Category: linterCategory,
					Message:  linterMessage,
				})
			}
		}
	}
}
