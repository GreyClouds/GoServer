package aoi

import (
	"testing"
)

func TestAddAndRemoveObject(t *testing.T) {
	tower := newTower()
	tower.Add(1)
	tower.Remove(1)
	tower.Add(2)
	if len(tower.GetObjects()) != 1 {
		t.Errorf("Expected 1: %d", len(tower.GetObjects()))
	}
}
