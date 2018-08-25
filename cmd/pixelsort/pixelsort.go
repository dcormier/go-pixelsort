package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"sort"

	_ "golang.org/x/image/tiff"

	"github.com/dcormier/go-pixelsort/combiner/perceivedoption2"
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

	combiner := perceivedoption2.New()

	buffer, bounds := sortablecolor.SortableBufferFromImage(img, combiner)

	fmt.Println("Image metadata:")
	fmt.Printf("    File:   %s\n", input)
	fmt.Printf("    Format:     % 5s\n", imgFmt)
	fmt.Printf("    Width:      % 5d\n", bounds.Dx())
	fmt.Printf("    Height:     % 5d\n", bounds.Dy())
	fmt.Printf("    Pixels: % 9d\n", bounds.Dx()*bounds.Dy())
	fmt.Println()

	sort.Sort(sort.Reverse(buffer))

	img2 := image.NewRGBA64(bounds)

	buffer.ToImage(img2)

	outFmt := imgFmt

	switch imgFmt {
	case formatJpeg:
		outFmt = formatJpeg

	case formatPng, "gif", "tiff":
		outFmt = formatPng

	default:
		fmt.Printf("Not sure what to do with image format %q. Defaulting to writing a PNG.", imgFmt)
		outFmt = formatPng
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
