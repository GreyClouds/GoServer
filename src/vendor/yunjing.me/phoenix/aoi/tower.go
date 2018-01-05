package aoi

type Tower struct {
	objects  Set
	watchers Set
}

func newTower() *Tower {
	return &Tower{
		objects:  NewUnsafeSet(),
		watchers: NewUnsafeSet(),
	}
}

func (tower *Tower) Add(id uint32) bool {
	if tower.objects.Contains(id) {
		return false
	}

	tower.objects.Add(id)

	return true
}

func (tower *Tower) Remove(id uint32) bool {
	if !tower.objects.Contains(id) {
		return false
	}

	tower.objects.Remove(id)

	return true
}

func (tower *Tower) GetObjects() []uint32 {
	return tower.objects.Values()
}

func (tower *Tower) AddWatcher(id uint32) bool {
	if tower.watchers.Contains(id) {
		return false
	}

	tower.watchers.Add(id)

	return true
}

func (tower *Tower) RemoveWatcher(id uint32) bool {
	if !tower.watchers.Contains(id) {
		return false
	}

	tower.watchers.Remove(id)

	return true
}

func (tower *Tower) GetWatchers() []uint32 {
	return tower.watchers.Values()
}

func (tower *Tower) GetDiffWatchers(other *Tower) []uint32 {
	return tower.watchers.Sub(other.watchers).Values()
}

func (tower *Tower) GetSameWatchers(other *Tower) []uint32 {
	return tower.watchers.Sum(other.watchers).Values()
}
