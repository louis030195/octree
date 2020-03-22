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
	if err != nil {
		return false
	}
	if !in {
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

func (o *OctreeNode) getNewRegion(branch int) *OctreeNode {
	minD := o.region.GetMin().Dimensions
	maxD := o.region.GetMax().Dimensions
	var newNode *OctreeNode
	switch branch {
	case TLF:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				minD[0],
				minD[1],
				minD[2]),
			*protometry.NewVectorN(
				maxD[0]/2,
				maxD[1]/2,
				maxD[2]/2)))
	case TRF:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				maxD[0]/2,
				minD[1],
				minD[2]),
			*protometry.NewVectorN(
				maxD[0],
				maxD[1]/2,
				maxD[2]/2)))
	case BRF:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				maxD[0]/2,
				maxD[1]/2,
				minD[2]),
			*protometry.NewVectorN(
				maxD[0],
				maxD[1],
				maxD[2]/2)))
	case BLF:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				minD[0],
				maxD[1]/2,
				minD[2]),
			*protometry.NewVectorN(
				maxD[0]/2,
				maxD[1],
				maxD[2]/2)))
	case TLB:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				minD[0],
				minD[1],
				maxD[2]/2),
			*protometry.NewVectorN(
				maxD[0]/2,
				maxD[1]/2,
				maxD[2])))
	case TRB:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				maxD[0]/2,
				minD[1],
				maxD[2]/2),
			*protometry.NewVectorN(
				maxD[0],
				maxD[1]/2,
				maxD[2])))
	case BRB:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				maxD[0]/2,
				maxD[1]/2,
				maxD[2]/2),
			*protometry.NewVectorN(
				maxD[0],
				maxD[1],
				maxD[2])))
	case BLB:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				minD[0],
				maxD[1]/2,
				maxD[2]/2),
			*protometry.NewVectorN(
				maxD[0]/2,
				maxD[1],
				maxD[2])))
	}
	return newNode
}

func (o *OctreeNode) getOctant(position protometry.VectorN) int {
	oct := 0 // Not sure this func is correct
	center := o.region.GetCenter()
	if position.Dimensions[0] > center.Dimensions[0] {
		oct |= 4
	}
	if position.Dimensions[1] > center.Dimensions[1] {
		oct |= 2
	}
	if position.Dimensions[2] > center.Dimensions[2] {
		oct |= 1
	}
	return oct
}

// ToString Get a human readable representation of the state of
// this node and its contents.
func (o *OctreeNode) ToString() string {
	return o.recursiveToString("", "  ")
}

func (o *OctreeNode) recursiveToString(curIndent, stepIndent string) string {
	singleIndent := curIndent + stepIndent

	// default values
	childStr := "nil"
	centerStr := "nil"
	dataStr := "nil"

	if o.children != nil {
		doubleIndent := singleIndent + stepIndent

		// accumulate child strings
		childStr = ""
		for i, child := range o.children {
			childStr = childStr + fmt.Sprintf("%v%d: %v,\n", doubleIndent, i, child.recursiveToString(doubleIndent, stepIndent))
		}

		childStr = fmt.Sprintf("[\n%v%v]", childStr, singleIndent)
	}

	if o.center != nil {
		centerStr = o.center.ToString()
	}

	if o.data != nil {
		// not stringifying elements since their type is unknown
		dataStr = fmt.Sprintf("[%d]", len(o.data))
	}

	return fmt.Sprintf("Node{\n%vcenter: %v,\n%vdata: %v,\n%vregion: %v,\n%vchildren: %v,%v\n%v}", singleIndent, centerStr, singleIndent, dataStr, singleIndent, o.region, singleIndent, childStr, singleIndent, curIndent)
}
*/
