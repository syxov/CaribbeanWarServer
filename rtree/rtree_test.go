package rtree

import (
	"math"
	"testing"
)

func TestIntersection(t *testing.T) {
	var (
		r1, r2 *Rect
	)
	r1, r2 = NewRect(Point{1, 1}, []float64{2, 2}, 0), NewRect(Point{2, 2}, []float64{2, 2}, 0)
	if r1.intersectRect(r2) == false {
		t.Error("1 test failing")
	}
	r1, r2 = NewRect(Point{1, 1}, []float64{2, 2}, math.Pi/4), NewRect(Point{2, 2}, []float64{2, 2}, 0)
	if r1.intersectRect(r2) == false {
		t.Error("2 test failing")
	}
	r1, r2 = NewRect(Point{1, 1}, []float64{2, 2}, math.Pi/4), NewRect(Point{2, 2}, []float64{2, 2}, math.Pi/4)
	if r1.intersectRect(r2) == false {
		t.Error("3 test failing")
	}
	r1, r2 = NewRect(Point{1, 1}, []float64{2, 2}, math.Pi/4), NewRect(Point{3, 3}, []float64{2, 2}, 0)
	if r1.intersectRect(r2) == true {
		t.Error("4 test failing")
	}
	r1, r2 = NewRect(Point{1, 1}, []float64{3, 3}, 0), NewRect(Point{2, 2}, []float64{1, 1}, math.Pi/4)
	if r1.intersectRect(r2) == false {
		t.Error("5 test failing")
	}
	r1, r2 = NewRect(Point{2, 2}, []float64{1, 1}, math.Pi/4), NewRect(Point{1, 1}, []float64{3, 3}, 0)
	if r1.intersectRect(r2) == false {
		t.Error("6 test failing")
	}
}
