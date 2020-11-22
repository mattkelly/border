package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func process(path string) error {
	mw := imagick.NewMagickWand()

	if err := mw.ReadImage(path); err != nil {
		// probably just not an image
		return nil
	}

	if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_DEACTIVATE); err != nil {
		panic(err)
	}

	bg := imagick.NewPixelWand()
	bg.SetColor("white")
	mw.SetBackgroundColor(bg)

	// When reading an image with layers, each layer is read as a separate
	// image into the wand. Reset to the first image before processing.
	mw.ResetIterator()
	mw = mw.MergeImageLayers(imagick.IMAGE_LAYER_MERGE)

	h := int(mw.GetImageHeight())
	w := int(mw.GetImageWidth())

	// use 20% of the longest side as the minimum border width
	var scaleTo int
	if h > w {
		scaleTo = h + int(0.2*float32(h))
	} else {
		scaleTo = w + int(0.2*float32(w))
	}

	mw.ExtentImage(uint(scaleTo), uint(scaleTo), -(scaleTo-w)/2, -(scaleTo-h)/2)

	if scaleTo > 2048 {
		// make smaller just to save space
		mw.ScaleImage(2048, 2048)
	}

	ext := filepath.Ext(path)
	base := path[0 : len(path)-len(ext)]

	// tiffs are weird; just convert to png instead
	if mw.GetImageFormat() == "TIFF" {
		mw.SetImageFormat("PNG")
		ext = ".png"
	}

	adjusted := fmt.Sprintf("%s-adjusted%s", base, ext)
	mw.WriteImage(adjusted)

	// delete original
	if err := os.Remove(path); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting %q: %v\n", path, err)
	}

	return nil
}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	dir := os.Args[1]
	newDir := dir + "-adjusted"
	if err := copy.Copy(dir, newDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error copying %q: %v\n", dir, err)
		os.Exit(2)
	}

	filepath.Walk(newDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Ignoring error: %v", err)
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if strings.Contains(filepath.Base(path), "adjusted") {
				// assume already processed
				return nil
			}

			return process(path)
		})
}
