package combiner

import (
	"image/color"
)

// Combiner represents a type that combines the channels of a color into a numerically sortable value.
type Combiner interface {
	Name() string
	Combine(c color.Color) uint64
}

var all []Combiner

// Register registers a combiner
func Register(cmd Combiner) {
	all = append(all, cmd)
}

// Registered returns all the registered combiners
func Registered() []Combiner {
	all2 := make([]Combiner, len(all))
	copy(all2, all)

	return all2
}
