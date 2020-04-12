package octree

import (
	"fmt"
	"github.com/The-Tensox/protometry"
)

// Octree ...
type Octree struct {
	root *Node
}

// NewOctree is a Octree constructor for ease of use
func NewOctree(region *protometry.Box) *Octree {
	return &Octree{
		root: &Node{region: *region},
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

func (o *Octree) GetNodes() []Node {
	return o.root.getNodes()
}

// getHeight debug function
func (o *Octree) getHeight() int {
	return o.root.getHeight()
}

// getNumberOfNodes debug function
func (o *Octree) getNumberOfNodes() int {
	return o.root.getNumberOfNodes()
}

// getNumberOfObjects debug function
func (o *Octree) getNumberOfObjects() int {
	return o.root.getNumberOfObjects()
}

// getUsage ...
func (o *Octree) getUsage() float64 {
	return float64(o.getNumberOfObjects()) / float64(o.getNumberOfNodes()*CAPACITY)
}

func (o *Octree) toString(verbose bool) string {
	return fmt.Sprintf( "Octree: {\n%v\n}", o.root.toString(verbose))
}