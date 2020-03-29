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

// Insert a object in the Octree, TODO: bool or object return?
func (o *Octree) Insert(object Object) bool {
	return o.root.insert(object)
}

// GetHeight debug function
func (o *Octree) GetHeight() int {
	return o.root.getHeight()
}

// GetNumberOfNodes debug function
func (o *Octree) GetNumberOfNodes() int {
	return o.root.getNumberOfNodes()
}

// GetNumberOfObjects debug function
func (o *Octree) GetNumberOfObjects() int {
	return o.root.getNumberOfObjects()
}

func (o *Octree) GetUsage() float64 {
	return float64(o.GetNumberOfObjects()) / float64(o.GetNumberOfNodes()*CAPACITY)
}

// // Get object(s) using their center, not their collider
// func (o *Octree) Get(dims ...float64) *[]Object {
// 	return o.root.get(dims...)
// }

// // Remove objects at position
// func (o *Octree) Remove(dims ...float64) *Object {
// 	return o.root.remove(dims...)
// }

// // Move object to a new position
// func (o *Octree) Move(object Object, newPosition ...float64) *Object {
// 	return o.root.move(object, newPosition...)
// }

// // Raycast get all objects colliding with an area
// func (o *Octree) Raycast(origin, direction protometry.VectorN, maxDistance float64) *[]Object {
// 	return o.root.raycast(origin, direction, maxDistance)
// }

/*
func (o *Octree) ToString() string {
	return o.root.ToString()
}
*/
