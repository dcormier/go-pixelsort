package sortablecolor

import (
	"image"
	"image/color"
	"math"
	"sort"
)

const cMax uint32 = 65535

// SortableColor is a color.Color that can be sorted by relative brightness as compared to other instances
type SortableColor struct {
	Color color.Color
	v     uint64
}

// Set assigns the color (and relative brightness) of this instance
func (sc *SortableColor) Set(c color.Color) {
	sc.Color = c

	r, g, b, a := c.RGBA()

	// Each channel only returns 16-bit colors, as per http://blog.golang.org/go-image-package
	// They'll safely fit in a single uint64
	//sc.v = uint64(r|1) * uint64(g|1) * uint64(b|1) * uint64(a|1)

	if false {
		// http://stackoverflow.com/a/3968341/297468
		// Do alpha blending
		sc.v = uint64(alphaBlend(r, a, cMax)*0.3 +
			alphaBlend(g, a, cMax)*0.59 +
			alphaBlend(b, a, cMax)*0.11)
	} else if false {
		// http://stackoverflow.com/a/6449381/297468
		// standard, objective
		sc.v = uint64(alphaBlend(r, a, cMax)*0.2126 +
			alphaBlend(g, a, cMax)*0.7152 +
			alphaBlend(b, a, cMax)*0.0722)
	} else if false {
		// http://stackoverflow.com/a/6449381/297468
		// perceived, option 1
		sc.v = uint64(alphaBlend(r, a, cMax)*0.299 +
			alphaBlend(g, a, cMax)*0.587 +
			alphaBlend(b, a, cMax)*0.114)
	} else if true {
		// http://stackoverflow.com/a/6449381/297468
		// perceived, option 2
		sc.v = uint64(math.Sqrt(math.Pow(alphaBlend(r, a, cMax)*0.241, 2) +
			math.Pow(alphaBlend(g, a, cMax)*0.691, 2) +
			math.Pow(alphaBlend(b, a, cMax)*0.068, 2)))
	} else {
		// Ignore alpha values

		// http://stackoverflow.com/a/6449381/297468
		// perceived, option 2
		sc.v = uint64(math.Sqrt(math.Pow(float64(r)*0.241, 2) +
			math.Pow(float64(g)*0.691, 2) +
			math.Pow(float64(b)*0.068, 2)))
	}
}

// alphaBlend helps convert RGBA color values to RGB
func alphaBlend(channelValue, alphaValue, backgroundValue uint32) (blended float64) {
	// http://stackoverflow.com/questions/2049230/convert-rgba-color-to-rgb
	// http://yolijn.com/convert-rgba-to-rgb

	alpha := float64(alphaValue / cMax)

	blended = (1.0-alpha)*float64(backgroundValue) + alpha*float64(channelValue)

	return
}

// Compare compares the relative brightness of SortableColor to another SortableColor
func (sc *SortableColor) Compare(sc2 SortableColor) int {
	if sc.v < sc2.v {
		return -1
	} else if sc.v > sc2.v {
		return 1
	} else {
		return 0
	}
}

var _ sort.Interface = SortableBuffer(nil)

// SortableBuffer is a []SortableColor that implements sort.Interface
// Implements http://golang.org/pkg/sort/#Interface
type SortableBuffer []SortableColor

// BufferFromImage reads in image into a SortableBuffer
func BufferFromImage(img image.Image) (SortableBuffer, image.Rectangle) {
	bounds := img.Bounds()

	// Allocate the memory for the buffer we're going to sort
	buffer := make(SortableBuffer, bounds.Dx()*bounds.Dy())

	// Read the image into the buffer
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			// buffer[x*bounds.Dy() : (x+1)*bounds.Dy()][y] = img.At(x, y)
			buffer[y*bounds.Dx():][x].Set(img.At(x, y))
		}
	}

	return buffer, bounds
}

func (b SortableBuffer) Len() int {
	return len(b)
}

func (b SortableBuffer) Less(i, j int) bool {
	return b[i].Compare(b[j]) < 0
}

func (b SortableBuffer) Swap(i, j int) {
	temp := b[j]
	b[j] = b[i]
	b[i] = temp
}
