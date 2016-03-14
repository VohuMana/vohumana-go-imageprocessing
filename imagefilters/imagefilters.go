package imagefilters

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

var ImageFilterMap = map[string]ImageFilter {
	"ExtractRedChannel": ExtractRedChannelFilter{},
	"ExtractGreenChannel": ExtractGreenChannelFilter{},
	"ExtractBlueChannel": ExtractBlueChannelFilter{},
	"HistogramNormalization": HistogramEqualizationFilter{},
	"FindEdgesWithSobel": SobelImageFilter{},
    "ConvertToGrayscale": ConvertToGrayscale{},
}

// COLOR FUNCTIONS

type HSL struct {
	H, S, L uint32
}

func (h HSL) RGBA() (r, g, b, a uint32) {
	return h.H, h.S, h.L, 0
}

func init() {
	// Package init stuff here
}

func threewayMax(a, b, c float64) float64 {
	return math.Max(math.Max(a, b), math.Max(b, c))
}

func threewayMin(a, b, c float64) float64 {
	return math.Min(math.Min(a, b), math.Min(b, c))
}

func rgbToHSL(r, g, b uint32) (float64, float64, float64) {	
	// Get rgb values from 0.0 to 1.0f
	redNormalized := float64(r) / float64(math.MaxUint16)
	greenNormailzed := float64(g) / float64(math.MaxUint16)
	blueNormalized := float64(b) / float64(math.MaxUint16)

	// Threeway max and min
	max := threewayMax(redNormalized, greenNormailzed, blueNormalized)
	min := threewayMin(redNormalized, greenNormailzed, blueNormalized)

	h := (max + min) / 2.0
	s := h
	l := h

	if (max == min)	{
		// achromatic
		h = 0
		s = 0
	} else {
		d := max - min
		
		if (l > 0.5) {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}

		if (max == redNormalized) {
			h = (greenNormailzed - blueNormalized) / d

			if (greenNormailzed < blueNormalized) {
				h += 6.0
			}
		} else if (max == greenNormailzed) {
			h = ((blueNormalized - redNormalized) / d) + 2.0
		} else {
			h = ((redNormalized - greenNormailzed) / d) + 4.0
		}

		h /= 6.0
	}


	// Convert to uint32 to conform to the go RGBA uint32 color model
	// H := uint32(h * math.MaxUint32)
	// S := uint32(s * math.MaxUint32)
	// L := uint32(l * math.MaxUint32)

	// fmt.Printf("RGBTOHSL - R: %v G: %v B: %v H: %v S: %v L: %v\n", uint8(redNormalized * math.MaxUint8), uint8(greenNormailzed * math.MaxUint8), uint8(blueNormalized * math.MaxUint8), h, s, l)
	
	return h,s,l//HSL{H, S, L}
}

func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	// H, S, L, _ := c.RGBA()

	// h := float64(H) / float64(math.MaxUint32)
	// s := float64(S) / float64(math.MaxUint32)
	// l := float64(L) / float64(math.MaxUint32)

	var r, g, b float64
	if (s == 0.0) {
		// achromatic
		r = l
		g = l
		b = l
	} else {
		var q float64

		if 	(l < 0.5) {
			q = l * (1.0 + s)
		} else {
			q = l + s - l * s
		}

		p := float64(2.0 * l - q)
		r = hueToRGB(p, q, h + (1.0 / 3.0))
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h - (1.0 / 3.0))
	}

	R := uint8(r * math.MaxUint8)
	G := uint8(g * math.MaxUint8)
	B := uint8(b * math.MaxUint8)

	// fmt.Printf("HSLTORGB - R: %v G: %v B: %v H: %v S: %v L: %v\n\n", R, G, B, h, s, l)

	return R, G, B //color.RGBA{R, G, B, math.MaxUint8}
}

