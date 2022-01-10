package blurify

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"math"
)

type FilterType string

var (
	Gaussian  FilterType = "GAUSSIAN"
	MovingAvg FilterType = "MOVING_AVG"
)

// BlurOptions are encoding parameters
// Quality ranges from 1 to 100 inclusive, higher is better.
type BlurOptions struct {
	KernelSize int
	Filter     FilterType

	// sigma is valid for gaussian filter only
	Sigma float64
}

func Blur(r io.Reader, w io.Writer, opts BlurOptions) error {
	img, err := jpeg.Decode(r)
	if opts.KernelSize == 0 {
		opts.KernelSize = 3
	}

	if opts.Sigma == 0 {
		opts.Sigma = 1
	}

	if opts.KernelSize%2 != 1 {
		return fmt.Errorf("kernel size must be odd number")
	}

	if err != nil {
		return fmt.Errorf("unable to decode image: %v", err)
	}
	bounds := img.Bounds()
	padding := opts.KernelSize / 2
	rect := image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Dx()+padding, bounds.Dy()+padding)
	newImg := image.NewRGBA(rect)
	finalImg := image.NewRGBA(bounds)

	// pad the matrix
	for i := 0; i < bounds.Max.X; i++ {
		for j := 0; j < bounds.Max.Y; j++ {
			newImg.Set(i+padding, j+padding, img.At(i, j))
		}
	}

	// TODO: use dynamic programming instead of padding
	for i := padding; i < rect.Dx(); i++ {
		for j := padding; j < rect.Dy(); j++ {

			tmp := make([][]color.Color, opts.KernelSize)
			for i := range tmp {
				tmp[i] = make([]color.Color, opts.KernelSize)
			}

			idx_x := 0
			for m := i - padding; m <= i+padding; m++ {
				idx_y := 0
				for n := j - padding; n <= j+padding; n++ {
					tmp[idx_x][idx_y] = newImg.At(m, n)
					idx_y++
				}
				idx_x++
			}

			finalImg.Set(i-padding, j-padding, convolve(tmp, opts.KernelSize, opts.Filter))
		}
	}
	err = jpeg.Encode(w, finalImg, &jpeg.Options{})
	return err
}

func convolve(pixel [][]color.Color, size int, filtertype FilterType) color.RGBA64 {
	//
	var filter [][]float64
	var divisor int
	switch filtertype {
	case Gaussian:
		filter, divisor = gaussianFilter(1.0, size), 1
	case MovingAvg:
		filter, divisor = filterMAKernel(size)
	default:
		log.Fatalf("Invalid filter type: %s", filtertype)
	}

	var sumR, sumG, sumB uint32
	for i := 0; i < len(pixel); i++ {
		for j := 0; j < len(pixel[i]); j++ {
			px := pixel[i][j]

			r, g, b, _ := px.RGBA()
			sumR += uint32(float64(r) * filter[i][j])
			sumG += uint32(float64(g) * filter[i][j])
			sumB += uint32(float64(b) * filter[i][j])

		}
	}

	r := uint16(sumR / uint32(divisor))
	g := uint16(sumG / uint32(divisor))
	b := uint16(sumB / uint32(divisor))

	return color.RGBA64{r, g, b, uint16(math.MaxUint16)}

}

// filterMAKernel uses the moving average as the kernel filter
func filterMAKernel(size int) ([][]float64, int) {
	kernel := make([][]float64, size)
	for i := range kernel {
		kernel[i] = make([]float64, size)
		for j := range kernel[i] {
			kernel[i][j] = 1
		}
	}
	return kernel, size * size
}

// binomialFilter uses the binomial filter
func binomialFilter(size int) [][]float64 {
	kernel := make([][]float64, size)
	for i := range kernel {
		kernel[i] = make([]float64, size)
		for j := range kernel[i] {
			kernel[i][j] = 1
		}
	}
	return kernel
}

// gaussianFilter implements the gaussian filter algorithm
func gaussianFilter(sigma float64, size int) [][]float64 {
	kernel := make([][]float64, size)
	for i := range kernel {
		kernel[i] = make([]float64, size)
		for j := range kernel[i] {
			rhosq := 2 * math.Pow(sigma, 2)
			ij_sq := math.Pow(float64(i), 2) + math.Pow(float64(j), 2)
			lhs := 1 / (math.Pi * rhosq)
			g := lhs * math.Exp(-(ij_sq / rhosq))
			kernel[i][j] = g
		}
	}
	return kernel
}
