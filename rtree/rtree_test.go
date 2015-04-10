package rtree

import (
	"CaribbeanWarServer/testUtil"
	"math"
	"testing"
)

func TestIntersection(t *testing.T) {
	var (
		r1, r2 *Rect
	)
	r1, r2 = NewRect(Point{1, 1}, []float64{2, 2}, 0), NewRect(Point{2, 2}, []float64{2, 2}, 0)
	if r1.intersectRect(r2) == false {
		t.Error("1 test failed")
	}
	r1, r2 = NewRect(Point{1, 1}, []float64{2, 2}, math.Pi/4), NewRect(Point{2, 2}, []float64{2, 2}, 0)
	if r1.intersectRect(r2) == false {
		t.Error("2 test failed")
	}
	r1, r2 = NewRect(Point{1, 1}, []float64{2, 2}, math.Pi/4), NewRect(Point{2, 2}, []float64{2, 2}, math.Pi/4)
	if r1.intersectRect(r2) == false {
		t.Error("3 test failed")
	}
	r1, r2 = NewRect(Point{1, 1}, []float64{2, 2}, math.Pi/4), NewRect(Point{3, 3}, []float64{2, 2}, 0)
	if r1.intersectRect(r2) == true {
		t.Error("4 test failed")
	}
	r1, r2 = NewRect(Point{1, 1}, []float64{3, 3}, 0), NewRect(Point{2, 2}, []float64{1, 1}, math.Pi/4)
	if r1.intersectRect(r2) == false {
		t.Error("5 test failed")
	}
	r1, r2 = NewRect(Point{2, 2}, []float64{1, 1}, math.Pi/4), NewRect(Point{1, 1}, []float64{3, 3}, 0)
	if r1.intersectRect(r2) == false {
		t.Error("6 test failed")
	}
}

func TestConvertionToRectangle(t *testing.T) {
	rect, expected := *NewRect(Point{2, 2}, []float64{1, 1}, 0.785398163).ToRectangle(), rectangle{point{2.5, 1.79289}, point{3.20711, 2.5}, point{2.5, 3.20711}, point{1.79289, 2.5}}
	if !rectangleAreEqual(rect, expected) {
		t.Error("Conversion test failed. Expect", rect, "to equal", expected)
	}
}

func rectangleAreEqual(r1, r2 rectangle) bool {
	for i := range r1 {
		if !pointAreEqual(r1[i], r2[i]) {
			return false
		}
	}
	return true
}

func pointAreEqual(p1, p2 point) bool {
	return testUtil.FloatEqual(p1.x, p2.x) && testUtil.FloatEqual(p1.y, p2.y)
}
