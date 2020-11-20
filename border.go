package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func process(path string) error {
	mw := imagick.NewMagickWand()

	if err := mw.ReadImage(path); err != nil {
		// probably just not an image
		return nil
	}

	height := mw.GetImageHeight()
	width := mw.GetImageWidth()

	if height == 650 && width == 650 {
		// assume already processed
		return nil
	}

	bg := imagick.NewPixelWand()
	bg.SetColor("white")

	var borderWidth uint
	if height > width {
		borderWidth = (height - width) / 2
		mw.SpliceImage(borderWidth, 0, int(width), 0)
		mw.SpliceImage(borderWidth, 0, 0, 0)
	} else {
		borderWidth = (width - height) / 2
		mw.SpliceImage(0, borderWidth, 0, int(height))
		mw.SpliceImage(0, borderWidth, 0, 0)
	}

	mw.ScaleImage(650, 650)

	ext := filepath.Ext(path)
	base := path[0 : len(path)-len(ext)]
	adjusted := fmt.Sprintf("%s-adjusted%s", base, ext)
	mw.WriteImage(adjusted)

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

	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("Ignoring error: %v", err)
				return nil
			}

			if info.IsDir() {
				return nil
			}

			return process(path)
		})
}
