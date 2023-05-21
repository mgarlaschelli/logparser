package config

import (
	"errors"
	"go/parser"
	"os"

	"github.com/mgarlaschelli/logparser/filter"
	"gopkg.in/yaml.v3"
)

type FilterConfig struct {
	Name       string
	FilterType string
	Pattern    string
}

type FiltersConfig struct {
	Filters []FilterConfig
}

func ParseConfigFile(configFileName string) ([]filter.LineFilter, error) {

	var filters []filter.LineFilter = make([]filter.LineFilter, 0)
	var err error

	f, err := os.ReadFile(configFileName)
	if err != nil {
		return filters, err
	}

	var filtersConfig FiltersConfig

	if err := yaml.Unmarshal(f, &filtersConfig); err != nil {
		return filters, err
	}

	// Generate filters
filterLoop:
	for _, filterConfig := range filtersConfig.Filters {

		switch filterConfig.FilterType {

		case filter.FilterTypeRegexp:

			var filter filter.LineFilter = &filter.RegExpFilter{
				RegExp: filterConfig.Pattern,
				Name:   filterConfig.Name,
			}

			filters = append(filters, filter)

		case filter.FilterTypeExpression:

			var filter filter.LineFilter

			filter, err = makeExpressionFilter(filterConfig.Name, filterConfig.Pattern)
			if err != nil {
				break filterLoop
			}

			filters = append(filters, filter)

		default:
			err = errors.New("filter " + filterConfig.Name + " has an invalid filter type: " + filterConfig.FilterType)
			break filterLoop
		}
	}

	return filters, err
}

func makeExpressionFilter(name, expression string) (*filter.ExpressionFilter, error) {

	tr, err := parser.ParseExpr(expression)

	if err != nil {
		return nil, err
	}

	return &filter.ExpressionFilter{
		ExpressionTree: tr,
		Name:           name,
	}, nil
}
