// Package perceivedoption2noalpha implements combiner.Combiner using
// http://stackoverflow.com/a/6449381/297468
// perceived, option 2; ignoring alpha values
package perceivedoption2noalpha

import (
	"image/color"
	"math"

	"github.com/dcormier/go-pixelsort/combiner"
)

func init() {
	combiner.Register(New())
}

var _ combiner.Combiner = (*perceivedOption2NoAlpha)(nil)

type perceivedOption2NoAlpha struct{}

// New creates a combiner.Combiner that uses perceived, option 1
func New() combiner.Combiner {
	return &perceivedOption2NoAlpha{}
}

func (*perceivedOption2NoAlpha) Name() string {
	return "perceived (option 2, no alpha)"
}

func (*perceivedOption2NoAlpha) Combine(c color.Color) uint64 {
	r, g, b, _ := c.RGBA()

	return uint64(math.Sqrt(math.Pow(float64(r)*0.241, 2) +
		math.Pow(float64(g)*0.691, 2) +
		math.Pow(float64(b)*0.068, 2)))
}
