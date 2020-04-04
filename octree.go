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

// GetColliding returns an array of objects that intersect with the specified bounds, if any.
// Otherwise returns an empty array.
func (o *Octree) GetColliding(bounds protometry.Box) []Object {
	return o.root.getColliding(bounds)
}

// Remove object
func (o *Octree) Remove(object Object) bool {
	return o.root.remove(object)
}

// Move object to a new Bounds, pass a pointer because we want to modify the passed object data
func (o *Octree) Move(object *Object, newBounds ...float64) bool {
	return o.root.move(object, newBounds...)
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

// GetUsage ...
func (o *Octree) GetUsage() float64 {
	return float64(o.GetNumberOfObjects()) / float64(o.GetNumberOfNodes()*CAPACITY)
}
