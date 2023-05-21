package filter

import "regexp"

type RegExpFilter struct {
	RegExp string
	Name   string
}

func (f *RegExpFilter) Match(line string) (bool, error) {
	return regexp.MatchString(f.RegExp, line)
}

func (f *RegExpFilter) GetName() string {
	return f.Name
}
