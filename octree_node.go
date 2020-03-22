package octree

import (
	"github.com/The-Tensox/protometry"
)

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

var (
	CAPACITY = 5
)

type Point struct {
	data     interface{}
	position protometry.VectorN
}

func NewPoint(x, y, z float64, data interface{}) Point {
	return Point{data: data, position: *protometry.NewVectorN(x, y, z)}
}

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
func (o *OctreeNode) Insert(point Point) bool {
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
			o.Insert(o.points[i])
		}
		// Empty it
		o.points = []Point{}
	}

	// Fourth case
	for i := range o.children {
		if o.children[i].Insert(point) {
			return true
		}
	}
	return false
}

// Range find all points that appear within a region
func (o *OctreeNode) Range(region protometry.Box) []Point {
	// Prepare an array of results
	var points []Point

	// Automatically abort if the range does not intersect this
	intersect, err := o.region.Intersects(region)
	if err != nil || !intersect {
		return points // Empty
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
		return points
	}

	// Otherwise, add the points from the children
	for i := range o.children {
		points = append(points, o.children[i].Range(region)...)
	}

	return points
}

/*
func (o *OctreeNode) Search(position protometry.VectorN) (*OctreeNode, error) {
	pos := o.getOctant(position)
	if o.children[pos].center == nil {
		// Region node
		return o.children[pos].Search(position)
	} else if o.children[pos].center.Dimensions[0] == math.MaxFloat64 {
		// Empty node
		return nil, ErrtreeFailedToFindNode
	}
	eq, err := o.children[pos].center.ApproxEqual(position)
	if err != nil {
		return nil, err
	}
	if eq {
		return o.children[pos], nil
	}
	return nil, ErrtreeFailedToFindNode
}

func (o *OctreeNode) Remove(position protometry.VectorN) error {
	n, err := o.Search(position)
	if err != nil {
		return err
	}
	if n != nil {
		*n = *NewEmptyOctreeNode()
	}
	return nil
}
*/

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
