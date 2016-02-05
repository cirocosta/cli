package filter

import (
	"github.com/launchpad-project/cli/launchpad/geo"
	"github.com/launchpad-project/cli/launchpad/range"
)

type Filter map[string]interface{}

type OperatorValue struct {
	Operator string      `json:"operator"`
	Value    interface{} `json:"value,omitempty"`
}

func (f *Filter) UnaryAdd(filter string) Filter {
	m := make(Filter)
	m[filter] = f
	return m
}

func New(field, operator string, value interface{}) Filter {
	var m = make(Filter)

	m[field] = &OperatorValue{
		Operator: operator,
		Value:    value,
	}

	return m
}

func Equal(field string, value interface{}) Filter {
	return New(field, "=", value)
}

func NotEqual(field string, value interface{}) Filter {
	return New(field, "!=", value)
}

func Gt(field string, value interface{}) Filter {
	return New(field, ">", value)
}

func Gte(field string, value interface{}) Filter {
	return New(field, ">=", value)
}

func Lt(field string, value interface{}) Filter {
	return New(field, "<", value)
}

func Lte(field string, value interface{}) Filter {
	return New(field, "=<", value)
}

func Regex(field string, value interface{}) Filter {
	return New(field, "~", value)
}

func None(field string, value interface{}) Filter {
	return New(field, "none", value)
}

func Any(field string, value ...interface{}) Filter {
	return New(field, "any", value)
}

func Add(operator string, filter ...Filter) Filter {
	m := make(Filter)
	m[operator] = filter
	return m
}

func And(filter ...Filter) Filter {
	return Add("and", filter...)
}

func Or(filter ...Filter) Filter {
	return Add("or", filter...)
}

func Exists(field string) Filter {
	return New(field, "exists", nil)
}

func Missing(field string) Filter {
	return New(field, "missing", nil)
}

func Match(value string) Filter {
	return New("*", "match", value)
}

func Phrase(value string) Filter {
	return New("*", "phrase", value)
}

func Prefix(value string) Filter {
	return New("*", "prefix", value)
}

func q(qType, fieldOrQuery string, query interface{}, fuzziness interface{}) Filter {
	var field string

	switch query {
	case nil:
		field = "*"
		query = fieldOrQuery
	default:
		field = fieldOrQuery
	}

	value := make(map[string]interface{})

	value["query"] = query

	if fuzziness != nil {
		value["fuzziness"] = fuzziness
	}

	return New(field, qType, value)
}

func Fuzzy(fieldOrQuery string, query interface{}, fuzziness interface{}) Filter {
	return q("fuzzy", fieldOrQuery, query, fuzziness)
}

func Similar(fieldOrQuery string, query interface{}, fuzziness interface{}) Filter {
	return q("similar", fieldOrQuery, query, fuzziness)
}

func Distance(field string, location interface{}, lr interface{}) Filter {
	value := make(map[string]interface{})

	switch location.(type) {
	case geo.Circle:
		geoCircles := location.(geo.Circle)
		value["location"] = geoCircles.Coordinates
		value["max"] = geoCircles.Radius
	default:
		value["location"] = location.(geo.Point)

		switch lr.(type) {
		case lrange.LRange:
			r := lr.(lrange.LRange)

			if r.From != nil {
				value["min"] = r.From
			}

			if r.To != nil {
				value["max"] = r.To
			}
		case int:
			value["max"] = lr.(int)
		}

	}

	return New(field, "gd", value)
}

func Range(field string, args ...interface{}) Filter {
	if len(args) == 1 {
		return New(field, "range", args[0].(lrange.LRange))
	}

	return New(field, "range", lrange.Between(args[0].(int), args[1].(int)))
}

func Shape(field string, shapes ...interface{}) Filter {
	value := make(map[string]interface{})
	value["type"] = "geometrycollection"
	value["geometries"] = shapes
	return New(field, "gs", value)
}

func Polygon(field string, points ...geo.Point) Filter {
	return New(field, "gp", points)
}

func BoundingBox(field string, boxOrUpperLeft interface{}, optLowerRight ...interface{}) Filter {
	var coords []geo.Point
	switch boxOrUpperLeft.(type) { // or len(optLowerRight)
	case geo.BoundingBox: // or 0
		coords = boxOrUpperLeft.(geo.BoundingBox).Coordinates
	default:
		coords = []geo.Point{boxOrUpperLeft.(geo.Point), optLowerRight[0].(geo.Point)}
	}

	return Polygon(field, coords...)
}
