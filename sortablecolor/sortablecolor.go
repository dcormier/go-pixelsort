package sortablecolor

import (
	"image"
	"image/color"
	"sort"

	"github.com/dcormier/go-pixelsort/combiner"
)

// SortableColor is a color.Color that can be sorted by relative brightness as compared to other instances
type SortableColor struct {
	Color color.Color

	// v is an absolute value representing the apparent brightness of this color
	v uint64
}

// Set assigns the color (and relative brightness) of this instance
func (sc *SortableColor) Set(c color.Color, combiner combiner.Combiner) {
	sc.Color = c

	sc.v = combiner.Combine(c)
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

// SortableBufferFromImage reads in image into a SortableBuffer
func SortableBufferFromImage(img image.Image, combiner combiner.Combiner) (SortableBuffer, image.Rectangle) {
	bounds := img.Bounds()

	// Allocate the memory for the buffer we're going to sort
	buffer := make(SortableBuffer, bounds.Dx()*bounds.Dy())

	// Read the image into the buffer
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			// buffer[x*bounds.Dy() : (x+1)*bounds.Dy()][y] = img.At(x, y)
			buffer[y*bounds.Dx():][x].Set(img.At(x, y), combiner)
		}
	}

	return buffer, bounds
}

// SettableImage represents an image.Image with the ability to set color at specific pixels
type SettableImage interface {
	image.Image
	Set(x, y int, c color.Color)
}

// ToImage writes the contents of the buffer out to the provided image using its bounds.
func (buf SortableBuffer) ToImage(img SettableImage) {
	bounds := img.Bounds()

	var c color.Color

	// Write it back out to the image
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			// c = buf[x*bounds.Dy() : (x+1)*bounds.Dy()][y]
			c = buf[y*bounds.Dx():][x].Color
			img.Set(x, y, img.ColorModel().Convert(c))
		}
	}
}

func (buf SortableBuffer) Len() int {
	return len(buf)
}

func (buf SortableBuffer) Less(i, j int) bool {
	return buf[i].Compare(buf[j]) < 0
}

func (buf SortableBuffer) Swap(i, j int) {
	temp := buf[j]
	buf[j] = buf[i]
	buf[i] = temp
}
