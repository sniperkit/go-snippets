package hoeffding_test

import (
	"fmt"

	common "github.com/bsm/reason/common/hoeffding"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/regression/hoeffding"
)

func Example() {
	model := core.NewModel(
		core.NewNumericalFeature("hours"),
		core.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"}),
		core.NewCategoricalFeature("temp", []string{"hot", "mild", "cool"}),
		core.NewCategoricalFeature("humidity", []string{"normal", "high"}),
		core.NewCategoricalFeature("windy", []string{"true", "false"}),
	)

	examples := []core.MapExample{
		{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "false", "hours": 25},
		{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "true", "hours": 30},
		{"outlook": "overcast", "temp": "hot", "humidity": "high", "windy": "false", "hours": 46},
		{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "false", "hours": 45},
		{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "false", "hours": 52},
		{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "true", "hours": 23},
		{"outlook": "overcast", "temp": "cool", "humidity": "normal", "windy": "true", "hours": 43},
		{"outlook": "rainy", "temp": "mild", "humidity": "high", "windy": "false", "hours": 35},
		{"outlook": "rainy", "temp": "cool", "humidity": "normal", "windy": "false", "hours": 38},
		{"outlook": "sunny", "temp": "mild", "humidity": "normal", "windy": "false", "hours": 46},
		{"outlook": "rainy", "temp": "mild", "humidity": "normal", "windy": "true", "hours": 48},
		{"outlook": "overcast", "temp": "mild", "humidity": "high", "windy": "true", "hours": 52},
		{"outlook": "overcast", "temp": "hot", "humidity": "normal", "windy": "false", "hours": 44},
		{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "true", "hours": 30},
	}

	// Init with a model
	tree, err := hoeffding.New(model, "hours", &hoeffding.Config{
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

	// Print mean value with weight
	fmt.Printf("mean: %.2f, weight: %.0f\n", prediction.Mean(), prediction.Weight)

	// Output:
	// mean: 42.67, weight: 6
}
