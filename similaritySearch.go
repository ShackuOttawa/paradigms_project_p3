package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Histo struct {
	Name        string
	H           []int
	HCompressed []float64
	Val         float64
}

// adapted from: first example at pkg.go.dev/image
func computeHistogram(imagePath string, depth int) (Histo, error) {
	// Open the JPEG file
	file, err := os.Open(imagePath)

	if err != nil {
		return Histo{}, err
	}
	defer file.Close()

	// Decode the JPEG image
	img, err := jpeg.Decode(file)
	if err != nil {
		return Histo{}, err
	}

	// Get the dimensions of the image
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	fmt.Println("width height depth: ", width, height, depth)

	h := Histo{imagePath, make([]int, 1<<(depth*3)), make([]float64, 1<<(depth*3)), 0.0}

	// Display RGB values for the first 5x5 pixels
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Convert the pixel to RGBA
			red, green, blue, _ := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 13 reduces this to the range [0, 8].
			red >>= 16 - depth
			blue >>= 16 - depth
			green >>= 16 - depth

			index := red*8*8 + green*8 + blue

			h.H[index]++
		}
	}

	for i := 0; i < len(h.H); i++ {
		h.HCompressed[i] = float64(h.H[i]) / float64(width*height)
	}

	return h, nil
}

func computeHistograms(imagePath []string, depth int, hChan chan<- Histo) {
	for i := range imagePath {
		h, _ := computeHistogram(imagePath[i], depth)
		hChan <- h
	}

}

func computeSimilarity(h1, h2 Histo) {
}

func readFiles(argmnnts string) {

	// read the directory name from command line
	//args := os.Args

	// sike. args is now a string
	args := argmnnts

	fmt.Println("args ", args, " end of args")

	//files, err := ioutil.ReadDir(args[1])
	files, err := ioutil.ReadDir(args)
	if err != nil {
		log.Fatal(err)
	}

	// Create an array to store filenames
	var filenames []string

	// get the list of jpg files
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jpg") {
			filenames = append(filenames, file.Name())
		}
	}

	// Print the array of filenames
	fmt.Println("List of jpg images:")
	for _, filename := range filenames {
		fmt.Println(filename)
	}
}

func main() {

	// change this to any of the files in queryimages
	fmt.Println(computeHistogram("C:\\Users\\User\\Downloads\\paradigms_project_p3\\queryImages\\q00.jpg", 3))
}
