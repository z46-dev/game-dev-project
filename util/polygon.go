package util

import "math"

type PolygonXYPoint struct {
	x, y float64
}

type Polygon struct {
	numPoints              int
	Reference, Points      []*Vector2D
	x, y, radius, rotation float64
	AABB                   *AABB
}

func NewPolygon(points []*Vector2D, position *Vector2D, radius, rotation float64) (p *Polygon) {
	p = &Polygon{
		numPoints: len(points),
		Reference: make([]*Vector2D, len(points)),
		Points:    make([]*Vector2D, len(points)),
		x:         0,
		y:         0,
		radius:    0,
		rotation:  0,
		AABB:      &AABB{},
	}

	for i := range points {
		p.Reference[i] = Vector(points[i].X, points[i].Y)
		p.Points[i] = Vector(0, 0)
	}

	p.Transform(position, radius, rotation)
	return
}

func (p *Polygon) Transform(position *Vector2D, radius, rotation float64) {
	if p.x == position.X && p.y == position.Y && p.radius == radius && p.rotation == rotation {
		return
	}

	var cos, sin float64 = math.Cos(rotation), math.Sin(rotation)
	for i := range p.numPoints {
		p.Points[i].X = position.X + (p.Reference[i].X*cos-p.Reference[i].Y*sin)*radius
		p.Points[i].Y = position.Y + (p.Reference[i].X*sin+p.Reference[i].Y*cos)*radius
	}

	p.x, p.y, p.radius, p.rotation = position.X, position.Y, radius, rotation
	p.updateAABB()
}

func (p *Polygon) updateAABB() {
	var minX, minY, maxX, maxY float64 = math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)

	for i := range p.numPoints {
		minX = min(minX, p.Points[i].X)
		minY = min(minY, p.Points[i].Y)
		maxX = max(maxX, p.Points[i].X)
		maxY = max(maxY, p.Points[i].Y)
	}

	p.AABB.X1, p.AABB.Y1, p.AABB.X2, p.AABB.Y2 = minX, minY, maxX, maxY
}

func makePolygonFromPoints(points []*Vector2D) (p *Polygon) {
	p = &Polygon{
		numPoints: len(points),
		Reference: make([]*Vector2D, len(points)),
		Points:    make([]*Vector2D, len(points)),
		AABB:      &AABB{},
	}

	for i := range points {
		p.Reference[i] = Vector(points[i].X, points[i].Y)
		p.Points[i] = Vector(points[i].X, points[i].Y)
	}

	p.updateAABB()
	return
}

func polygonArea(points []*Vector2D) (area float64) {
	for i := range points {
		var j int = (i + 1) % len(points)
		area += points[i].X*points[j].Y - points[j].X*points[i].Y
	}

	return area * 0.5
}

