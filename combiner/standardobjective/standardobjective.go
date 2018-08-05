// Package standardobjective implements combiner.Combiner based on
// http://stackoverflow.com/a/6449381/297468 standard, objective
package standardobjective

import (
	"image/color"

	"github.com/dcormier/go-pixelsort/combiner"
)

func init() {
	combiner.Register(New())
}

var _ combiner.Combiner = (*standardObjective)(nil)

type standardObjective struct{}

// New creates a new combiner.Combiner that uses standard, objective processing
func New() combiner.Combiner {
	return &standardObjective{}
}

func (*standardObjective) Name() string {
	return "standard objective"
}

func (*standardObjective) Combine(c color.Color) uint64 {
	r, g, b, a := c.RGBA()

	return uint64(combiner.AlphaBlend(r, a, combiner.CMax)*0.2126 +
		combiner.AlphaBlend(g, a, combiner.CMax)*0.7152 +
		combiner.AlphaBlend(b, a, combiner.CMax)*0.0722)
}
