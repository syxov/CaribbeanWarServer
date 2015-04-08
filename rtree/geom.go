// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"math"
)

// Point represents a point in n-dimensional Euclidean space.
type Point []float64

// minDist computes the square of the distance from a point to a rectangle.
// If the point is contained in the rectangle then the distance is zero.
//
// Implemented per Definition 2 of "Nearest Neighbor Queries" by
// N. Roussopoulos, S. Kelley and F. Vincent, ACM SIGMOD, pages 71-79, 1995.
func (p Point) minDist(r *Rect) (sum float64) {
	for i, pi := range p {
		if pi < r.P[i] {
			sum += math.Pow(pi-r.P[i], 2)
		} else if pi > r.Q[i] {
			d := math.Pow(pi-r.Q[i], 2)
			sum += d * d
		}
	}
	return
}

// minMaxDist computes the minimum of the maximum distances from p to points
// on r.  If r is the bounding box of some geometric objects, then there is
// at least one object contained in r within minMaxDist(p, r) of p.
//
// Implemented per Definition 4 of "Nearest Neighbor Queries" by
// N. Roussopoulos, S. Kelley and F. Vincent, ACM SIGMOD, pages 71-79, 1995.
func (p Point) minMaxDist(r *Rect) float64 {
	// by definition, MinMaxDist(p, r) =
	// min{1<=k<=n}(|pk - rmk|^2 + sum{1<=i<=n, i != k}(|pi - rMi|^2))
	// where rmk and rMk are defined as follows:

	rm := func(k int) float64 {
		if p[k] <= (r.P[k]+r.Q[k])/2 {
			return r.P[k]
		}
		return r.Q[k]
	}

	rM := func(k int) float64 {
		if p[k] >= (r.P[k]+r.Q[k])/2 {
			return r.P[k]
		}
		return r.Q[k]
	}

	// This formula can be computed in linear time by precomputing
	// S = sum{1<=i<=n}(|pi - rMi|^2).

	S := 0.0
	for i := range p {
		d := p[i] - rM(i)
		S += d * d
	}

	// Compute MinMaxDist using the precomputed S.
	min := math.MaxFloat64
	for k := range p {
		d1 := p[k] - rM(k)
		d2 := p[k] - rm(k)
		d := S - d1*d1 + d2*d2
		if d < min {
			min = d
		}
	}

	return min
}

// Rect represents a subset of n-dimensional Euclidean space of the form
// [a1, b1] x [a2, b2] x ... x [an, bn], where ai < bi for all 1 <= i <= n.
type Rect struct {
	P, Q  Point // Enforced by NewRect: p[i] <= q[i] for all i.
	angle float64
	size  float64
}

// NewRect constructs and returns a pointer to a Rect given a corner point and
// the lengths of each dimension.  The point p should be the most-negative point
// on the rectangle (in every dimension) and every length should be positive.
func NewRect(p Point, lengths []float64, angle ...float64) *Rect {
	r := new(Rect)
	r.P = p
	r.Q = make([]float64, 2, 2)
	r.Q[0] = p[0] + lengths[0]
	r.Q[1] = p[1] + lengths[1]
	r.size = (r.Q[0] - r.P[0]) * (r.Q[1] - r.P[1])
	if len(angle) != 0 {
		r.angle = angle[0]
	}
	return r
}

// containsRect tests whether r2 is is located inside r1.
func (r1 *Rect) containsRect(r2 *Rect) bool {
	for i, a1 := range r1.P {
		b1, a2, b2 := r1.Q[i], r2.P[i], r2.Q[i]
		// enforced by constructor: a1 <= b1 and a2 <= b2.
		// so containment holds if and only if a1 <= a2 <= b2 <= b1
		// for every dimension.
		if a1 > a2 || b2 > b1 {
			return false
		}
	}

	return true
}

type point struct {
	x, y float64
}

type rectangle []point

const length = 4

