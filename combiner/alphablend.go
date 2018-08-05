package combiner

import (
	"math"
)

// TODO:
//
// There are problems here related to alpha-premultiplication.
// The channel value comes from color.Color.RGBA(), which gives back premultiplied values (per
// https://blog.golang.org/go-image-package).
//
// We need either non-premultiplied values, or to figure out how to take advantage of that
// premultiplication (what's the math for that? possible answer in Porter and Duff's forumula, here:
// https://blog.golang.org/go-image-package).
//
// If we need to convert to non-premultiplied, just do `channelValue/alphaValue` (per bottom of
// https://microsoft.github.io/Win2D/html/PremultipliedAlpha.htm).

// CMax is the maximum value of a color (as used for alpha blending)
const CMax uint32 = math.MaxUint8

// AlphaBlend helps convert RGBA color values to RGB
func AlphaBlend(channelValue, alphaValue, backgroundValue uint32) (blended float64) {
	// http://stackoverflow.com/questions/2049230/convert-rgba-color-to-rgb
	// http://yolijn.com/convert-rgba-to-rgb

	alpha := float64(alphaValue) / float64(CMax)

	blended = ((1.0 - alpha) * float64(backgroundValue)) + (alpha * float64(channelValue))

	return
}
