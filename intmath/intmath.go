package intmath

func Doz(x, y int) int {
	if x > y {
		return x - y
	}
	return 0
}

func Max(x, y int) int {
	return y + Doz(x, y)
}

func Min(x, y int) int {
	return x - Doz(x, y)
}

func Doz32(x, y int32) int32 {
	if x > y {
		return x - y
	}
	return 0
}

func Max32(x, y int32) int32 {
	return y + Doz32(x, y)
}

func Min32(x, y int32) int32 {
	return x - Doz32(x, y)
}
