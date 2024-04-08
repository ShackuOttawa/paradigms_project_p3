package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Histo struct {
	Name string
	H    []int
}

// adapted from: first example at pkg.go.dev/image
func computeHistogram(imagePath string, depth int) (Histo, error) {
	// Open the JPEG file
	file, err := os.Open(imagePath)

	if err != nil {
		return Histo{"opening issue", nil}, err
	}
	defer file.Close()

	// Decode the JPEG image
	img, _, err := image.Decode(file)
	if err != nil {
		return Histo{"decoding issue", nil}, err
	}

	// Get the dimensions of the image
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	fmt.Println("width height depth: ", width, height, depth)

	h := Histo{imagePath, make([]int, depth)}

	// Display RGB values for the first 5x5 pixels
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			// Convert the pixel to RGBA
			red, green, blue, _ := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 13 reduces this to the range [0, 8].
			red >>= 13
			blue >>= 13
			green >>= 13

			// Display the RGB values
			fmt.Printf("Pixel at (%d, %d): R=%d, G=%d, B=%d\n", x, y, red, green, blue)

			index := red*8*8 + green*8 + blue

			fmt.Println(index)

			h.H[index]++
		}
	}

	fmt.Println()

	return h, nil
}

/*
func readImage() {
	// read the image name from command line
	args := os.Args

	// Call the function to display RGB values of some pixels
	_, err := displayRGBValues(args[1], 10)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
}
*/

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

	// displays all files in the directory
	readFiles("./queryimages")

	// change this to any of the files in queryimages
	fmt.Println(computeHistogram("./queryimages/q00.jpg", 512))
}
