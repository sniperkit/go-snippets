package testdata

import "github.com/bsm/reason/core"

type RegressionScore struct {
	R2, RMSE float64
}

func RegressionModel() *core.Model {
	return core.NewModel(
		core.NewNumericalFeature("hours"),
		core.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"}),
		core.NewCategoricalFeature("temp", []string{"hot", "mild", "cool"}),
		core.NewNumericalFeature("humidity"),
		core.NewCategoricalFeature("windy", []string{"true", "false"}),
	)
}

func RegressionData() []core.Example {
	return []core.Example{
		core.MapExample{"outlook": "rainy", "temp": "hot", "humidity": 60, "windy": "false", "hours": 25},
		core.MapExample{"outlook": "rainy", "temp": "hot", "humidity": 60, "windy": "true", "hours": 30},
		core.MapExample{"outlook": "overcast", "temp": "hot", "humidity": 60, "windy": "false", "hours": 46},
		core.MapExample{"outlook": "sunny", "temp": "mild", "humidity": 60, "windy": "false", "hours": 45},
		core.MapExample{"outlook": "sunny", "temp": "cool", "humidity": 40, "windy": "false", "hours": 52},
		core.MapExample{"outlook": "sunny", "temp": "cool", "humidity": 40, "windy": "true", "hours": 23},
		core.MapExample{"outlook": "overcast", "temp": "cool", "humidity": 40, "windy": "true", "hours": 43},
		core.MapExample{"outlook": "rainy", "temp": "mild", "humidity": 60, "windy": "false", "hours": 35},
		core.MapExample{"outlook": "rainy", "temp": "cool", "humidity": 40, "windy": "false", "hours": 38},
		core.MapExample{"outlook": "sunny", "temp": "mild", "humidity": 40, "windy": "false", "hours": 46},
		core.MapExample{"outlook": "rainy", "temp": "mild", "humidity": 40, "windy": "true", "hours": 48},
		core.MapExample{"outlook": "overcast", "temp": "mild", "humidity": 60, "windy": "true", "hours": 52},
		core.MapExample{"outlook": "overcast", "temp": "hot", "humidity": 40, "windy": "false", "hours": 44},
		core.MapExample{"outlook": "sunny", "temp": "mild", "humidity": 60, "windy": "true", "hours": 30},
	}
}
