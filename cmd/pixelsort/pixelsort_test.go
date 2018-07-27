package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

// scaleFrac is used to scale a fraction. It can take 2/8 and scale it to 25/100 or 1/4.
// The return value is the scaled numerator to match the scaled denonimator passed in.
func scaleFrac(tb testing.TB, num, den, scaleToDen int64) float64 {
	return float64(scaleToDen) *
		(float64(num) * float64(scaleToDen)) /
		(float64(den) * float64(scaleToDen))
}

func TestScaleFrac(t *testing.T) {
	t.Parallel()

	require.Equal(t, float64(25), scaleFrac(t, 2, 8, 100), "2/8 == 25/100")
	require.Equal(t, float64(1), scaleFrac(t, 2, 8, 4), "2/8 == 1/4")
	require.Equal(t, float64(0.5), scaleFrac(t, 2, 8, 2), "2/8 == 0.5/2")
}

// testImage creates a test image to use for tests.
func testImage(tb testing.TB) image.Image {
	rowColors := []*color.NRGBA{
		&color.NRGBA{R: 255, A: 255},
		&color.NRGBA{G: 255, A: 255},
		&color.NRGBA{B: 255, A: 255},
		&color.NRGBA{A: 255},
		&color.NRGBA{R: 64, G: 64, B: 64, A: 255},
		&color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	}

	swatchSize := 128

	img := image.NewRGBA(image.Rect(0, 0, swatchSize, (len(rowColors)*swatchSize)-1))

	tb.Logf("Max: %+v; Size: %+v", img.Bounds().Max, img.Bounds().Size())

	maxX := img.Bounds().Dx()
	maxY := img.Bounds().Dy()
	rows := maxY + 1
	rowsPerColor := rows / len(rowColors)

	tb.Logf("rows = %d; rowColors = %d; rowsPerColor = %v",
		rows, len(rowColors), rowsPerColor)

	maxAlpha := float32(math.MaxUint8)
	maxXTimesMaxAlpha := float32(maxX) * maxAlpha
	var rowColor *color.NRGBA
	x := 0
	y := 0
	for x = 0; x <= maxX; x++ {
		// math.MaxUint8 is opaque. 0 is transparent. We want this image to be written out so that
		// the pixels on the left are opaque and they gradually become fully transparent on the right.
		columnAlpha :=
			uint8(maxAlpha - (maxAlpha * (float32(x) * maxAlpha) / maxXTimesMaxAlpha))

		// tb.Logf("columnAlpha for column %d: %d", x, columnAlpha)
		for y = 0; y <= maxY; y++ {
			rowColor = rowColors[y/rowsPerColor]
			rowColor.A = columnAlpha
			img.Set(x, y, *rowColor)
		}
	}

	return img
}

func savePNG(tb testing.TB, target string, img image.Image) {
	f, err := os.Create(target)
	require.NoError(tb, err)

	defer f.Close()

	err = png.Encode(f, img)
	require.NoError(tb, err)
}

func TestImage(t *testing.T) {
	t.Parallel()

	img := testImage(t)

	savePNG(t, path.Join("testdata", "golden.png"), img)
}
