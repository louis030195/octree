package octree

import (
	"math"

	"github.com/The-Tensox/protometry"
)

// FIXME
const (
	TLF = iota // top left front
	TRF        // top right front
	BRF        // bottom right front
	BLF        // bottom left front
	TLB        // top left back
	TRB        // top right back
	BRB        // bottom right back
	BLB        // bottom left back
)

// FIXME
var (
	CAPACITY = 5
)

// Point stores position, data and collider about the object
type Point struct {
	data     interface{}
	collider protometry.Box
	// position refers to the center of the Point
	position protometry.VectorN
}

// NewPoint is a Point constructor for ease of use
func NewPoint(data interface{}, dims ...float64) *Point {
	if len(dims) != 3 {
		return nil
	}
	return &Point{data: data, position: *protometry.NewVectorN(dims...)}
}

// NewPointCollide is a Point constructor with collider for ease of use
func NewPointCollide(data interface{}, collider protometry.Box, position protometry.VectorN) *Point {
	if len(position.Dimensions) != 3 {
		return nil
	}
	return &Point{data: data, collider: collider, position: position}
}

// OctreeNode ...
type OctreeNode struct {
	points   []Point
	region   protometry.Box
	children *[8]OctreeNode
}

// Insert ...
// First case: position isn't in region => return error
// Second case: number of points < CAPACITY and children is nil => add in points
// Third case: number of points >= CAPACITY and children is nil => create children and add in children
// Fourth case: children isn't nil => add in children
func (o *OctreeNode) insert(point Point) bool {
	// First case
	in, err := point.position.In(o.region)
	if err != nil || !in {
		return false
	}

	// Second case
	if len(o.points) < CAPACITY && o.children == nil {
		o.points = append(o.points, point)
		return true
	}

	// Third case
	if len(o.points) >= CAPACITY && o.children == nil {
		subBoxes := o.region.MakeSubBoxes() // Can fail
		o.children = &[8]OctreeNode{}
		for i := range subBoxes {
			o.children[i] = OctreeNode{region: *subBoxes[i]}
		}
		// Move old points to children
		for i := range o.points {
			o.insert(o.points[i])
		}
		// Empty it
		o.points = []Point{}
	}

	// Fourth case
	for i := range o.children {
		if o.children[i].insert(point) {
			return true
		}
	}
	return false
}

// getMultiple return a list of points which have they center found inside the defined region TODO: FIXME COMPLEXITY
func (o *OctreeNode) getMultiple(dims ...float64) *[]Point {
	// Prepare an array of results
	var points []Point
	region := *protometry.NewBox(dims[0], dims[1], dims[2], dims[3], dims[4], dims[5]) // Ugly

	// Automatically abort if the range does not intersect this
	intersect, err := o.region.Intersects(region)
	if err != nil || !intersect {
		return nil // Empty
	}

	// Check objects at this level
	for i := range o.points {
		in, err := o.points[i].position.In(region)
		if err == nil && in {
			points = append(points, o.points[i])
		}
	}

	// Terminate here, if there are no children
	if o.children == nil {
		return &points
	}

	// Otherwise, add the points from the children
	for i := range o.children {
		if n := o.children[i].getMultiple(dims...); n != nil {
			points = append(points, *n...)
		}
	}

	return &points
}

// get return a Point, TODO: FIXME COMPLEXITY
func (o *OctreeNode) get(dims ...float64) *Point {
	position := *protometry.NewVectorN(dims...)
	in, err := position.In(o.region)
	if err != nil || !in {
		return nil
	}
	for i := range o.points {
		eq, err := o.points[i].position.ApproxEqual(position)
		if err != nil {
			return nil
		}
		if eq {
			return &o.points[i]
		}
	}
	if o.children != nil {
		for i := range o.children {
			if n := o.children[i].get(dims...); n != nil {
				return n
			}
		}
	}
	return nil
}

// TODO: test, probably incorrect, impl for getmultiple, maybe return point
func (o *OctreeNode) remove(dims ...float64) *Point {
	if len(dims) != 3 {
		return nil
	}
	// FIXME
	return o.get(dims...)
}

func (o *OctreeNode) move(point Point, newPosition ...float64) *Point {
	if len(newPosition) != 3 {
		return nil
	}
	// FIXME
	n := o.remove(point.position.Dimensions...)
	if n == nil {
		return n
	}
	newPoint := NewPoint(point.data, newPosition...)
	if res := o.insert(*newPoint); res {
		return newPoint
	}
	return nil
}

func (o *OctreeNode) raycast(origin, direction protometry.VectorN, maxDistance float64) *[]Point {
	// Prepare an array of results
	var points []Point
	var destination *protometry.VectorN = protometry.NewVectorN(maxDistance, maxDistance, maxDistance)
	if maxDistance != math.MaxFloat64 {
		destination = direction.Mul(maxDistance)
	}
	destination = destination.Add(origin)
	ray := *protometry.NewBox(origin.Get(0), origin.Get(1), origin.Get(2), destination.Get(0), destination.Get(1), destination.Get(2))

	// Check objects at this level
	for i := range o.points {
		in, err := o.points[i].collider.Intersects(ray)
		if err == nil && in {
			points = append(points, o.points[i])
		}
	}

	// Terminate here, if there are no children
	if o.children == nil {
		return &points
	}

	// Otherwise, add the points from the children
	for i := range o.children {
		// TODO: we can just move origin now
		if n := o.children[i].raycast(origin, direction, maxDistance); n != nil {
			points = append(points, *n...)
		}
	}

	return &points
}

/*
// ToString Get a human readable representation of the state of
// this node and its contents.
func (o *OctreeNode) ToString() string {
	return o.recursiveToString("", "  ")
}

func (o *OctreeNode) recursiveToString(curIndent, stepIndent string) string {
	singleIndent := curIndent + stepIndent

	// default values
	childStr := "nil"
	pointsStr := "nil"

	if o.children != nil {
		doubleIndent := singleIndent + stepIndent

		// accumulate child strings
		childStr = ""
		for i, child := range o.children {
			childStr = childStr + fmt.Sprintf("%v%d: %v,\n", doubleIndent, i, child.recursiveToString(doubleIndent, stepIndent))
		}

		childStr = fmt.Sprintf("[\n%v%v]", childStr, singleIndent)
	}

	for _, point := range o.points {
		pointsStr = pointsStr + fmt.Sprintf("%v%v", singleIndent+stepIndent, point)
	}

	return fmt.Sprintf("Node{\n%vregion: %v,\n%vpoints: %v,\n%vchildren: %v,%v\n%v}", singleIndent, o.region, singleIndent, pointsStr, singleIndent, childStr, singleIndent, curIndent)
}
*/
