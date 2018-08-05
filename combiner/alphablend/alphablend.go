// Package alphablend implements combiner.Combiner using
// http://stackoverflow.com/a/3968341/297468
// Assumes a white background to blend with.
package alphablend

import (
	"image/color"

	"github.com/dcormier/go-pixelsort/combiner"
)

func init() {
	combiner.Register(New())
}

var _ combiner.Combiner = (*alphaBlend)(nil)

type alphaBlend struct{}

// New creates a new combiner.Combiner that uses alpha blending (assumes a white background)
func New() combiner.Combiner {
	return &alphaBlend{}
}

func (*alphaBlend) Name() string {
	return "alpha blend"
}

func (*alphaBlend) Combine(c color.Color) uint64 {
	r, g, b, a := c.RGBA()

	return uint64(combiner.AlphaBlend(r, a, combiner.CMax)*0.3 +
		combiner.AlphaBlend(g, a, combiner.CMax)*0.59 +
		combiner.AlphaBlend(b, a, combiner.CMax)*0.11)
}
