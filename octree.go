package octree

import (
	"github.com/The-Tensox/protometry"
)

// Octree ...
type Octree struct {
	root *OctreeNode
}

// NewOctree is a Octree constructor for ease of use
func NewOctree(region *protometry.Box) *Octree {
	return &Octree{
		root: &OctreeNode{region: *region},
	}
}

// Insert a point in the Octree, TODO: bool or point return?
func (o *Octree) Insert(point Point) bool {
	return o.root.insert(point)
}

// Get point(s) using their center, not their collider
func (o *Octree) Get(dims ...float64) *[]Point {
	if len(dims) == 3 {
		if p := o.root.get(dims...); p != nil {
			return &[]Point{*p}
		}
	} else if len(dims) == 6 {
		return o.root.getMultiple(dims...)
	}
	return nil
}

// Remove points at position
func (o *Octree) Remove(dims ...float64) *Point {
	return o.root.remove(dims...)
}

// Move point to a new position
func (o *Octree) Move(point Point, newPosition ...float64) *Point {
	return o.root.move(point, newPosition...)
}

// Raycast get all points colliding with an area
func (o *Octree) Raycast(origin, direction protometry.VectorN, maxDistance float64) *[]Point {
	return o.root.raycast(origin, direction, maxDistance)
}

/*
func (o *Octree) ToString() string {
	return o.root.ToString()
}
*/
