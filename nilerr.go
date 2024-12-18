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

func getCheckedErrValue(binOp *ssa.BinOp) ssa.Value {
	if isErrType(binOp.X) && isConstNil(binOp.Y) {
		return binOp.X
	}
	if isErrType(binOp.Y) && isConstNil(binOp.X) {
		return binOp.Y
	}
	return nil
}

type checkedErr struct {
	err     ssa.Value
	nilness nilness
}

func getLatestNonnilValue(errors []checkedErr, res ssa.Value) ssa.Value {
	if len(errors) == 0 {
		return nil
	}

	for j := len(errors) - 1; j >= 0; j-- {
		last := errors[j]
		if last.err == res {
			return nil
		} else {
			if last.nilness == isnonnil {
				return last.err
			}
		}
	}

	return nil
}

func checkNilnesserr(pass *analysis.Pass, b *ssa.BasicBlock, errors []checkedErr, isNilnees func(value ssa.Value) bool) {
	for i := range b.Instrs {
		instr, ok := b.Instrs[i].(*ssa.Return)
		if !ok {
			continue
		}

		for _, res := range instr.Results {
			if !isErrType(res) || isConstNil(res) || !isNilnees(res) {
				continue
			}
			// check the latestValue error that is isnonnil
			latestValue := getLatestNonnilValue(errors, res)
			if latestValue == nil {
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
