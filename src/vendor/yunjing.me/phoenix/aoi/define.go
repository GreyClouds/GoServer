package aoi

type Point struct {
	X int32
	Y int32
}

type Listener interface {
	OnObjectAppear([]uint32, uint32)
	OnObjectDisappear([]uint32, uint32)
}
