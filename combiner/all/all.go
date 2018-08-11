// Package all simple imports all the combiners so that they're all registered
package all

import (
	// These are blank imports just to get their init functions to run to register them.
	"github.com/dcormier/go-pixelsort/combiner"
	"github.com/dcormier/go-pixelsort/combiner/alphablend"
	"github.com/dcormier/go-pixelsort/combiner/basic"
	"github.com/dcormier/go-pixelsort/combiner/perceivedoption1"
	"github.com/dcormier/go-pixelsort/combiner/perceivedoption2"
	"github.com/dcormier/go-pixelsort/combiner/perceivedoption2noalpha"
	"github.com/dcormier/go-pixelsort/combiner/standardobjective"
)

// All retuns all the known combiner.Combiners
func All() []combiner.Combiner {
	return []combiner.Combiner{
		alphablend.Combiner,
		basic.Combiner,
		perceivedoption1.Combiner,
		perceivedoption2.Combiner,
		perceivedoption2noalpha.Combiner,
		standardobjective.Combiner,
	}
}
