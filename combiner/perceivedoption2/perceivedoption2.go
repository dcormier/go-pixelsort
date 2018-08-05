// Package perceivedoption2 implements combiner.Combiner using
// http://stackoverflow.com/a/6449381/297468
// perceived, option 2
package perceivedoption2

import (
	"image/color"
	"math"

	"github.com/dcormier/go-pixelsort/combiner"
)

func init() {
	combiner.Register(New())
}

var _ combiner.Combiner = (*perceivedOption2)(nil)

type perceivedOption2 struct{}

// New creates a combiner.Combiner that uses perceived, option 2
func New() combiner.Combiner {
	return &perceivedOption2{}
}

func (*perceivedOption2) Name() string {
	return "perceived (option 2)"
}

func (*perceivedOption2) Combine(c color.Color) uint64 {
	r, g, b, a := c.RGBA()

	return uint64(math.Sqrt(math.Pow(combiner.AlphaBlend(r, a, combiner.CMax)*0.241, 2) +
		math.Pow(combiner.AlphaBlend(g, a, combiner.CMax)*0.691, 2) +
		math.Pow(combiner.AlphaBlend(b, a, combiner.CMax)*0.068, 2)))
}
