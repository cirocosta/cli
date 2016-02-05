package aggregation

import "github.com/launchpad-project/cli/launchpad/range"

type Aggregation map[string]interface{}

type OperatorValue struct {
	Name     string      `json:"name"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value,omitempty"`
}

func New(name, field, operator string, value interface{}) *Aggregation {
	var m = make(Aggregation)

	m[field] = &OperatorValue{
		Name:     name,
		Operator: operator,
		Value:    value,
	}

	return &m
}

func Avg(name, field string) *Aggregation {
	return New(name, field, "avg", nil)
}

func Count(name, field string) *Aggregation {
	return New(name, field, "count", nil)
}

func ExtendedStats(name, field string) *Aggregation {
	return New(name, field, "extendedStats", nil)
}

func Histogram(name, field string, interval int) *Aggregation {
	return New(name, field, "histogram", interval)
}

func Max(name, field string) *Aggregation {
	return New(name, field, "max", nil)
}

func Min(name, field string) *Aggregation {
	return New(name, field, "min", nil)
}

func Missing(name, field string) *Aggregation {
	return New(name, field, "missing", nil)
}

func Stats(name, field string) *Aggregation {
	return New(name, field, "stats", nil)
}

func Sum(name, field string) *Aggregation {
	return New(name, field, "sum", nil)
}

func Terms(name, field string) *Aggregation {
	return New(name, field, "terms", nil)
}

func Distance(name, field string, location interface{}, lr ...lrange.LRange) *Aggregation {
	value := make(map[string]interface{})

	value["location"] = location
	value["ranges"] = lr

	return New(name, field, "geoDistance", value)
}

func (a *Aggregation) getFieldName() string {
	var field string
	for k := range *a {
		field = k
	}
	return field
}

func (a *Aggregation) Unit(unit string) *Aggregation {
	var i = (*a)[a.getFieldName()].(*OperatorValue).Value
	i.(map[string]interface{})["unit"] = unit
	return a
}

func (a *Aggregation) Range(args ...interface{}) *Aggregation {
	var i = (*a)[a.getFieldName()].(*OperatorValue).Value
	var ra lrange.LRange

	switch len(args) {
	case 2:
		ra = lrange.Between(args[0].(int), args[1].(int))
	default:
		ra = args[0].(lrange.LRange)
	}

	var x = i.(map[string]interface{})["ranges"]
	var r = x.([]lrange.LRange)
	i.(map[string]interface{})["ranges"] = append(r, ra)

	return a
}
