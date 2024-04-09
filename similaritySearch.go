package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
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
	var sum float64
	for i := range h1.HCompressed {
		sum += (min(h1.HCompressed[i], h2.HCompressed[i]))
	}

	h2.Val = sum

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

func Images(folderPath string) ([]string, error) {
	var imageNames []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Verifie si le fichier est un .jpg
		if filepath.Ext(path) == ".jpg" {
			imageNames = append(imageNames, filepath.Base(path))
		}

		return nil
	})
	for i := range imageNames {
		imageNames[i] = folderPath + "/" + imageNames[i]
	}
	return imageNames, err
}

func slices(list []string, k int) [][]string {
	var divided [][]string
	n := len(list)
	for i := 0; i < k; i++ {
		start := i * n / k
		end := (i + 1) * n / k
		divided = append(divided, list[start:end])
	}
	return divided
}

func main() {
	start := time.Now()

	h := make(chan Histo)

	image, _ := Images(os.Args[2])
	k := 1048
	imageS := slices(image, k)
	for i := 0; i < k; i++ {
		go computeHistograms(imageS[i], 3, h)
	}

	t := make(chan Histo, 1)
	go func() {
		his, _ := computeHistogram("queryImages/"+os.Args[1], 3)
		t <- his
	}()
	h1 := <-t
	similarImages := make([]Histo, 5)
	for i := 0; i < len(image); i++ {
		h2 := <-h
		computeSimilarity(h1, h2)
		if len(similarImages) < 5 {
			similarImages = append(similarImages, h2)
			sort.Slice(similarImages, func(i, j int) bool {
				return similarImages[i].Val > similarImages[j].Val
			})
		} else if similarImages[4].Val < h2.Val {
			similarImages[4] = h2
			sort.Slice(similarImages, func(i, j int) bool {
				return similarImages[i].Val > similarImages[j].Val
			})
		}

	}

	fmt.Println("The 5 most similar images are: ")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d. %s with a similarity of : %f \n", i+1, similarImages[i].Name, similarImages[i].Val)
		fmt.Println()
	}
	end := time.Now()
	fmt.Println("Time taken: ", end.Sub(start))
	close(h)
	close(t)
}
