package main

import (
	"context"
	"io/ioutil"
	"log"
	"path/filepath"

	scryfall "github.com/BlueMonday/go-scryfall"
	flag "github.com/ogier/pflag"
	"github.com/otiai10/gosseract"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

type configuration struct {
	InferenceGraphPath      string
	ImageInputPath          string
	SetSymbolDirectory      string
	SetSymbolBackgroundPath string
	GetNewSetSymbols        bool
}

func main() {
	log.Println(`Starting mox...
		Tesseract is version`, gosseract.Version(), `(this should be 4.x).
		TensorFlow is version`, tf.Version(), ` (this should be 1.x).
		GoCV is version`, gocv.Version(), ` (who knows with this one).`)

	var config configuration

	flag.StringVarP(&config.InferenceGraphPath, "inference-graph", "g", "./inference_graph/", "Path to the directory containing the inference graph.")
	flag.StringVarP(&config.ImageInputPath, "image-input", "i", "./image.jpg", "Path to an image to infer partial card information from.")
	flag.StringVarP(&config.SetSymbolDirectory, "setsymbol-dir", "s", "./set_symbols/", "A directory containing set symbols as PNGs.")
	flag.StringVarP(&config.SetSymbolBackgroundPath, "set-symbol-background", "b", "./set_symbols/background.jpg",
		"Path to the background to composite card images over when finding set symbols. THE BACKGROUND MUST BE LARGER THAN IMAGES LIKE THIS: https://img.scryfall.com/cards/png/en/ddt/60.png")
	flag.BoolVar(&config.GetNewSetSymbols, "get-new-set-symbols", false, "If set to true, mox will get new set symbols from the internet and load them.")
	flag.Parse()

	// Construct the graph
	graph, err := newDetectionGraph(filepath.Join(config.InferenceGraphPath, "frozen_inference_graph.pb"))
	if err != nil {
		log.Panic(err)
	}

	// Create a session to run inference with
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()

	// Set up Scryfall
	context := context.Background()
	scryfallClient, err := scryfall.NewClient()
	if err != nil {
		log.Panic(err)
	}

	// Generate/get a list of set symbols
	setSymbols, err := setupSetSymbols(context, scryfallClient, config.SetSymbolDirectory, config.SetSymbolBackgroundPath, session, graph, config.GetNewSetSymbols)
	if err != nil {
		log.Panic(err)
	}

	// Create a Tesseract client to OCR the text
	tessClient := gosseract.NewClient()
	defer tessClient.Close()
	tessClient.SetLanguage("eng")                 // Set the language to English (you may need to install this language for Tesseract)
	tessClient.SetPageSegMode(gosseract.PSM_AUTO) // Set the segmentation mode

	// Setup a SURF context
	surf := contrib.NewSURF()
	defer surf.Close()

	imageBytes, err := ioutil.ReadFile(config.ImageInputPath)
	if err != nil {
		log.Panic(err)
	}

	pcd, err := inferPartialCardData(imageBytes, session, graph, tessClient)
	if err != nil {
		log.Panic(err)
	}

	log.Println(pcd.String())

	card, err := pcd.findClosestCard(context, scryfallClient, surf, setSymbols, config.SetSymbolDirectory)
	if err != nil {
		log.Panic(err)
	}

	log.Println(card)
}
