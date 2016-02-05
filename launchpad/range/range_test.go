package lrange

import (
	"encoding/json"
	"testing"
)

func TestFromTo(t *testing.T) {
	var r = Between(10, 20)
	var expected = `{"from":10,"to":20}`

	bin, err := json.Marshal(r)
	var res = string(bin)

	if err != nil {
		t.Error(err)
	}

	if *r.From != 10 {
		t.Errorf("Expected range.From = %d, got %d instead", 10, &r.From)
	}

	if *r.To != 20 {
		t.Errorf("Expected range.To = %d, got %d instead", 20, &r.To)
	}

	if res != expected {
		t.Errorf("Expected JSON response %s, got %s instead", expected, res)
	}
}

func TestFromOnly(t *testing.T) {
	var r = From(10)
	var expected = `{"from":10}`

	bin, err := json.Marshal(r)
	var res = string(bin)

	if err != nil {
		t.Error(err)
	}

	if *r.From != 10 {
		t.Errorf("Expected range.From = %d, got %d instead", 10, &r.From)
	}

	if r.To != nil {
		t.Errorf("Expected range.To = %d, got %d instead", 20, &r.To)
	}

	if res != expected {
		t.Errorf("Expected JSON response %s, got %s instead", expected, res)
	}
}

func TestToOnly(t *testing.T) {
	var r = To(20)
	var expected = `{"to":20}`

	bin, err := json.Marshal(r)
	var res = string(bin)

	if err != nil {
		t.Error(err)
	}

	if r.From != nil {
		t.Errorf("Expected range.From = %d, got %d instead", 10, &r.From)
	}

	if *r.To != 20 {
		t.Errorf("Expected range.To = %d, got %d instead", 20, &r.To)
	}

	if res != expected {
		t.Errorf("Expected JSON response %s, got %s instead", expected, res)
	}
}
