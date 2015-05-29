// Copyright 2013 Google. All rights reserved.
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

// The example tool demonstrates use of the gifresize package by resizing an
// image.
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/nfnt/resize"
	"willnorris.com/go/gifresize"
)

var width, height uint

func init() {
	flag.UintVar(&width, "width", 0, "resize width")
	flag.UintVar(&height, "height", 0, "resize height")
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage of gifresize/example:
example [-width <width>] [-height <height>] src > dest

src can be a local image file or a URL of a remote image.  The transformed image 
will be written to stdout, so output should be redirected.

Flags:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	r, err := read(flag.Arg(0))
	if err != nil {
		log.Fatalf("error reading input: %v", err)
	}

	tx := func(m image.Image) image.Image {
		return resize.Resize(width, height, m, resize.NearestNeighbor)
	}

	err = gifresize.Process(os.Stdout, r, tx)
	if err != nil {
		log.Fatalf("error processing image: %v", err)
	}
}

// read the image file specified by input.  input may be an absolute URL of a
// remote image, or a path to a file on disk.
func read(input string) (io.Reader, error) {
	if f, err := os.Open(input); err == nil {
		return f, nil
	}

	if u, err := url.Parse(input); err == nil && u.IsAbs() {
		if u.Scheme != "http" && u.Scheme != "https" {
			return nil, errors.New("only http(s) URLs are supported")
		}
		resp, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}
		return resp.Body, nil
	}

	return nil, errors.New("input must specify a file or http(s) URL")
}
