package main

import (
	"github.com/MetalBlueberry/go-plotly-express/pkg/express"
	"github.com/MetalBlueberry/go-plotly/pkg/offline"
)

func main() {
	X := []int{1, 2, 3, 4, 5, 6}
	Y := []float64{2, 1, 1, 2, 3, 3}
	Cat := []string{"a", "a", "b", "b", "c", "c"}

	fig := express.NewScatter(X, Y).
		// WithXTitle("letters").
		// WithYTitle("numbers").
		// WithAnimationFrames(Cat, []int{0, 8}, []float64{0, 4}).
		WithCategories(Cat).
		// WithColor([]types.ColorWithColorScale{
		// 	types.UseColor("red"),
		// 	types.UseColor("red"),
		// 	types.UseColor("red"),
		// 	types.UseColor("blue"),
		// 	types.UseColor("blue"),
		// 	types.UseColor("blue"),
		// }).
		// WithSize([]types.NumberType{
		// 	types.N(10),
		// 	types.N(10),
		// 	types.N(10),
		// 	types.N(20),
		// 	types.N(20),
		// 	types.N(20),
		// }).
		WithTitle("Playground for plotly express").
		Fig()

	offline.Serve(fig)

}
