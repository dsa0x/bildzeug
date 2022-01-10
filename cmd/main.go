package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dsa0x/bildzeug/blurify"
)

func main() {

	imgFile, err := os.Open("img1.png")
	if err != nil {
		fmt.Printf("unable to open file: %v", err)
	}

	imgStats, err := imgFile.Stat()
	if err != nil {
		fmt.Printf("unable to get file stats: %v", err)
	}

	imgName := strings.Split(imgStats.Name(), ".")[0]
	file, err := os.Create(fmt.Sprintf("%v_blur.jpg", imgName))
	if err != nil {
		fmt.Printf("unable to create blur file")
	}

	opts := blurify.BlurOptions{KernelSize: 3, Sigma: 3, Filter: blurify.Gaussian}
	err = blurify.Blur(imgFile, file, opts)
	if err != nil {
		fmt.Println(err)
	}
}
