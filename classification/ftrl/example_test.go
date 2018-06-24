package ftrl_test

import (
	"fmt"

	"github.com/bsm/reason/classification/ftrl"
	"github.com/bsm/reason/core"
)

func Example() {
	model := core.NewModel(
		core.NewCategoricalFeature("play", []string{"yes", "no"}),
		core.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"}),
		core.NewCategoricalFeature("temp", []string{"hot", "mild", "cool"}),
		core.NewCategoricalFeature("humidity", []string{"normal", "high"}),
		core.NewCategoricalFeature("windy", []string{"true", "false"}),
	)

	examples := []core.MapExample{
		{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "false", "play": "no"},
		{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "true", "play": "no"},
		{"outlook": "overcast", "temp": "hot", "humidity": "high", "windy": "false", "play": "yes"},
		{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "false", "play": "yes"},
		{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "false", "play": "yes"},
		{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "true", "play": "no"},
		{"outlook": "overcast", "temp": "cool", "humidity": "normal", "windy": "true", "play": "yes"},
		{"outlook": "rainy", "temp": "mild", "humidity": "high", "windy": "false", "play": "no"},
		{"outlook": "rainy", "temp": "cool", "humidity": "normal", "windy": "false", "play": "yes"},
		{"outlook": "sunny", "temp": "mild", "humidity": "normal", "windy": "false", "play": "yes"},
		{"outlook": "rainy", "temp": "mild", "humidity": "normal", "windy": "true", "play": "yes"},
		{"outlook": "overcast", "temp": "mild", "humidity": "high", "windy": "true", "play": "yes"},
		{"outlook": "overcast", "temp": "hot", "humidity": "normal", "windy": "false", "play": "yes"},
		{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "true", "play": "no"},
	}

	// Init with a model
	opt, err := ftrl.New(model, "play", nil)
	if err != nil {
		panic(err)
	}

	// Train
	for epoch := 0; epoch < 2000; epoch++ {
		for _, x := range examples {
			opt.Train(x, 1.0)
		}
	}

	// Predict
	prediction := opt.Predict(core.MapExample{
		"outlook":  "rainy",
		"temp":     "mild",
		"humidity": "high",
		"windy":    "false",
	})

	// Print categories with probabilities
	fmt.Printf("yes: %.2f\n", 1-prediction)
	fmt.Printf(" no: %.2f\n", prediction)

	// Output:
	// yes: 0.41
	//  no: 0.59
}
