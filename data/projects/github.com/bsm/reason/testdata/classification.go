package testdata

import (
	"github.com/bsm/reason/core"
)

type ClassificationScore struct {
	Accuracy, Kappa, LogLoss float64
}

func ClassificationModel() *core.Model {
	return core.NewModel(
		core.NewCategoricalFeature("play", []string{"yes", "no"}),
		core.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"}),
		core.NewCategoricalFeature("temp", []string{"hot", "mild", "cool"}),
		core.NewCategoricalFeature("humidity", []string{"high", "normal"}),
		core.NewCategoricalFeature("windy", []string{"true", "false"}),
	)
}

func ClassificationData() []core.Example {
	return []core.Example{
		core.MapExample{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "false", "play": "no"},
		core.MapExample{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "true", "play": "no"},
		core.MapExample{"outlook": "overcast", "temp": "hot", "humidity": "high", "windy": "false", "play": "yes"},
		core.MapExample{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "false", "play": "yes"},
		core.MapExample{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "false", "play": "yes"},
		core.MapExample{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "true", "play": "no"},
		core.MapExample{"outlook": "overcast", "temp": "cool", "humidity": "normal", "windy": "true", "play": "yes"},
		core.MapExample{"outlook": "rainy", "temp": "mild", "humidity": "high", "windy": "false", "play": "no"},
		core.MapExample{"outlook": "rainy", "temp": "cool", "humidity": "normal", "windy": "false", "play": "yes"},
		core.MapExample{"outlook": "sunny", "temp": "mild", "humidity": "normal", "windy": "false", "play": "yes"},
		core.MapExample{"outlook": "rainy", "temp": "mild", "humidity": "normal", "windy": "true", "play": "yes"},
		core.MapExample{"outlook": "overcast", "temp": "mild", "humidity": "high", "windy": "true", "play": "yes"},
		core.MapExample{"outlook": "overcast", "temp": "hot", "humidity": "normal", "windy": "false", "play": "yes"},
		core.MapExample{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "true", "play": "no"},
	}
}
