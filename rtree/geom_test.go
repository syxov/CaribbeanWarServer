package rtree

import (
	"testing"
)

func TestGeom(t *testing.T) {
	if sameSign(-1, 2) {
		t.Error("Expected that -1 and 2 not has same signs")
	}
	if sameSign(1, -2) {
		t.Error("Expected that 1 and -2 not has same signs")
	}
	if !sameSign(1, 2) {
		t.Error("Expected that 1 and 2 has same signs")
	}
	if !sameSign(-1, -2) {
		t.Error("Expected that 1 and 2 has same signs")
	}
	if !sameSign(0, 0) {
		t.Error("Expected that zeroes are equal")
	}
}