func (r1 *Rect) intersectRect(r2 *Rect) bool {
	rect1, rect2 := *(r1.ToRectangle()), *(r2.ToRectangle())
	for _, p := range rect1 {
		inReactangle := true
		ali := align(&rect2[0], &rect2[1], &p)
		for i := 1; i < length; i++ {
			ali_t := align(&rect2[i], &rect2[(i+1)%length], &p)
			if !sameSign(ali_t, ali) {
				inReactangle = false
				break
			}
		}
		if inReactangle {
			return true
		}
	}
	return false
}

func sameSign(_x, _y float64) bool {
	x, y := int(_x), int(_y)
	t := (x ^ y) >> 31
	return ((x + t) ^ t) == x
}

func align(a, b, c *point) float64 {
	return a.x*(b.y-c.y) + a.y*(c.x-b.x) + b.x*c.y - c.x*b.y
}

func (self *Rect) ToRectangle() *rectangle {
	return &rectangle{
		rotatePoint(&point{x: self.P[0], y: self.P[1]}, self.angle),
		rotatePoint(&point{x: self.P[0] + self.Q[0], y: self.P[1]}, self.angle),
		rotatePoint(&point{x: self.P[0] + self.Q[0], y: self.P[1] + self.Q[1]}, self.angle),
		rotatePoint(&point{x: self.P[0], y: self.P[1] + self.Q[1]}, self.angle),
	}

}

func rotatePoint(p *point, angle float64) point {
	return point{x: p.x*math.Cos(angle) - p.y*math.Sin(angle), y: p.x*math.Sin(angle) + p.y*math.Cos(angle)}
}

// intersect computes the intersection of two rectangles.  If no intersection
// exists, the intersection is nil.
const dim = 2

func intersect(r1, r2 *Rect) *Rect {
	// There are four cases of overlap:
	//
	//     1.  a1------------b1
	//              a2------------b2
	//              p--------q
	//
	//     2.       a1------------b1
	//         a2------------b2
	//              p--------q
	//
	//     3.  a1-----------------b1
	//              a2-------b2
	//              p--------q
	//
	//     4.       a1-------b1
	//         a2-----------------b2
	//              p--------q
	//
	// Thus there are only two cases of non-overlap:
	//
	//     1. a1------b1
	//                    a2------b2
	//
	//     2.             a1------b1
	//        a2------b2
	//
	// Enforced by constructor: a1 <= b1 and a2 <= b2.  So we can just
	// check the endpoints.

	p := make([]float64, dim)
	q := make([]float64, dim)
	for i := range p {
		a1, b1, a2, b2 := r1.P[i], r1.Q[i], r2.P[i], r2.Q[i]
		if b2 <= a1 || b1 <= a2 {
			return nil
		}
		p[i] = math.Max(a1, a2)
		q[i] = math.Min(b1, b2)
	}
	return NewRect(p, q)
}

// ToRect constructs a rectangle containing p with side lengths 2*tol.
func (p Point) ToRect(tol float64) *Rect {
	a, b := make([]float64, dim), make([]float64, dim)
	for i := range p {
		a[i] = p[i] - tol
		b[i] = p[i] + tol
	}
	return NewRect(a, b)
}

// boundingBox constructs the smallest rectangle containing both r1 and r2.
func boundingBox(r1, r2 *Rect) (bb *Rect) {
	bb = new(Rect)
	bb.P, bb.Q = make([]float64, dim), make([]float64, dim)
	for i := 0; i < dim; i++ {
		if r1.P[i] <= r2.P[i] {
			bb.P[i] = r1.P[i]
		} else {
			bb.P[i] = r2.P[i]
		}
		if r1.Q[i] <= r2.Q[i] {
			bb.Q[i] = r2.Q[i]
		} else {
			bb.Q[i] = r1.Q[i]
		}
	}
	return
}

// boundingBoxN constructs the smallest rectangle containing all of r...
func boundingBoxN(rects ...*Rect) (bb *Rect) {
	if len(rects) == 1 {
		bb = rects[0]
		return
	}
	bb = boundingBox(rects[0], rects[1])
	for _, rect := range rects[2:] {
		bb = boundingBox(bb, rect)
	}
	return
}