func hueToRGB(p, q, t float64) float64 {
	if (t < 0.0) {
		t += 1.0
	}

	if (t > 1.0) {
		t -= 1.0
	}
	// fmt.Printf("HUETORGB - P: %v Q: %v T: %v\n", p, q, t)
	if (t < (1.0 / 6.0)) {
		return p + (q - p) * 6.0 * t
	} else if (t < (1.0 / 2.0)) {
		return q
	} else if (t < (2.0 / 3.0)) {
		return p + (q - p) * (2.0 / 3.0 - t) * 6.0
	}

	return p
}


// FILTERS

type ImageFilter interface {
	Apply(img image.Image) image.Image
}

type HistogramEqualizationFilter struct {
}

type ExtractRedChannelFilter struct {
}

type ExtractGreenChannelFilter struct {
}

type ExtractBlueChannelFilter struct {
}

type SobelImageFilter struct {
}

type ConvertToGrayscale struct {
}

func (f HistogramEqualizationFilter) Apply(img image.Image) image.Image {
	normalizedImage := image.NewRGBA(img.Bounds())

	var histogram [256]uint

	// Convert image to HSL color space and generate the histogram
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			_, _, L := rgbToHSL(r, g, b)
			newL := uint8(L * float64(math.MaxUint8))
			histogram[newL]++
			// fmt.Printf("H: %v S: %v L: %v\n", H, S, L)
		}
	}
	fmt.Println("Histogram calculated")

	var normailzedChannel [256]float64
	totalPixels := b.Max.X * b.Max.Y

	for i := 0; i < 256; i++ {
		normailzedChannel[i] = float64(histogram[i]) / float64(totalPixels)
	}

	var newHistogram [256]float64
	for i := 0; i < 256; i++ {
		for j := 0; j <= i; j++ {
			newHistogram[i] += normailzedChannel[j]
		}
	}

	var newValues [256]uint8
	for i := 0; i < 256; i++ {
		newValues[i] = uint8((newHistogram[i] * math.MaxUint8) + 0.5)
		// fmt.Printf("Value: %v OldHist: %v NewHist: %v normailzedChannel: %v\n", i, histogram[i], newHistogram[i], normailzedChannel[i])
	}

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			H, S, L := rgbToHSL(r, g, b)
			newL := uint8(L * float64(math.MaxUint8))
			L = (float64(newValues[newL]) / float64(math.MaxUint8))
			R, G, B := hslToRGB(H, S, L)
			normalizedImage.Set(x, y, color.RGBA{R, G, B, math.MaxUint8}) 
		}
	}

	return normalizedImage
}

func (f ExtractRedChannelFilter) Apply(img image.Image) image.Image {
	redChannel := image.NewRGBA(img.Bounds())

	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, _, _, _ := img.At(x, y).RGBA()

			redChannel.Set(x, y, color.RGBA{uint8(r), 0, 0, math.MaxUint8})
		}
	}

	return redChannel
}

func (f ExtractGreenChannelFilter) Apply(img image.Image) image.Image {
	greenChannel := image.NewRGBA(img.Bounds())

	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, g, _, _ := img.At(x, y).RGBA()

			greenChannel.Set(x, y, color.RGBA{0, uint8(g), 0, math.MaxUint8})
		}
	}

	return greenChannel
}

func (f ExtractBlueChannelFilter) Apply(img image.Image) image.Image {
	blueChannel := image.NewRGBA(img.Bounds())

	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, blue, _ := img.At(x, y).RGBA()

			blueChannel.Set(x, y, color.RGBA{0, 0, uint8(blue), math.MaxUint8})
		}
	}

	return blueChannel
}

func (f SobelImageFilter) Apply(img image.Image) image.Image {
	return img
}

func (f ConvertToGrayscale) Apply(img image.Image) image.Image {
    grayscaleImage := image.NewRGBA(img.Bounds())
    
    bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			red, green, blue, _ := img.At(x, y).RGBA()
            _, _, L := rgbToHSL(red, green, blue)
			r, g, b := hslToRGB(0.0, 0.0, L)
            
            grayscaleImage.Set(x, y, color.RGBA{r, g, b, math.MaxUint8})
		}
	}
    
    return grayscaleImage
}