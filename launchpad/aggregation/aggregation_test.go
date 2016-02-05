package aggregation

import (
	"testing"

	"github.com/launchpad-project/cli/launchpad/geo"
	"github.com/launchpad-project/cli/launchpad/range"
	"github.com/launchpad-project/cli/util"
)

func TestAverage(t *testing.T) {
	var aggregation = Avg("myName", "myField")
	var want = `{"myField":{"operator":"avg","name":"myName"}}`
	util.AssertJSON(t, want, aggregation)
}

func TestCount(t *testing.T) {
	var aggregation = Count("myName", "myField")
	var want = `{"myField":{"operator":"count","name":"myName"}}`
	util.AssertJSON(t, want, aggregation)
}

func TestExtendedStats(t *testing.T) {
	var aggregation = ExtendedStats("myName", "myField")
	var want = `{"myField":{"operator":"extendedStats","name":"myName"}}`
	util.AssertJSON(t, want, aggregation)
}

func TestHistogram(t *testing.T) {
	var aggregation = Histogram("myName", "myField", 10)
	var want = `{"myField":{"operator":"histogram","name":"myName","value":10}}`
	util.AssertJSON(t, want, aggregation)
}

func TestMax(t *testing.T) {
	var aggregation = Max("myName", "myField")
	var want = `{"myField":{"operator":"max","name":"myName"}}`
	util.AssertJSON(t, want, aggregation)
}

func TestMin(t *testing.T) {
	var aggregation = Min("myName", "myField")
	var want = `{"myField":{"operator":"min","name":"myName"}}`
	util.AssertJSON(t, want, aggregation)
}

func TestMissing(t *testing.T) {
	var aggregation = Missing("myName", "myField")
	var want = `{"myField":{"operator":"missing","name":"myName"}}`
	util.AssertJSON(t, want, aggregation)
}

func TestStats(t *testing.T) {
	var aggregation = Stats("myName", "myField")
	var want = `{"myField":{"operator":"stats","name":"myName"}}`
	util.AssertJSON(t, want, aggregation)
}

func TestSum(t *testing.T) {
	var aggregation = Sum("myName", "myField")
	var want = `{"myField":{"operator":"sum","name":"myName"}}`
	util.AssertJSON(t, want, aggregation)
}

func TestTerms(t *testing.T) {
	var aggregation = Terms("myName", "myField")
	var want = `{"myField":{"operator":"terms","name":"myName"}}`
	util.AssertJSON(t, want, aggregation)
}

func TestDistance(t *testing.T) {
	var want = `{"myField":{"operator":"geoDistance","name":"myName","value":{"location":[0,0],"ranges":[{"from":0},{"to":0}]}}}`
	var distance = Distance("myName", "myField", geo.NewPoint(0, 0), lrange.From(0), lrange.To(0))
	util.AssertJSON(t, want, distance)
}

func TestDistanceWithUnit(t *testing.T) {
	var want = `{"myField":{"operator":"geoDistance","name":"myName","value":{"location":[0,0],"unit":"km","ranges":[{"from":0},{"to":0}]}}}`
	var distance = Distance("myName", "myField", geo.NewPoint(0, 0), lrange.From(0), lrange.To(0)).Unit("km")

	util.AssertJSON(t, want, distance)
}

func TestDistanceWithUnitAndExtraRange(t *testing.T) {
	var want = `{"myField":{"operator":"geoDistance","name":"myName","value":{"location":[0,0],"unit":"km","ranges":[{"from":0},{"to":0},{"to":0},{"from":0,"to":1},{"from":1}]}}}`
	var distance = Distance("myName", "myField", geo.NewPoint(0, 0), lrange.From(0), lrange.To(0))

	distance.Range(lrange.To(0))
	distance.Range(lrange.Between(0, 1))
	distance.Range(lrange.From(1)).Unit("km")

	util.AssertJSON(t, want, distance)
}

func TestDistanceWithUnitAndExtraNumericBetweenRange(t *testing.T) {
	var want = `{"myField":{"operator":"geoDistance","name":"myName","value":{"location":[0,0],"ranges":[{"from":0},{"to":0},{"from":0,"to":1}]}}}`
	var distance = Distance("myName", "myField", geo.NewPoint(0, 0), lrange.From(0), lrange.To(0))

	distance.Range(0, 1)

	util.AssertJSON(t, want, distance)
}
