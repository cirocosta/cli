package filter

import (
	"testing"

	"github.com/launchpad-project/cli/launchpad/geo"
	"github.com/launchpad-project/cli/launchpad/range"
	"github.com/launchpad-project/cli/util"
)

func TestFilterLessThanValue(t *testing.T) {
	var want = `{"age":{"operator":"\u003e","value":12}}`
	var got = New("age", ">", 12)
	util.AssertJSON(t, want, got)
}

func TestFilterMoreThanValue(t *testing.T) {
	var want = `{"age":{"operator":"\u003c","value":0}}`
	var got = New("age", "<", 0)
	util.AssertJSON(t, want, got)
}

func TestFilterEqualValue(t *testing.T) {
	var want = `{"age":{"operator":"=","value":12}}`
	var got = Equal("age", 12)
	util.AssertJSON(t, want, got)
}

func TestFilterNotEqualValue(t *testing.T) {
	var want = `{"age":{"operator":"!=","value":12}}`
	var got = NotEqual("age", 12)
	util.AssertJSON(t, want, got)
}

func TestFilterGreaterValue(t *testing.T) {
	var want = `{"age":{"operator":">","value":12}}`
	var got = Gt("age", 12)
	util.AssertJSON(t, want, got)
}

func TestFilterGreaterOrEqualValue(t *testing.T) {
	var want = `{"age":{"operator":">=","value":12}}`
	var got = Gte("age", 12)
	util.AssertJSON(t, want, got)
}

func TestFilterLessValue(t *testing.T) {
	var want = `{"age":{"operator":"<","value":12}}`
	var got = Lt("age", 12)
	util.AssertJSON(t, want, got)
}

func TestFilterLessOrEqualValue(t *testing.T) {
	var want = `{"age":{"operator":"=\u003c","value":12}}`
	var got = Lte("age", 12)
	util.AssertJSON(t, want, got)
}

func TestFilterRegex(t *testing.T) {
	var want = `{"age":{"operator":"~","value":12}}`
	var got = Regex("age", 12)
	util.AssertJSON(t, want, got)
}

func TestFilterNone(t *testing.T) {
	var want = `{"age":{"operator":"none","value":12}}`
	var got = None("age", 12)
	util.AssertJSON(t, want, got)
}

func TestFilterNilValue(t *testing.T) {
	var want = `{"age":{"operator":"="}}`
	var got = Equal("age", nil)
	util.AssertJSON(t, want, got)
}

func TestComposingAndFilters(t *testing.T) {
	var want = `{"and":[{"age":{"operator":"\u003e","value":12}},{"age":{"operator":"\u003c","value":15}},{"name":{"operator":"=","value":"foo"}}]}`
	var got = And(New("age", ">", 12),
		New("age", "<", 15),
		Equal("name", "foo"),
	)
	util.AssertJSON(t, want, got)
}

func TestComposingOrFilters(t *testing.T) {
	var want = `{"or":[{"age":{"operator":"\u003e","value":12}},{"age":{"operator":"\u003c","value":15}},{"name":{"operator":"=","value":"foo"}}]}`
	var got = Or(New("age", ">", 12),
		New("age", "<", 15),
		Equal("name", "foo"),
	)
	util.AssertJSON(t, want, got)
}

func TestNegation(t *testing.T) {
	var want = `{"not":{"age":{"operator":"\u003e","value":12}}}`
	var ageFilter = New("age", ">", 12)
	var got = ageFilter.UnaryAdd("not")
	util.AssertJSON(t, want, got)
}

func TestAny(t *testing.T) {
	var want = `{"age":{"operator":"any","value":[12,21,25]}}`
	var got = Any("age", 12, 21, 25)
	util.AssertJSON(t, want, got)
}

func TestExists(t *testing.T) {
	var want = `{"age":{"operator":"exists"}}`
	var got = Exists("age")
	util.AssertJSON(t, want, got)
}

func TestMissing(t *testing.T) {
	var want = `{"age":{"operator":"missing"}}`
	var got = Missing("age")
	util.AssertJSON(t, want, got)
}

func TestMatch(t *testing.T) {
	var want = `{"*":{"operator":"match","value":"foo"}}`
	var got = Match("foo")
	util.AssertJSON(t, want, got)
}

func TestPhrase(t *testing.T) {
	var want = `{"*":{"operator":"phrase","value":"foo"}}`
	var got = Phrase("foo")
	util.AssertJSON(t, want, got)
}

func TestPrefix(t *testing.T) {
	var want = `{"*":{"operator":"prefix","value":"myPrefix"}}`
	var got = Prefix("myPrefix")
	util.AssertJSON(t, want, got)
}

func TestFuzyQuery(t *testing.T) {
	var want = `{"*":{"operator":"fuzzy","value":{"query":"foo"}}}`
	var got = Fuzzy("foo", nil, nil)
	util.AssertJSON(t, want, got)
}

