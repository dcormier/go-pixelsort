package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dcormier/go-pixelsort/combiner"
	"github.com/dcormier/go-pixelsort/combiner/all"
	"github.com/dcormier/go-pixelsort/combiner/perceivedoption2noalpha"
	"github.com/dcormier/go-pixelsort/sortablecolor"
)

// scaleFrac is used to scale a fraction. It can take 2/8 and scale it to 25/100 or 1/4.
// The return value is the scaled numerator to match the scaled denonimator passed in.
func scaleFrac(tb testing.TB, num, den, scaleToDen int) float64 {
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

	rowColorsCount := len(rowColors)
	swatchSize := 128

	img := image.NewNRGBA(image.Rect(0, 0, swatchSize, (rowColorsCount*swatchSize)-1))

	tb.Logf("Max: %+v; Size: %+v", img.Bounds().Max, img.Bounds().Size())

	maxX := img.Bounds().Dx()
	maxY := img.Bounds().Dy()
	rows := maxY + 1
	rowsPerColor := rows / rowColorsCount

	tb.Logf("rows = %d; rowColorsCount = %d; rowsPerColor = %v",
		rows, rowColorsCount, rowsPerColor)

	var rowColor *color.NRGBA
	x := 0
	y := 0
	for x = 0; x <= maxX; x++ {
		// math.MaxUint8 is opaque. 0 is transparent. We want this image to be written out so that
		// the pixels on the left are opaque and they gradually become fully transparent on the
		// right.
		columnAlpha := math.MaxUint8 - uint8(scaleFrac(tb, x, maxX, math.MaxUint8))
		// tb.Logf("columnAlpha for column %d: %d", x, columnAlpha)

		for y = 0; y <= maxY; y++ {
			rowColor = rowColors[int(scaleFrac(tb, y, rows, rowColorsCount))]
			rowColor.A = columnAlpha
			// rowColor.A = 255
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

func imageFromFile(tb testing.TB, file string) image.Image {
	f, err := os.Open(file)
	require.NoError(tb, err)

	img, _, err := image.Decode(f)
	require.NoError(tb, err)

	return img
}

func sortImage(tb testing.TB, srcFile, destFile string, combiner combiner.Combiner) {
	srcImg := imageFromFile(tb, srcFile)

	buffer, bounds := sortablecolor.SortableBufferFromImage(srcImg, combiner)

	sort.Sort(sort.Reverse(buffer))

	destImg := image.NewNRGBA64(bounds)

	buffer.ToImage(destImg)

	savePNG(tb, destFile, destImg)
}

func TestImage(t *testing.T) {
	// t.SkipNow()
	t.Parallel()

	img := testImage(t)

	savePNG(t, path.Join("testdata", "golden.png"), img)

	// buffer, bounds := sortablecolor.BufferFromImage(img, sortablecolor.DefaultCombiner)
	buffer, bounds := sortablecolor.SortableBufferFromImage(img, perceivedoption2noalpha.Combiner)
	sort.Sort(sort.Reverse(buffer))
	img2 := image.NewNRGBA64(bounds)

	buffer.ToImage(img2)

	savePNG(t, path.Join("testdata", "output.png"), img2)
}

func TestImageDemo(t *testing.T) {
	t.Parallel()

	for _, sourceFile := range []string{
		path.Join("testdata", "2G7kAHr.jpg"),
		path.Join("testdata", "starbound_by_steelsoldat-d71fm1o.jpg"),
	} {
		t.Run(sourceFile, func(t *testing.T) {
			sourceFile := sourceFile

			t.Parallel()

			ext := filepath.Ext(sourceFile)
			imgBaseName := sourceFile[:len(sourceFile)-len(ext)]

			for _, cmb := range all.All() {
				t.Run(cmb.Name(), func(t *testing.T) {
					cmb := cmb

					t.Parallel()

					sortImage(t, sourceFile, imgBaseName+"."+cmb.Name()+".png", cmb)
				})
			}
		})
	}
}
