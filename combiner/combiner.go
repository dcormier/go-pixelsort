package combiner

import (
	"image/color"
)

// Combiner represents a type that combines the channels of a color into a numerically sortable value.
type Combiner interface {
	Name() string
	Combine(c color.Color) uint64
}
