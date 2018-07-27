package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"sort"

	_ "golang.org/x/image/tiff"

	"github.com/dcormier/go-pixelsort/sortablecolor"
)

const (
	formatJpeg = "jpeg"
	formatPng  = "png"
	formatGif  = "gif"
	formatTiff = "tiff"
)

func writeHelp(prog string) {
	fmt.Printf("%v <input> [output_sorted.png]\n", prog)
}

func getArgs() (input, output string, err error) {
	switch len(os.Args) {
	case 2:
		input = os.Args[1]
		ext := path.Ext(input)
		output = input[:len(input)-len(ext)] + "_sorted"
		break

	case 3:
		input = os.Args[1]
		ext := path.Ext(os.Args[2])
		output = input[:len(os.Args[2])-len(ext)] + "_sorted"
		break

	default:
		writeHelp(os.Args[0])
		err = errors.New("")
		break
	}

	return
}

func main() {
	input, output, err := getArgs()
	if err != nil {
		return
	}

	reader, err := os.Open(input)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	img, imgFmt, err := image.Decode(reader)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	reader.Close()

	buffer, bounds := readImage(img)

	fmt.Println("Image metadata:")
	fmt.Printf("    File:   %s\n", input)
	fmt.Printf("    Format:     % 5s\n", imgFmt)
	fmt.Printf("    Width:      % 5d\n", bounds.Dx())
	fmt.Printf("    Height:     % 5d\n", bounds.Dy())
	fmt.Printf("    Pixels: % 9d\n", bounds.Dx()*bounds.Dy())
	fmt.Println()

	sort.Sort(sort.Reverse(sortablecolor.SortableBuffer(buffer)))

	img2 := image.NewRGBA64(bounds)

	writeImage(img2, buffer)

	outFmt := imgFmt

	switch imgFmt {
	case formatJpeg:
		outFmt = formatJpeg
		break
		// fallthrough
	case formatPng:
		// No change
		// break
		fallthrough

	case "gif":
		fallthrough
	case "tiff":
		outFmt = formatPng
		break

	default:
		fmt.Printf("Not sure what to do with image format %q. Defaulting to writing a PNG.", imgFmt)
		outFmt = formatPng
		break
	}

	fmt.Printf("Output format is %v\n", outFmt)

	switch outFmt {
	case formatJpeg:
		output += ".jpg"
		break

	case formatPng:
		output += ".png"
		break
	}

	fmt.Printf("Output will be written to: %v\n", output)
	fmt.Println()

	writer, err := os.Create(output)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	switch outFmt {
	case formatJpeg:
		var opts jpeg.Options
		opts.Quality = 100
		err = jpeg.Encode(writer, img2, &opts)
		break

	case formatPng:
		err = png.Encode(writer, img2)
		break
	}

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	writer.Close()
}

func readImage(img image.Image) (buffer []sortablecolor.SortableColor, bounds image.Rectangle) {
	bounds = img.Bounds()

	// Allocate the memory for the buffer we're going to sort
	buffer = make([]sortablecolor.SortableColor, bounds.Dx()*bounds.Dy())

	// Read the image into the buffer
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			// buffer[x*bounds.Dy() : (x+1)*bounds.Dy()][y] = img.At(x, y)
			buffer[y*bounds.Dx():][x].Set(img.At(x, y))
		}
	}

	return
}

// SettableImage represents an image.Image with the ability to set color at specific pixels
type SettableImage interface {
	image.Image
	Set(x, y int, c color.Color)
}

func writeImage(img SettableImage, buffer []sortablecolor.SortableColor) {
	bounds := img.Bounds()

	var c color.Color

	// Write it back out to the image
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			// c = buffer[x*bounds.Dy() : (x+1)*bounds.Dy()][y]
			c = buffer[y*bounds.Dx():][x].Color
			img.Set(x, y, img.ColorModel().Convert(c))
		}
	}
}

func sortImg(buffer []sortablecolor.SortableColor) {
	// Bubble sort http://www.sorting-algorithms.com/bubble-sort
	// http://www.planet-source-code.com/vb/scripts/ShowCode.asp?txtCodeId=2966&lngWId=3
	// Reversed for more favorable memory access

	var temp sortablecolor.SortableColor

	for i := len(buffer) - 2; i >= 0; i-- {
		for j := i; j < len(buffer)-1; j++ {
			if buffer[j].Compare(buffer[j+1]) < 0 {
				temp = buffer[j]
				buffer[j] = buffer[j+1]
				buffer[j+1] = temp
			}
		}
	}
}
