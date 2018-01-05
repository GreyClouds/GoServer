package aoi

import (
	"math"
	"testing"
)

func TestGetPosLimit1(t *testing.T) {
	start, end := getPosLimit(Point{X: 1, Y: 1}, &Point{X: 12, Y: 321}, 4)
	if start.X != 0 || start.Y != 0 || end.X != 8 || end.Y != 8 {
		t.Errorf("Expected (0, 0) ~ (8, 8): (%d, %d) ~ (%d, %d)", start.X, start.Y, end.X, end.Y)
	}
}

func TestGetPosLimit2(t *testing.T) {
	start, end := getPosLimit(Point{X: 8, Y: 4}, &Point{X: 12, Y: 321}, 4)
	if start.X != 4 || start.Y != 0 || end.X != 12 || end.Y != 8 {
		t.Errorf("Expected (4, 0) ~ (8, 8): (%d, %d) ~ (%d, %d)", start.X, start.Y, end.X, end.Y)
	}
}

func TestNewTowerNum(t *testing.T) {
	for i := int32(0); i <= int32(1000000); i++ {
		calc := int32((math.Ceil(float64(i+1) / float64(65536))))
		num := int32(i>>16) + 1
		if num != calc {
			t.Errorf("Expected %d: %d", calc, num)
		}
	}
}
