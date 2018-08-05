// Package all simple imports all the combiners so that they're all registered
package all

import (
	// These are blank imports just to get their init functions to run to register them.
	_ "github.com/dcormier/go-pixelsort/combiner/alphablend"
	_ "github.com/dcormier/go-pixelsort/combiner/basic"
	_ "github.com/dcormier/go-pixelsort/combiner/perceivedoption1"
	_ "github.com/dcormier/go-pixelsort/combiner/perceivedoption2"
	_ "github.com/dcormier/go-pixelsort/combiner/perceivedoption2noalpha"
	_ "github.com/dcormier/go-pixelsort/combiner/standardobjective"
)
