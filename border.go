package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func process(path string) error {
	mw := imagick.NewMagickWand()

	if err := mw.ReadImage(path); err != nil {
		// probably just not an image
		return nil
	}

	height := int(mw.GetImageHeight())
	width := int(mw.GetImageWidth())

	bg := imagick.NewPixelWand()
	bg.SetColor("white")
	mw.SetImageBackgroundColor(bg)

	const scaleTo = 2048
	mw.ExtentImage(scaleTo, scaleTo, -(scaleTo-width)/2, -(scaleTo-height)/2)

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

			if strings.Contains(path, "adjusted") {
				// assume already processed
				return nil
			}

			return process(path)
		})
}
