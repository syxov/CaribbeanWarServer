package structs

import (
	"CaribbeanWarServer/point"
	"CaribbeanWarServer/rtree"
	"CaribbeanWarServer/testUtil"
	"encoding/json"
	"math"
	"testing"
)

func TestCore(t *testing.T) {
	core := NewCore(point.Point3D{0, 0, 0}, math.Pi/4, 0, -1, 0)
	rect := core.GetBounds()
	expectedValue := rtree.NewRect(rtree.Point{-0.5, -0.5}, []float64{1, 1})
	if !rectEqual(rect, expectedValue) {
		expectedValueMarshal, _ := json.Marshal(expectedValue)
		rectMarshal, _ := json.Marshal(rect)
		t.Error("Expected that before update", string(rectMarshal), "is equal to", string(expectedValueMarshal))
	}
	core.UpdatePosition(1)
	rect = core.GetBounds()
	expectedValue = rtree.NewRect(rtree.Point{-0.5, -0.5}, []float64{101, 1})
	if !rectEqual(rect, expectedValue) {
		expectedValueMarshal, _ := json.Marshal(expectedValue)
		rectMarshal, _ := json.Marshal(rect)
		t.Error("Expected that after update", string(rectMarshal), "is equal to", string(expectedValueMarshal))
	}
}

func rectEqual(r1, r2 *rtree.Rect) bool {
	return testUtil.FloatEqual(r1.P[0], r2.P[0]) &&
		testUtil.FloatEqual(r1.P[1], r2.P[1]) &&
		testUtil.FloatEqual(r1.Q[0], r2.Q[0]) &&
		testUtil.FloatEqual(r1.Q[1], r2.Q[1])
}
