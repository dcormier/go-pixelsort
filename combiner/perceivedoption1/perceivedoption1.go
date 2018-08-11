// Package perceivedoption1 implements combiner.Combiner using
// http://stackoverflow.com/a/6449381/297468
// perceived, option 1
package perceivedoption1

import (
	"image/color"

	"github.com/dcormier/go-pixelsort/combiner"
)

// Combiner is an instance of this Combiner
var Combiner = New()

var _ combiner.Combiner = (*perceivedOption1)(nil)

type perceivedOption1 struct{}

// New creates a combiner.Combiner that uses perceived, option 1
func New() combiner.Combiner {
	return &perceivedOption1{}
}

func (*perceivedOption1) Name() string {
	return "perceived (option 1)"
}

func (*perceivedOption1) Combine(c color.Color) uint64 {
	r, g, b, a := c.RGBA()

	return uint64(combiner.AlphaBlend(r, a, combiner.CMax)*0.299 +
		combiner.AlphaBlend(g, a, combiner.CMax)*0.587 +
		combiner.AlphaBlend(b, a, combiner.CMax)*0.114)
}
