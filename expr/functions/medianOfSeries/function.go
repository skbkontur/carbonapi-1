package medianOfSeries

import (
	"github.com/bookingcom/carbonapi/expr/helper"
	"github.com/bookingcom/carbonapi/expr/interfaces"
	"github.com/bookingcom/carbonapi/expr/types"
	"github.com/bookingcom/carbonapi/pkg/parser"
)

type medianOfSeries struct {
	interfaces.FunctionBase
}

func GetOrder() interfaces.Order {
	return interfaces.Any
}

func New(configFile string) []interfaces.FunctionMetadata {
	res := make([]interfaces.FunctionMetadata, 0)
	f := &medianOfSeries{}
	functions := []string{"median"}
	for _, n := range functions {
		res = append(res, interfaces.FunctionMetadata{Name: n, F: f})
	}
	return res
}

// medianOfSeries(seriesList, n, interpolate=False)
func (f *medianOfSeries) Do(e parser.Expr, from, until int32, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	// TODO(dgryski): make sure the arrays are all the same 'size'
	args, err := helper.GetSeriesArg(e.Args()[0], from, until, values)
	if err != nil {
		return nil, err
	}

	interpolate, err := e.GetBoolNamedOrPosArgDefault("interpolate", 1, false)
	if err != nil {
		return nil, err
	}

	return helper.AggregateSeries(e, args, func(values []float64) float64 {
		return helper.Percentile(values, 50.0, interpolate)
	})
}

// Description is auto-generated description, based on output of https://github.com/graphite-project/graphite-web
func (f *medianOfSeries) Description() map[string]types.FunctionDescription {
	return map[string]types.FunctionDescription{
		"medianOfSeries": {
			Description: "medianOfSeries returns a single series which is composed of the 50-percentile\nvalues taken across a wildcard series at each point. Unless `interpolate` is\nset to True, percentile values are actual values contained in one of the\nsupplied series.",
			Function:    "medianOfSeries(seriesList, interpolate=False)",
			Group:       "Combine",
			Module:      "graphite.render.functions",
			Name:        "medianOfSeries",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Default: types.NewSuggestion(false),
					Name:    "interpolate",
					Type:    types.Boolean,
				},
			},
		},
	}
}
