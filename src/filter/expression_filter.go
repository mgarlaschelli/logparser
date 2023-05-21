package filter

import (
	"errors"
	"go/ast"
	"go/token"
	"strings"
)

type ExpressionFilter struct {
	ExpressionTree ast.Node
	Name           string
}

func (f *ExpressionFilter) GetName() string {
	return f.Name
}

func (f *ExpressionFilter) Match(line string) (bool, error) {

	match, err := f.evalNode(f.ExpressionTree, line)

	return match, err
}

func (f *ExpressionFilter) evalNode(node ast.Node, line string) (bool, error) {

	var err error
	var bVal, bVal1, bVal2 bool

	switch nod := node.(type) {

	case *ast.BasicLit: // string

		// leading and trailing " must be removed
		escapedValue := nod.Value

		if strings.HasPrefix(escapedValue, `"`) && strings.HasSuffix(escapedValue, `"`) {
			escapedValue = escapedValue[1 : len(escapedValue)-1]
		}

		bVal = strings.Contains(line, escapedValue)

	case *ast.ParenExpr:
		bVal, err = f.evalNode(nod.X, line)

	case *ast.BinaryExpr:
		bVal1, err = f.evalNode(nod.X, line)
		if err != nil {
			return false, err
		}

		bVal2, err = f.evalNode(nod.Y, line)
		if err != nil {
			return false, err
		}

		bVal, err = f.operate(bVal1, bVal2, nod.Op)

	case *ast.UnaryExpr:
		bVal1, err = f.evalNode(nod.X, line)
		if err != nil {
			return false, err
		}

		bVal, err = f.operate(bVal1, false, nod.Op)
	}

	return bVal, err
}

func (f *ExpressionFilter) operate(value1, value2 bool, op token.Token) (bool, error) {

	var err error
	var bVal bool

	switch op {
	case token.LAND:
		bVal = value1 && value2

	case token.LOR:
		bVal = value1 || value2

	case token.NOT:
		bVal = !value1

	default:
		bVal = false
		err = errors.New("operator " + op.String() + " not supported ")
	}

	return bVal, err
}