func TestFuzyFieldAndQuery(t *testing.T) {
	var want = `{"name":{"operator":"fuzzy","value":{"query":"foo"}}}`
	var got = Fuzzy("name", "foo", nil)
	util.AssertJSON(t, want, got)
}

func TestFuzyFieldAndFuziness(t *testing.T) {
	var want = `{"*":{"operator":"fuzzy","value":{"query":"foo","fuzziness":0.8}}}`
	var got = Fuzzy("foo", nil, 0.8)
	util.AssertJSON(t, want, got)
}

func TestSimilarQuery(t *testing.T) {
	var want = `{"*":{"operator":"similar","value":{"query":"foo"}}}`
	var got = Similar("foo", nil, nil)
	util.AssertJSON(t, want, got)
}

func TestSimilarFieldAndQuery(t *testing.T) {
	var want = `{"name":{"operator":"similar","value":{"query":"foo"}}}`
	var got = Similar("name", "foo", nil)
	util.AssertJSON(t, want, got)
}

func TestSimilarFieldAndFuziness(t *testing.T) {
	var want = `{"*":{"operator":"similar","value":{"query":"foo","fuzziness":0.8}}}`
	var got = Similar("foo", nil, 0.8)
	util.AssertJSON(t, want, got)
}

func TestDistanceCircle(t *testing.T) {
	var want = `{"point":{"operator":"gd","value":{"location":[0,0],"max":"2km"}}}`
	var got = Distance("point", geo.NewCircle(geo.NewPoint(0, 0), "2km"), nil)
	util.AssertJSON(t, want, got)
}

func TestDistancePointFrom(t *testing.T) {
	var want = `{"point":{"operator":"gd","value":{"location":[0,0],"min":1}}}`
	var got = Distance("point", geo.NewPoint(0, 0), lrange.From(1))
	util.AssertJSON(t, want, got)
}

func TestDistancePointToImplicit(t *testing.T) {
	var want = `{"point":{"operator":"gd","value":{"location":[0,0],"max":2}}}`
	var got = Distance("point", geo.NewPoint(0, 0), 2)
	util.AssertJSON(t, want, got)
}

func TestDistancePointTo(t *testing.T) {
	var want = `{"point":{"operator":"gd","value":{"location":[0,0],"max":2}}}`
	var got = Distance("point", geo.NewPoint(0, 0), lrange.To(2))
	util.AssertJSON(t, want, got)
}

func TestDistancePointRange(t *testing.T) {
	var want = `{"point":{"operator":"gd","value":{"location":[0,0],"min":1,"max":2}}}`
	var got = Distance("point", geo.NewPoint(0, 0), lrange.Between(1, 2))
	util.AssertJSON(t, want, got)
}

func TestRange(t *testing.T) {
	var want = `{"age":{"operator":"range","value":{"from":12,"to":15}}}`
	var got = Range("age", 12, 15)
	util.AssertJSON(t, want, got)
}

func TestRangeLRange(t *testing.T) {
	var want = `{"age":{"operator":"range","value":{"from":12,"to":15}}}`
	var got = Range("age", lrange.Between(12, 15))
	util.AssertJSON(t, want, got)
}

func TestShapes(t *testing.T) {
	var want = `{"xshape":{"operator":"gs","value":{"type":"geometrycollection",` +
		`"geometries":[{"type":"circle","coordinates":[0,0],"radius":"2km"},` +
		`{"type":"envelope","coordinates":[[20,0],[0,20]]}]}}}`

	var got = Shape("xshape", geo.NewCircle(geo.NewPoint(0, 0), "2km"),
		geo.NewBoundingBox(geo.NewPoint(20, 0), geo.NewPoint(0, 20)))
	util.AssertJSON(t, want, got)
}

func TestPolygon(t *testing.T) {
	var want = `{"xshape":{"operator":"gp","value":[[10,0],[20,0],[15,10]]}}`
	var got = Polygon("xshape", geo.NewPoint(10, 0), geo.NewPoint(20, 0), geo.NewPoint(15, 10))
	util.AssertJSON(t, want, got)
}

func TestBoundingBox(t *testing.T) {
	var want = `{"shape":{"operator":"gp","value":[[20,0],[0,20]]}}`
	var got = BoundingBox("shape", geo.NewBoundingBox(geo.NewPoint(20, 0), geo.NewPoint(0, 20)))
	util.AssertJSON(t, want, got)
}

func TestBoundingBoxByGeoPoints(t *testing.T) {
	var want = `{"xshape":{"operator":"gp","value":[[20,0],[0,20]]}}`
	var got = BoundingBox("xshape", geo.NewPoint(20, 0), geo.NewPoint(0, 20))
	util.AssertJSON(t, want, got)
}
