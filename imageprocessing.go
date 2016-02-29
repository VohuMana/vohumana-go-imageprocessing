package main 

import
(
	"flag"
	filters "github.com/vohumana/imageprocessing/imagefilters"
	"image"
	"image/png"
	"image/jpeg"
	"log"
	"os"
)

func checkError(err error) {
	if (err != nil) {
		log.Fatal(err)
	}
}

func main() {
	var imageName string
	var outputFormat string
	var outputName string
	
	// Get command line parameter for html file
	flag.StringVar(&imageName, "img", "", "filename of the image to be loaded and modified")
	flag.StringVar(&outputFormat, "fmt", "", "output format. Possible options are png or jpg")
	flag.StringVar(&outputName, "out", "", "output filename without trailing .png or .jpg")
	flag.Parse()

	if (imageName == "" || outputName == "" || (outputFormat != "png" && outputFormat != "jpg")) {
		flag.PrintDefaults()
		return
	}

	file, err := os.Open(imageName)
	checkError(err)

	var histogramFilter filters.HistogramEqualizationFilter
	img, _, err := image.Decode(file)
	checkError(err)

	newImage := histogramFilter.Apply(img)

	switch (outputFormat) {
	case "png":
		outFile, err := os.Create(outputName + ".png")
		checkError(err)
		defer outFile.Close()

		err = png.Encode(outFile, newImage)
		checkError(err)
		break
	
	case "jpg":
		outFile, err := os.Create(outputName + ".jpeg")
		checkError(err)
		defer outFile.Close()
		err = jpeg.Encode(outFile, newImage, &jpeg.Options{90})
		checkError(err)
		break
	}
}