func isConvexPolygon(points []*Vector2D) (convex bool) {
	if len(points) < 4 {
		return true
	}

	const eps float64 = 1e-9
	var sign float64
	for i := range points {
		var (
			a *Vector2D = points[i]
			b *Vector2D = points[(i+1)%len(points)]
			c *Vector2D = points[(i+2)%len(points)]
			cross       float64 = (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
		)

		if math.Abs(cross) <= eps {
			continue
		}

		if sign == 0 {
			sign = math.Copysign(1, cross)
			continue
		}

		if sign*cross < 0 {
			return false
		}
	}

	return true
}

func pointInTriangle(point, a, b, c *Vector2D, ccw bool) (inside bool) {
	const eps float64 = 1e-9
	var (
		ab float64 = (b.X-a.X)*(point.Y-a.Y) - (b.Y-a.Y)*(point.X-a.X)
		bc float64 = (c.X-b.X)*(point.Y-b.Y) - (c.Y-b.Y)*(point.X-b.X)
		ca float64 = (a.X-c.X)*(point.Y-c.Y) - (a.Y-c.Y)*(point.X-c.X)
	)

	if ccw {
		return ab >= -eps && bc >= -eps && ca >= -eps
	}

	return ab <= eps && bc <= eps && ca <= eps
}

func convexParts(points []*Vector2D) (parts []*Polygon) {
	if len(points) < 4 || isConvexPolygon(points) {
		return []*Polygon{makePolygonFromPoints(points)}
	}

	var (
		indices []int = make([]int, len(points))
		ccw     bool  = polygonArea(points) >= 0
	)

	for i := range points {
		indices[i] = i
	}

	const eps float64 = 1e-9
	var guard int
	for len(indices) > 3 && guard < len(points)*len(points) {
		guard++
		var earFound bool
		for i := range indices {
			var (
				prev int = indices[(i-1+len(indices))%len(indices)]
				curr int = indices[i]
				next int = indices[(i+1)%len(indices)]

				a *Vector2D = points[prev]
				b *Vector2D = points[curr]
				c *Vector2D = points[next]
				cross       float64 = (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
			)

			if math.Abs(cross) <= eps {
				continue
			}

			if ccw && cross <= 0 {
				continue
			}
			if !ccw && cross >= 0 {
				continue
			}

			var anyInside bool
			for _, idx := range indices {
				if idx == prev || idx == curr || idx == next {
					continue
				}
				if pointInTriangle(points[idx], a, b, c, ccw) {
					anyInside = true
					break
				}
			}

			if anyInside {
				continue
			}

			parts = append(parts, makePolygonFromPoints([]*Vector2D{a, b, c}))
			indices = append(indices[:i], indices[i+1:]...)
			earFound = true
			break
		}

		if !earFound {
			break
		}
	}

	if len(indices) == 3 {
		parts = append(parts, makePolygonFromPoints([]*Vector2D{
			points[indices[0]], points[indices[1]], points[indices[2]],
		}))
	}

	if len(parts) == 0 {
		return []*Polygon{makePolygonFromPoints(points)}
	}

	return parts
}

func (p *Polygon) PointIsInside(point *Vector2D) (inside bool) {
	var x1, y1 float64 = p.Points[p.numPoints-1].X, p.Points[p.numPoints-1].Y

	for i := range p.numPoints {
		var x2, y2 float64 = p.Points[i].X, p.Points[i].Y

		if (point.Y < y1) != (point.Y < y2) && (point.X < (x2-x1)*(point.Y-y1)/(y2-y1)+x1) {
			inside = !inside
		}

		x1, y1 = x2, y2
	}

	return
}

func (p *Polygon) CircleIntersectsEdge(p1, p2 *Vector2D, circlePoint *Vector2D, circleRadius float64) (intersects bool) {
	var (
		ABx, ABy float64 = p2.X - p1.X, p2.Y - p1.Y
		ACx, ACy float64 = circlePoint.X - p1.X, circlePoint.Y - p1.Y
		t        float64 = max(0, min(1, (ABx*ACx+ABy*ACy)/(ABx*ABx+ABy*ABy)))
		dx, dy   float64 = (p1.X + ABx*t) - circlePoint.X, (p1.Y + ABy*t) - circlePoint.Y
	)

	intersects = (dx*dx + dy*dy) < (circleRadius * circleRadius)
	return
}

func (p *Polygon) CircleIntersects(circlePoint *Vector2D, radius float64) (intersects bool) {
	if intersects = p.PointIsInside(circlePoint); intersects {
		return
	}

	for i := range p.numPoints {
		if intersects = p.CircleIntersectsEdge(p.Points[i], p.Points[(i+1)%p.numPoints], circlePoint, radius); intersects {
			return
		}
	}

	return
}

func (p *Polygon) GetClosestPointOnEdge(p1, p2, p3 *Vector2D) (point *Vector2D) {
	var (
		ABx, ABy float64 = p2.X - p1.X, p2.Y - p1.Y
		ACx, ACy float64 = p3.X - p1.X, p3.Y - p1.Y
		t        float64 = max(0, min(1, (ABx*ACx+ABy*ACy)/(ABx*ABx+ABy*ABy)))
	)

	point = Vector(p1.X+ABx*t, p1.Y+ABy*t)
	return
}

func (p *Polygon) GetAxes() (axes []*Vector2D) {
	axes = make([]*Vector2D, p.numPoints)

	for i := range p.numPoints {
		var (
			x1, y1 float64 = p.Points[i].X, p.Points[i].Y
			x2, y2 float64 = p.Points[(i+1)%p.numPoints].X, p.Points[(i+1)%p.numPoints].Y
			xE, yE float64 = x2 - x1, y2 - y1
		)

		axes[i] = Vector(-yE, xE).Normalize()
	}

	return
}

func (p *Polygon) ProjectOnto(axis *Vector2D) (minimum, maximum float64) {
	minimum, maximum = math.Inf(1), math.Inf(-1)

	for i := range p.numPoints {
		var dotProduct float64 = p.Points[i].Dot(axis)
		minimum, maximum = min(minimum, dotProduct), max(maximum, dotProduct)
	}

	return
}

func ResolveCirclePolygon(circlePoint *Vector2D, circleRadius float64, polygon *Polygon) (point *Vector2D, angle float64) {
	circleRadius += 1

	var (
		closestDistance float64   = math.Inf(1)
		closestPoint    *Vector2D = nil
	)

	for i := range polygon.numPoints {
		var (
			point *Vector2D = polygon.GetClosestPointOnEdge(polygon.Points[i], polygon.Points[(i+1)%polygon.numPoints], circlePoint)
			dist  float64   = Distance(point, circlePoint)
		)

		if dist < closestDistance {
			closestDistance = dist
			closestPoint = point
		}
	}

	angle = AngleBetween(closestPoint, circlePoint)
	var newPoint *Vector2D = Vector(circlePoint.X-circleRadius*math.Cos(angle), circlePoint.Y-circleRadius*math.Sin(angle))
	angle = AngleBetween(newPoint, closestPoint)
	if polygon.PointIsInside(newPoint) {
		angle += math.Pi
	}

	point = Vector(closestPoint.X+circleRadius*math.Cos(angle), closestPoint.Y+circleRadius*math.Sin(angle))
	return
}

func TwoPolygonsIntersect(p1, p2 *Polygon) (intersects bool) {
	var parts1 []*Polygon = convexParts(p1.Points)
	var parts2 []*Polygon = convexParts(p2.Points)
	for _, a := range parts1 {
		for _, b := range parts2 {
			if twoPolygonsIntersectConvex(a, b) {
				return true
			}
		}
	}

	return false
}

func twoPolygonsIntersectConvex(p1, p2 *Polygon) (intersects bool) {
	for i := range p1.numPoints {
		var (
			p1X1, p1Y1             float64   = p1.Points[i].X, p1.Points[i].Y
			p1X2, p1Y2             float64   = p1.Points[(i+1)%p1.numPoints].X, p1.Points[(i+1)%p1.numPoints].Y
			normal                 *Vector2D = Vector(-(p1X2 - p1X1), p1Y2-p1Y1).Normalize()
			min1, max1, min2, max2 float64
		)

		min1, max1 = p1.ProjectOnto(normal)
		min2, max2 = p2.ProjectOnto(normal)

		if max1 < min2 || max2 < min1 {
			intersects = false
			return
		}
	}

	for i := range p1.numPoints {
		if intersects = p2.PointIsInside(p1.Points[i]); intersects {
			return
		}
	}

	for i := range p2.numPoints {
		if intersects = p1.PointIsInside(p2.Points[i]); intersects {
			return
		}
	}

	return
}

func ResolveTwoPolygons(p1, p2 *Polygon) (resolution *Vector2D) {
	var (
		parts1    []*Polygon = convexParts(p1.Points)
		parts2    []*Polygon = convexParts(p2.Points)
		bestMTV   *Vector2D  = nil
		bestScore float64    = math.Inf(1)
	)

	for _, a := range parts1 {
		for _, b := range parts2 {
			var mtv *Vector2D = resolveTwoPolygonsConvex(a, b)
			if mtv == nil {
				continue
			}

			var score float64 = mtv.SquaredMagnitude()
			if score < bestScore {
				bestScore = score
				bestMTV = mtv
			}
		}
	}

	return bestMTV
}

func resolveTwoPolygonsConvex(p1, p2 *Polygon) (resolution *Vector2D) {
	var (
		mtv        *Vector2D = nil
		minOverlap float64   = math.Inf(1)
	)

	var axes []*Vector2D = append(p1.GetAxes(), p2.GetAxes()...)
	for _, axis := range axes {
		var (
			min1, max1 float64 = p1.ProjectOnto(axis)
			min2, max2 float64 = p2.ProjectOnto(axis)
			overlap    float64 = min(max1, max2) - max(min1, min2)
		)

		if overlap <= 0 {
			return nil
		}

		if overlap < minOverlap {
			minOverlap = overlap
			var flip float64 = 1
			if min1 < min2 {
				flip = -1
			}

			mtv = axis.Copy().Scale(overlap * flip)
		}
	}

	return mtv
}
