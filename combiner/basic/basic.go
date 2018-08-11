// Package basic implements Combiner using the most basic of methods. Not a good visual sort.
package basic

import (
	"image/color"

	"github.com/dcormier/go-pixelsort/combiner"
)

// Combiner is an instance of this Combiner
var Combiner = New()

var _ combiner.Combiner = (*basic)(nil)

type basic struct{}

// New creates a new, basic, combiner.Combiner
func New() combiner.Combiner {
	return &basic{}
}

func (*basic) Name() string {
	return "basic"
}

func (*basic) Combine(c color.Color) uint64 {
	r, g, b, a := c.RGBA()

	// Each channel only returns 16-bit colors, as per http://blog.golang.org/go-image-package
	// They'll safely fit in a single uint64
	return uint64(r|1) * uint64(g|1) * uint64(b|1) * uint64(a|1)
}
