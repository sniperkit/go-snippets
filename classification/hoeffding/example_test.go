package hoeffding_test

import (
	"fmt"

	"github.com/bsm/reason/classification/hoeffding"
	common "github.com/bsm/reason/common/hoeffding"
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
	tree, err := hoeffding.New(model, "play", &hoeffding.Config{
		Config: common.Config{
			GracePeriod:     2,
			SplitConfidence: 0.1,
		},
	})
	if err != nil {
		panic(err)
	}

	// Train
	for _, x := range examples {
		tree.Train(x, 1.0)
	}

	// Predict
	prediction := tree.Predict(nil, core.MapExample{
		"outlook":  "rainy",
		"temp":     "mild",
		"humidity": "high",
		"windy":    "false",
	}).Best()

	// Print categories with probabilities
	fmt.Printf("yes: %.2f\n", prediction.P(0))
	fmt.Printf(" no: %.2f\n", prediction.P(1))

	// Output:
	// yes: 0.40
	//  no: 0.60
}
