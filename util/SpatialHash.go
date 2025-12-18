package util

const SHIFT_A = 10
const SHIFT_B = 10

type AABB struct {
	X1, Y1, X2, Y2 float64
}

func (a *AABB) Contains(x, y float64) bool {
	return x >= a.X1 && x <= a.X2 && y >= a.Y1 && y <= a.Y2
}

func (a *AABB) Intersects(b *AABB) bool {
	return a.X1 <= b.X2 && a.X2 >= b.X1 && a.Y1 <= b.Y2 && a.Y2 >= b.Y1
}

func (a *AABB) GetCenter() *Vector2D {
	return &Vector2D{
		X: (a.X1 + a.X2) / 2,
		Y: (a.Y1 + a.Y2) / 2,
	}
}

type Collidable interface {
	GetAABB() *AABB
}

type SpatialHash[T Collidable] struct {
	grid map[int][]T
}

func NewSpatialHash[T Collidable]() *SpatialHash[T] {
	return &SpatialHash[T]{
		grid: make(map[int][]T),
	}
}

func (sh *SpatialHash[T]) Clear() {
	sh.grid = make(map[int][]T)
}

func (sh *SpatialHash[T]) Insert(item T) {
	aabb := item.GetAABB()
	x1 := int(aabb.X1) >> SHIFT_A
	y1 := int(aabb.Y1) >> SHIFT_A
	x2 := int(aabb.X2) >> SHIFT_A
	y2 := int(aabb.Y2) >> SHIFT_A

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			key := x | (y << SHIFT_B)
			sh.grid[key] = append(sh.grid[key], item)
		}
	}
}

func (sh *SpatialHash[T]) Retrieve(aabb *AABB) (results []T) {
	x1 := int(aabb.X1) >> SHIFT_A
	y1 := int(aabb.Y1) >> SHIFT_A
	x2 := int(aabb.X2) >> SHIFT_A
	y2 := int(aabb.Y2) >> SHIFT_A

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			key := x | (y << SHIFT_B)
			if items, found := sh.grid[key]; found {
				for _, item := range items {
					if item.GetAABB().Intersects(aabb) {
						results = append(results, item)
					}
				}
			}
		}
	}

	return
}

func (sh *SpatialHash[T]) RetrieveAround(x, y, radius float64) []T {
	return sh.Retrieve(&AABB{
		X1: x - radius,
		Y1: y - radius,
		X2: x + radius,
		Y2: y + radius,
	})
}

func (sh *SpatialHash[T]) All() (results []T) {
	for _, items := range sh.grid {
		results = append(results, items...)
	}

	return
}
