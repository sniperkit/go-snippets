package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/xiam/exif"
)

func getDefaultOutputFilename(inputImage string) string {
	inputFileName := filepath.Base(inputImage)
	doesFileNameHaveExtension := strings.Contains(inputFileName, ".")
	outputFileName := ""
	if doesFileNameHaveExtension == true {
		outputStr := strings.Split(inputFileName, ".")
		outputFileName = outputStr[0] + "_clean." + outputStr[1]
	}
	if doesFileNameHaveExtension == false {
		outputFileName = inputFileName + "_clean"
	}
	return outputFileName
}

func doCleaning(inputImage string, cleanImage bool, outputImage string, jsonOutput bool) error {
	// Check user provided an input file
	if inputImage == "" {
		err := errors.New("user did not provide an image")
		flag.Usage()
		return err
	}

	// Check user provided an actual file
	if _, err := os.Stat(inputImage); os.IsNotExist(err) {
		return err
	}

	inputFile, err := os.Open(inputImage)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	imageData, imageType, err := image.Decode(inputFile)
	if err != nil {
		return err
	}

	// Name output file if no name provided by user
	if outputImage == "" {
		outputImage = getDefaultOutputFilename(inputImage)
	}

	if cleanImage == true {
		// See if image type is not supported
		if imageType != "png" && imageType != "jpeg" && imageType != "jpg" {
			fmt.Println(imageType)
			err := errors.New("image type not yet supported")
			return err
		}

		// Create output file
		outputFile, err := os.Create(outputImage)
		if err != nil {
			return err
		}

		if imageType == "png" {
			// AFAIK PNG does not have EXIF data, but let's make a new file for giggles
			png.Encode(outputFile, imageData)
		}
		if imageType == "jpeg" || imageType == "jpg" {
			jpeg.Encode(outputFile, imageData, nil)
		}
		fmt.Println("Clean version of", inputImage, "saved in the current directory as", outputImage)
		outputFile.Close()
	}

	// Optionally output metadata as JSON
	if jsonOutput == false {
		return nil
	}
	data, err := exif.Read(inputImage)
	// We expect png to error as there is no EXIF data
	if err != nil && imageType != "png" {
		return err
	}

	if data != nil {
		metadata := data.Tags
		metadata["Filename"] = inputImage
		jsonMetadata, err := json.MarshalIndent(metadata, "", "   ")
		if err != nil {
			return err
		}
		os.Stdout.Write(jsonMetadata)
	}
	if data == nil {
		fmt.Printf("{\"Filename\": \"%v\"}\n", inputImage)
	}

	return nil
}

func main() {
	inputImagePtr := flag.String("input", "", "Image file to scrub metadata from (required)")
	cleanImagePtr := flag.Bool("clean", true, "Generate cleaned image (optional)")
	outputImagePtr := flag.String("output", "", "Output file of cleaned image (optional)")
	jsonOutputPtr := flag.Bool("json", false, "Print JSON metadata to stdout (optional)")
	flag.Parse()
	err := doCleaning(*inputImagePtr, *cleanImagePtr, *outputImagePtr, *jsonOutputPtr)
	if err != nil {
		fmt.Println(os.Stderr, "[!]", err)
		os.Exit(1)
	}
}
