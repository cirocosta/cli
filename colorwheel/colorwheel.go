package colorwheel

import "github.com/wedeploy/cli/color"

// Wheel for sequentially repeatable colors for coloring messages related
// grouping by a given id - color relation
type Wheel struct {
	palette [][]color.Attribute
	hm      map[string][]color.Attribute
	next    int
}

// New create a color Wheel
func New(palette [][]color.Attribute) Wheel {
	return Wheel{
		palette: palette,
	}
}

// Get a color for a given id
func (w *Wheel) Get(id string) []color.Attribute {
	if w.hm == nil {
		w.hm = map[string][]color.Attribute{}
	}

	var _, ok = w.hm[id]

	if ok {
		return w.hm[id]
	}

	w.hm[id] = w.palette[w.next]
	w.nextColor()

	return w.hm[id]
}

func (w *Wheel) nextColor() {
	if w.next == len(w.palette)-1 {
		w.next = 0
	} else {
		w.next++
	}
}
