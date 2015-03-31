package structs

import (
	"CaribbeanWarServer/point"
	"CaribbeanWarServer/rtree"
	"encoding/json"
	"math"
	"testing"
)

func TestCore(t *testing.T) {
	core := NewCore(point.Point3D{0, 10, 0}, 0, math.Pi/4, -1, 0)
	rect := core.GetBounds()
	expectedValue := rtree.NewRect(rtree.Point{-0.5, -0.5}, []float64{1, 1})
	if !rectEqual(rect, expectedValue) {
		expectedValueMarshal, _ := json.Marshal(expectedValue)
		rectMarshal, _ := json.Marshal(rect)
		t.Error("Expected that", string(rectMarshal), "is equal to", string(expectedValueMarshal))
	}
	core.UpdatePosition(1)
	rect = core.GetBounds()
	expectedValue = rtree.NewRect(rtree.Point{-0.5, -0.5}, []float64{101, 1})
	if !rectEqual(rect, expectedValue) {
		expectedValueMarshal, _ := json.Marshal(expectedValue)
		rectMarshal, _ := json.Marshal(rect)
		t.Error("Expected that", string(rectMarshal), "is equal to", string(expectedValueMarshal))
	}
}

func rectEqual(r1, r2 *rtree.Rect) bool {
	return floatEqual(r1.P[0], r2.P[0]) &&
		floatEqual(r1.P[1], r2.P[1]) &&
		floatEqual(r1.Q[0], r2.Q[0]) &&
		floatEqual(r1.Q[1], r2.Q[1])
}

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.0000001
}
