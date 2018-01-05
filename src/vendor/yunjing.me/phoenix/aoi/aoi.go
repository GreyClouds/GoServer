package aoi

const (
	kMaxRangeV = 5
)

type TowerAOI struct {
	width       int32
	height      int32
	towerWidth  uint8
	towerHeight uint8
	towers      [][]*Tower
	max         *Point
	listener    Listener
}

func New(width, height int32, towerWidth, towerHeight uint8, listener Listener) *TowerAOI {
	aoi := &TowerAOI{
		width:       width,
		height:      height,
		towerWidth:  towerWidth,
		towerHeight: towerHeight,
		towers:      [][]*Tower{},
		listener:    listener,
	}

	aoi.initialize()

	return aoi
}

func (aoi *TowerAOI) initialize() {
	aoi.max = &Point{
		X: int32(aoi.width>>aoi.towerWidth) + 1,
		Y: int32(aoi.height>>aoi.towerHeight) + 1,
	}

	for i := int32(0); i <= aoi.max.X; i++ {
		group := []*Tower{}
		for j := int32(0); j <= aoi.max.Y; j++ {
			group = append(group, newTower())
		}
		aoi.towers = append(aoi.towers, group)
	}
}

// 获取指定视野内的特定类型物件
func (self *TowerAOI) GetViewObjects(pos Point, view int32) []uint32 {
	if !self.checkPos(pos) {
		return []uint32{}
	}

	if view > kMaxRangeV {
		view = kMaxRangeV
	}

	result := NewUnsafeSet()
	start, end := getPosLimit(self.transPos(pos), self.max, view)

	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			objs := self.towers[i][j].GetObjects()
			if objs == nil {
				continue
			}

			if len(objs) == 0 {
				continue
			}

			for _, v := range objs {
				result.Add(v)
			}
		}
	}

	return result.Values()
}

// 获取指定位置与半径内的指定类型物体
func (self *TowerAOI) GetObjects(pos Point, radius int32) []uint32 {
	if !self.checkPos(pos) {
		return []uint32{}
	}

	x1 := pos.X - radius
	if x1 < 0 {
		x1 = 0
	}

	y1 := pos.Y - radius
	if y1 < 0 {
		y1 = 0
	}

	x2 := pos.X + radius
	if x2 > self.width {
		x2 = self.width
	}

	y2 := pos.Y + radius
	if y2 > self.height {
		y2 = self.height
	}

	result := NewUnsafeSet()
	pt1, pt2 := self.transPos(Point{X: x1, Y: y1}), self.transPos(Point{X: x2, Y: y2})
	for i := pt1.X; i <= pt2.X; i++ {
		for j := pt1.Y; j <= pt2.Y; j++ {
			objs := self.towers[i][j].GetObjects()
			if objs == nil {
				continue
			}

			if len(objs) == 0 {
				continue
			}

			for _, v := range objs {
				result.Add(v)
			}
		}
	}

	return result.Values()
}

func (aoi *TowerAOI) AddObject(id uint32, pos Point) bool {
	if !aoi.checkPos(pos) {
		return false
	}

	pt := aoi.transPos(pos)
	success := aoi.towers[pt.X][pt.Y].Add(id)
	if success {
		watchers := aoi.towers[pt.X][pt.Y].GetWatchers()
		if len(watchers) > 0 {
			aoi.listener.OnObjectAppear(watchers, id)
		}
	}

	return success
}

func (aoi *TowerAOI) RemoveObject(id uint32, pos Point) bool {
	if !aoi.checkPos(pos) {
		return false
	}

	pt := aoi.transPos(pos)
	success := aoi.towers[pt.X][pt.Y].Remove(id)
	if success {
		watchers := aoi.towers[pt.X][pt.Y].GetWatchers()
		if len(watchers) > 0 {
			aoi.listener.OnObjectDisappear(watchers, id)
		}
	}

	return success
}

func (aoi *TowerAOI) UpdateObject(id uint32, oldPos, newPos Point) bool {
	if !aoi.checkPos(oldPos) || !aoi.checkPos(newPos) {
		return false
	}

	pt1 := aoi.transPos(oldPos)
	pt2 := aoi.transPos(newPos)

	if pt1.X == pt2.X && pt1.Y == pt2.Y {
		return true
	}

	tower1 := aoi.towers[pt1.X][pt1.Y]
	tower2 := aoi.towers[pt2.X][pt2.Y]

	tower1.Remove(id)
	tower2.Add(id)

	if watchers := tower1.GetDiffWatchers(tower2); len(watchers) > 0 {
		aoi.listener.OnObjectDisappear(watchers, id)
	}

	if watchers := tower2.GetDiffWatchers(tower1); len(watchers) > 0 {
		aoi.listener.OnObjectAppear(watchers, id)
	}

	// if watchers := tower1.GetSameWatchers(tower2); len(watchers) > 0 {
	// 	aoi.listener.OnObjectUpdate(watchers, id)
	// }

	return true
}

func (aoi TowerAOI) checkPos(pos Point) bool {
	if pos.X < 0 || pos.Y < 0 || pos.X > aoi.width || pos.Y > aoi.height {
		return false
	}

	return true
}

func (aoi TowerAOI) transPos(pos Point) Point {
	return Point{
		X: (pos.X >> aoi.towerWidth),
		Y: (pos.Y >> aoi.towerHeight),
	}
}

