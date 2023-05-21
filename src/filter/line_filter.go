package filter

const FilterTypeRegexp string = "REGEXP"
const FilterTypeExpression string = "EXPRESSION"

type LineFilter interface {
	Match(line string) (bool, error)
	GetName() string
}
