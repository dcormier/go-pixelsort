package main

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dcormier/go-pixelsort/combiner"
	"github.com/dcormier/go-pixelsort/combiner/all"
	"github.com/dcormier/go-pixelsort/sortablecolor"
)

var update = flag.Bool("update", false, "update .golden files")

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

func imageFromFile(tb testing.TB, file string) image.Image {
	f, err := os.Open(file)
	require.NoError(tb, err)

	img, _, err := image.Decode(f)
	require.NoError(tb, err)

	return img
}

func sortImage(tb testing.TB, srcImg image.Image, combiner combiner.Combiner) image.Image {
	buffer, bounds := sortablecolor.SortableBufferFromImage(srcImg, combiner)

	sort.Sort(sort.Reverse(buffer))

	destImg := image.NewNRGBA64(bounds)

	buffer.ToImage(destImg)

	return destImg
}

func savePNG(tb testing.TB, target string, img image.Image) {
	f, err := os.Create(target)
	require.NoError(tb, err)

	defer f.Close()

	err = png.Encode(f, img)
	require.NoError(tb, err)
}

func compareImages(tb testing.TB, expectedFilename string, expectedImg, actualImg []byte) {
	if !assert.EqualValues(tb, expectedImg, actualImg) {
		baseName, _ := splitFilename(expectedFilename)

		if baseBaseName, ext2 := splitFilename(baseName); ext2 == ".golden" {
			baseName = baseBaseName
		}

		err := ioutil.WriteFile(baseName+".failure.png", actualImg, 0)
		require.NoError(tb, err)

		if tb.Failed() {
			tb.FailNow()
		}
	}
}

func compareImageToFile(tb testing.TB, expectedFilename string, actualImg []byte) {
	expectedBytes, err := ioutil.ReadFile(expectedFilename)
	require.NoError(tb, err)

	compareImages(tb, expectedFilename, expectedBytes, actualImg)
}

func testAllCombiners(t *testing.T, srcImgName string, srcImg image.Image) {
	imgBaseName, _ := splitFilename(srcImgName)

	for _, cmb := range all.All() {
		t.Run(cmb.Name(), func(t *testing.T) {
			cmb := cmb

			t.Parallel()

			destImg := sortImage(t, srcImg, cmb)

			goldenName := imgBaseName + "." + cmb.Name() + ".golden.png"

			if *update {
				savePNG(t, goldenName, destImg)
				return
			}

			destImgBuf := &bytes.Buffer{}
			err := png.Encode(destImgBuf, destImg)
			require.NoError(t, err)

			compareImageToFile(t, goldenName, destImgBuf.Bytes())
		})
	}
}

func TestImage(t *testing.T) {
	// t.SkipNow()
	t.Parallel()

	img := testImage(t)
	imgFilename := path.Join("testdata", "colors.png")

	if *update {
		savePNG(t, imgFilename, img)
	}

	testAllCombiners(t, imgFilename, img)
}

func splitFilename(filename string) (base, ext string) {
	ext = filepath.Ext(filename)
	base = filename[:len(filename)-len(ext)]

	return base, ext
}

func TestImageDemo(t *testing.T) {
	t.Parallel()

	for _, srcFile := range []string{
		path.Join("testdata", "2G7kAHr.jpg"),
		path.Join("testdata", "starbound_by_steelsoldat-d71fm1o.jpg"),
	} {
		t.Run(srcFile, func(t *testing.T) {
			srcFile := srcFile

			t.Parallel()

			srcImg := imageFromFile(t, srcFile)
			testAllCombiners(t, srcFile, srcImg)
		})
	}
}