func (aoi *TowerAOI) GetWatchers(pos Point) []uint32 {
	if !aoi.checkPos(pos) {
		return nil
	}

	pt := aoi.transPos(pos)
	return aoi.towers[pt.X][pt.Y].GetWatchers()
}

func (aoi *TowerAOI) AddWatcher(id uint32, pos Point, view int32) {
	if view < 0 {
		return
	}

	if view > kMaxRangeV {
		view = kMaxRangeV
	}

	start, end := getPosLimit(aoi.transPos(pos), aoi.max, view)

	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			if objs := aoi.towers[i][j].GetObjects(); len(objs) > 0 {
				for _, id2 := range objs {
					if id2 == id {
						continue
					}

					aoi.listener.OnObjectAppear([]uint32{id}, id2)
				}
			}

			aoi.towers[i][j].AddWatcher(id)
		}
	}
}

func (aoi *TowerAOI) RemoveWatcher(id uint32, pos Point, view int32) {
	if view < 0 {
		return
	}

	if view > kMaxRangeV {
		view = kMaxRangeV
	}

	start, end := getPosLimit(aoi.transPos(pos), aoi.max, view)

	for i := start.X; i <= end.X; i++ {
		for j := start.Y; j <= end.Y; j++ {
			aoi.towers[i][j].RemoveWatcher(id)

			if objs := aoi.towers[i][j].GetObjects(); len(objs) > 0 {
				for _, id2 := range objs {
					if id2 == id {
						continue
					}

					aoi.listener.OnObjectDisappear([]uint32{id}, id2)
				}
			}
		}
	}
}

func (aoi *TowerAOI) UpdateWatcher(id uint32, oldPos, newPos Point, oldRange, newRange int32) bool {
	if !aoi.checkPos(oldPos) || !aoi.checkPos(newPos) {
		return false
	}

	pt1 := aoi.transPos(oldPos)
	pt2 := aoi.transPos(newPos)

	if oldRange < 0 || newRange < 0 {
		return false
	}

	if oldRange > kMaxRangeV {
		oldRange = kMaxRangeV
	}

	if newRange > kMaxRangeV {
		newRange = kMaxRangeV
	}

	if pt1.X == pt2.X && pt1.Y == pt2.Y && oldRange == newRange {
		return true
	}

	addTowers, removeTowers, _ := getChangedTowers(pt1, pt2, oldRange, newRange, aoi.towers, aoi.max)
	if len(removeTowers) > 0 {
		for _, v := range removeTowers {
			v.RemoveWatcher(id)

			if objs := v.GetObjects(); len(objs) > 0 {
				for _, id2 := range objs {
					if id2 == id {
						continue
					}

					aoi.listener.OnObjectDisappear([]uint32{id}, id2)
				}
			}
		}
	}

	if len(addTowers) > 0 {
		for _, v := range addTowers {
			v.AddWatcher(id)

			if objs := v.GetObjects(); len(objs) > 0 {
				for _, id2 := range objs {
					if id2 == id {
						continue
					}

					aoi.listener.OnObjectAppear([]uint32{id}, id2)
				}
			}
		}
	}

	return true
}

// Get changed towers for given pos
func getChangedTowers(pt1, pt2 Point, rv1, rv2 int32, towers [][]*Tower, max *Point) ([]*Tower, []*Tower, []*Tower) {
	start1, end1 := getPosLimit(pt1, max, rv1)
	start2, end2 := getPosLimit(pt2, max, rv2)

	addTowers, removeTowers, unchangeTowers := []*Tower{}, []*Tower{}, []*Tower{}

	for x := start1.X; x <= end1.X; x++ {
		for y := start1.Y; y <= end1.Y; y++ {
			if isInRect(x, y, start2, end2) {
				unchangeTowers = append(unchangeTowers, towers[x][y])
			} else {
				removeTowers = append(removeTowers, towers[x][y])
			}
		}
	}

	for x := start2.X; x <= end2.X; x++ {
		for y := start2.Y; y <= end2.Y; y++ {
			if !isInRect(x, y, start1, end1) {
				addTowers = append(addTowers, towers[x][y])
			}
		}
	}

	return addTowers, removeTowers, unchangeTowers
}

// Get the postion limit of gived range
func getPosLimit(pos Point, max *Point, rv int32) (start, end Point) {
	if (pos.X - rv) < 0 {
		start.X = 0
		end.X = 2 * rv
	} else if (pos.X + rv) > max.X {
		end.X = max.X
		start.X = max.X - 2*rv
	} else {
		start.X = pos.X - rv
		end.X = pos.X + rv
	}

	if (pos.Y - rv) < 0 {
		start.Y = 0
		end.Y = 2 * rv
	} else if (pos.Y + rv) > max.Y {
		end.Y = max.Y
		start.Y = max.Y - 2*rv
	} else {
		start.Y = pos.Y - rv
		end.Y = pos.Y + rv
	}

	if start.X < 0 {
		start.X = 0
	}

	if end.X > max.X {
		end.X = max.X
	}

	if start.Y < 0 {
		start.Y = 0
	}

	if end.Y > max.Y {
		end.Y = max.Y
	}

	return
}

func isInRect(x, y int32, start, end Point) bool {
	return x >= start.X && x <= end.X && y >= start.Y && y <= end.Y
}
