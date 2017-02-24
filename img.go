package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"net/http"
	"os"
	"path"
)

// img represents an individual gif to be processed.
type img struct {
	url string
	gif *gif.GIF
	err error
}

// Instantiate a new img
func newImg(url string) *img {
	return &img{
		url: url,
	}
}

// Download, reflect, and save an image.
func (m *img) process(dir string, trans bool, results chan<- *img) {
	in, err := m.download()
	if err != nil {
		m.err = err
		results <- m
		return
	}

	m.reflect(in, trans)
	err = m.save(dir)
	if err != nil {
		m.err = err
	}
	results <- m
	return
}

// Download image given a url.
func (m *img) download() (*gif.GIF, error) {
	res, err := http.Get(m.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return gif.DecodeAll(res.Body)
}

// Append a reversed version of the original gif.
func (m *img) reflect(in *gif.GIF, trans bool) {
	inLen := len(in.Image)
	outLen := inLen*2 - 2

	out := &gif.GIF{
		Image:           make([]*image.Paletted, outLen),
		Delay:           make([]int, outLen),
		LoopCount:       in.LoopCount,
		Disposal:        make([]byte, outLen),
		Config:          in.Config,
		BackgroundIndex: in.BackgroundIndex,
	}

	reversed := make([]*image.Paletted, inLen)
	for i := 0; i < inLen; i++ {
		out.Image[i] = in.Image[i]
		out.Delay[i] = in.Delay[i]
		out.Disposal[i] = in.Disposal[i]

		if !trans || i == 0 || in.Image[i].Opaque() {
			reversed[i] = clonePaletted(in.Image[i])
		} else {
			reversed[i] = clonePaletted(reversed[i-1])
			draw.Over.Draw(reversed[i], reversed[i].Bounds(), in.Image[i], image.Point{})
		}
	}

	if trans {
		optimizeTransparencyReversed(reversed)
	}

	for i := inLen; i < outLen; i++ {
		j := outLen - i
		out.Image[i] = reversed[j]
		out.Delay[i] = in.Delay[j]
		out.Disposal[i] = gif.DisposalNone
	}
	m.gif = out
}

// Save image to the given directory.
func (m *img) save(dir string) error {
	filename := path.Base(m.url)
	path := fmt.Sprintf("%s/%s", dir, filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gif.EncodeAll(file, m.gif)
}

// Log the result of image processing.
func (m *img) logResult() {
	fmt.Println(m.url)
	if m.err != nil {
		logError(m.err)
	} else {
		fmt.Println("Processed")
	}
}
