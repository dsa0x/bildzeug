package main

import (
	"fmt"
	"os"

	"github.com/dsa0x/bildzeug/blurify"
)

func main() {

	imgFile, err := os.Open("vini.jpeg")
	if err != nil {
		fmt.Printf("unable to open file: %v", err)
	}

	file, err := os.Create("vini_blur.jpeg")
	if err != nil {
		fmt.Printf("unable to create blur file")
	}

	opts := blurify.BlurOptions{KernelSize: 3, Sigma: 3, Filter: blurify.Gaussian}
	err = blurify.Blur(imgFile, file, opts)
	if err != nil {
		fmt.Println(err)
	}
}
