// Copyright 2013 Daniel Pupius. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

// This demonstrates some problems when resizing animated gifs.

// Frames in an animated gif aren't necessarily the same size, subsequent
// frames are overlayed on previous frames. Therefore, resizing the frames
// individually may cause problems due to aliasing of transparent pixels. This
// example tries to avoid this by building frames from all previous frames and
// resizing the frames as RGB.
//
// There is still a problem when resizing shapes.gif results in frames with
// color artifacts. e.g. See frame 19 in shapes.out.gif vs. frames/shapes.19.jpg

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func main() {
	process("shapes")
	process("blob")
}

func process(filename string) {

	// Open image file.
	f, err := os.Open(filename + ".gif")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	// Decode the original gif.
	im, err := gif.DecodeAll(f)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create a new RGBA image to hold the incremental frames.
	firstFrame := im.Image[0].Bounds()
	b := image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy())
	img := image.NewRGBA(b)

	// Resize each frame.
	for index, frame := range im.Image {
		bounds := frame.Bounds()
		draw.Draw(img, bounds, frame, bounds.Min, draw.Src)
		im.Image[index] = ImageToPaletted(ProcessImage(img))
	}

	// Write resized gif.
	out, err := os.Create(filename + ".out.gif")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out.Close()
	gif.EncodeAll(out, im)

	// Write each frame to a jpeg.
	for i := 0; i < len(im.Image); i++ {
		jout, err := os.Create(fmt.Sprintf("frames/%s.%d.jpg", filename, i))
		if err != nil {
			log.Fatal(err.Error())
		}
		defer jout.Close()
		jpeg.Encode(jout, im.Image[i], &jpeg.Options{90})
	}
}

func ProcessImage(img image.Image) image.Image {
	return resize.Resize(250, 0, img, resize.NearestNeighbor)
}

func ImageToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
}
