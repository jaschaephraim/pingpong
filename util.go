package main

import (
	"fmt"
	"image"
	"image/color"
)

func logError(err error) {
	fmt.Printf("ERROR: %s\n", err)
}

// Add transparency to a sequence of images to be played in reverse order.
func optimizeTransparencyReversed(imgs []*image.Paletted) {
	alpha := uint8(imgs[0].Palette.Index(color.Transparent))
	for i := 0; i < len(imgs)-1; i++ {
		curr := imgs[i].Pix
		prev := imgs[i+1].Pix
		for j := 0; j < len(curr) && j < len(prev); j++ {
			if curr[j] == prev[j] {
				curr[j] = alpha
			}
		}
	}
}

// Clone a paletted image.
func clonePaletted(src *image.Paletted) *image.Paletted {
	dst := image.NewPaletted(src.Bounds(), src.Palette)
	copy(dst.Pix, src.Pix)
	return dst
}
