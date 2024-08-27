package express

import (
	"fmt"
	"reflect"

	grob "github.com/MetalBlueberry/go-plotly/generated/v2.31.1/graph_objects"
	"github.com/MetalBlueberry/go-plotly/pkg/types"
	"golang.org/x/exp/constraints"
)

type Scatter[XType any, YType any] struct {
	title  string
	titleX string
	titleY string
	x      []XType
	y      []YType

	categories      []string
	color           []types.ColorWithColorScale
	size            []types.NumberType
	animationFrames []string
	xRange          []XType
	yRange          []YType
}

func NewScatter[XType any, YType any](X []XType, Y []YType) *Scatter[XType, YType] {
	return &Scatter[XType, YType]{
		x: X,
		y: Y,
	}
}

func (s *Scatter[XType, YType]) WithXTitle(title string) *Scatter[XType, YType] {
	s.titleX = title
	return s
}
func (s *Scatter[XType, YType]) WithYTitle(title string) *Scatter[XType, YType] {
	s.titleY = title
	return s
}

func (s *Scatter[XType, YType]) WithTitle(title string) *Scatter[XType, YType] {
	s.title = title
	return s
}

func (s *Scatter[XType, YType]) WithCategories(categories []string) *Scatter[XType, YType] {
	s.categories = categories
	return s
}
func (s *Scatter[XType, YType]) WithColor(color []types.ColorWithColorScale) *Scatter[XType, YType] {
	s.color = color
	return s
}

func (s *Scatter[XType, YType]) WithSize(size []types.NumberType) *Scatter[XType, YType] {
	s.size = size
	return s
}

func (s *Scatter[XType, YType]) WithAnimationFrames(frames []string, xRange []XType, yRange []YType) *Scatter[XType, YType] {
	s.animationFrames = frames
	s.xRange = xRange
	s.yRange = yRange
	return s
}

func generateTraces[XType any, YType any](x []XType, y []YType, categories []string, color []types.ColorWithColorScale, size []types.NumberType) []types.Trace {

	keys := []string{""}
	xCategories := [][]XType{x}
	yCategories := [][]YType{y}
	colorCategories := [][]types.ColorWithColorScale{color}
	sizeCategories := [][]types.NumberType{size}

	if categories != nil {
		indices := [][]int{}
		indices, keys = findIndices(categories)
		xCategories = splitByIndices(x, indices, keys)
		yCategories = splitByIndices(y, indices, keys)
		colorCategories = splitByIndices(color, indices, keys)
		sizeCategories = splitByIndices(size, indices, keys)
	}

	traces := []types.Trace{}
	for i := range keys {
		trace := &grob.Scatter{
			X:    types.DataArray(xCategories[i]),
			Y:    types.DataArray(yCategories[i]),
			Name: types.S(keys[i]),
		}

		if size != nil {
			if trace.Marker == nil {
				trace.Marker = &grob.ScatterMarker{}
			}
			trace.Marker.Size = types.ArrayOKArray(sizeCategories[i]...)
		}
		if color != nil {
			if trace.Marker == nil {
				trace.Marker = &grob.ScatterMarker{}
			}
			trace.Marker.Color = types.ArrayOKArray(colorCategories[i]...)
		}
		traces = append(traces, trace)
	}
	return traces
}

func (s *Scatter[XType, YType]) Fig() *grob.Fig {

	keys := []string{""}
	xFrames := [][]XType{s.x}
	yFrames := [][]YType{s.y}
	categoryFrames := [][]string{s.categories}
	colorFrames := [][]types.ColorWithColorScale{s.color}
	sizeFrames := [][]types.NumberType{s.size}

	if s.animationFrames != nil {
		indices := [][]int{}
		indices, keys = findIndices(s.animationFrames)
		xFrames = splitByIndices(s.x, indices, keys)
		yFrames = splitByIndices(s.y, indices, keys)
		categoryFrames = splitByIndices(s.categories, indices, keys)
		colorFrames = splitByIndices(s.color, indices, keys)
		sizeFrames = splitByIndices(s.size, indices, keys)
	}

	frames := []grob.Frame{}
	for i := range keys {
		traces := generateTraces(xFrames[i], yFrames[i], categoryFrames[i], colorFrames[i], sizeFrames[i])
		frames = append(frames, grob.Frame{
			Name: types.S(keys[i]),
			Data: traces,
		})
	}

	fig := &grob.Fig{
		Data: frames[0].Data,
		Layout: &grob.Layout{
			Xaxis: &grob.LayoutXaxis{
				Title: &grob.LayoutXaxisTitle{
					Text: types.StringType(s.titleX),
				},
			},
			Yaxis: &grob.LayoutYaxis{
				Title: &grob.LayoutYaxisTitle{
					Text: types.StringType(s.titleY),
				},
			},
		},
	}
	if s.title != "" {
		fig.Layout.Title = &grob.LayoutTitle{
			Text: types.StringType(s.title),
		}
	}

	if len(keys) > 1 {
		fig.Frames = frames
		fig.Layout.Updatemenus = []grob.LayoutUpdatemenu{
			{
				Type:       grob.LayoutUpdatemenuTypeButtons,
				Showactive: types.False,
				Buttons: []grob.LayoutUpdatemenuButton{
					{
						Label:  types.S("Play"),
						Method: grob.LayoutUpdatemenuButtonMethodAnimate,
						Args: []*map[string]interface{}{
							nil,
							{
								"mode":        "immediate",
								"fromcurrent": false,
							},
						},
					},
				},
			},
		}
		fig.Layout.Xaxis.Range = s.xRange
		fig.Layout.Yaxis.Range = s.yRange
	}

	return fig
}

func iterateOverSlice(data interface{}, iter func(int, interface{})) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		panic(fmt.Errorf("data is not a slice, is a %T", data))
	}

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		iter(i, elem)
	}

	return nil
}

func findIndices[T constraints.Ordered](input []T) ([][]int, []T) {
	indexMap := make(map[T][]int)
	var keys []T

	// Populate the map with indices grouped by the value
	for i, val := range input {
		if _, found := indexMap[val]; !found {
			keys = append(keys, val)
		}
		indexMap[val] = append(indexMap[val], i)
	}

	// Collect the grouped indices into a result slice in the order of first appearance
	var result [][]int
	for _, key := range keys {
		result = append(result, indexMap[key])
	}

	return result, keys
}

func splitByIndices[T any, Y constraints.Ordered](orginal []T, indices [][]int, keys []Y) [][]T {
	if len(keys) != len(indices) {
		panic("Length doesn't match")
	}
	if orginal == nil {
		return make([][]T, len(keys))
	}

	result := [][]T{}
	for i, section := range indices {
		result = append(result, []T{})
		for _, value := range section {
			result[i] = append(result[i], orginal[value])
		}
	}

	return result
}

func ArrayOKAppend[T any](array *types.ArrayOK[*T], values ...T) *types.ArrayOK[*T] {
	for _, el := range values {
		array.Array = append(array.Array, &el)
	}
	return array
}
