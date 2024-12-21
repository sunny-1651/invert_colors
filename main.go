package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func invertColor(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(255 - r/257),
		G: uint8(255 - g/257),
		B: uint8(255 - b/257),
		A: uint8(a / 257),
	}
}

func invertImageColors(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	inverted := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			inverted.Set(x, y, invertColor(originalColor))
		}
	}

	return inverted
}

func calculate_op_path(oldPath string) string {
	ext := filepath.Ext(oldPath)
	base := strings.TrimSuffix(oldPath, ext)

	return fmt.Sprintf("%s_inverted%s", base, ext)
}

func process(dirName string, inputFile string) {
	outputFile := calculate_op_path(filepath.Join(dirName, "inverted", inputFile))

	dirNameNew := filepath.Join(dirName, "inverted")

	if _, err := os.Stat(dirNameNew); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(dirNameNew, 0755) // 0755 is the permission for the new directory
		if err != nil {
			fmt.Printf("Error creating directory: %s\n", err)
			return
		}
		fmt.Printf("Directory '%s' created successfully.\n", dirNameNew)
	} else {
		fmt.Printf("Directory '%s' already exists.\n", dirNameNew)
	}

	file, err := os.Open(filepath.Join(dirName, inputFile))
	if err != nil {
		log.Fatalf("Failed to open input file: %v\n", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v\n", err)
	}

	inverted := invertImageColors(img)

	output, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v\n", err)
	}
	defer output.Close()

	ext := filepath.Ext(outputFile)
	switch ext {
	case ".jpg", ".jpeg":
		if err := jpeg.Encode(output, inverted, nil); err != nil {
			log.Fatalf("Failed to save JPEG image: %v\n", err)
		}
	case ".png":
		if err := png.Encode(output, inverted); err != nil {
			log.Fatalf("Failed to save PNG image: %v\n", err)
		}
	default:
		log.Fatalf("Unsupported output format: %s\n", ext)
	}

	log.Println("Image inversion completed successfully.")
}

func isImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <Directory_Name>\n", os.Args[0])
	}

	inputDir := os.Args[1]
	// process("/mnt/c/Users/saror/OneDrive/Pictures/web/Screenshots/Screenshot 2024-12-21 171918.png")

	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isImageFile(path) {
			process(inputDir, info.Name())
		} else {
			fmt.Println("Skipping ", info.Name(), " file in ", inputDir, " dir")
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path: %s", err)
	}
}
