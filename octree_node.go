package octree

import (
	"github.com/The-Tensox/protometry"
)

// FIXME
var (
	CAPACITY = 5
)

// OctreeNode ...
type OctreeNode struct {
	objects  []Object
	region   protometry.Box
	children *[8]OctreeNode
}

// Insert ...
func (o *OctreeNode) insert(object Object) bool {
	// Object Bounds doesn't fit in node region => return false
	if !object.Bounds.Fit(o.region) {
		return false
	}

	// Number of objects < CAPACITY and children is nil => add in objects
	if len(o.objects) < CAPACITY && o.children == nil {
		o.objects = append(o.objects, object)
		return true
	}

	// Number of objects >= CAPACITY and children is nil => create children,
	// try to move all objects in children
	// and try to add in children otherwise add in objects
	if len(o.objects) >= CAPACITY && o.children == nil {
		o.split()

		objects := o.objects
		o.objects = []Object{}

		// Move old objects to children
		for i := range objects {
			o.insert(objects[i])
		}

	}

	// Children isn't nil => try to add in children otherwise add in objects
	for i := range o.children {
		if o.children[i].insert(object) {
			return true
		}
	}
	o.objects = append(o.objects, object)
	return true
}

func (o *OctreeNode) remove(object Object) bool {
	removedObject := false

	// Object outside Bounds
	if !object.Bounds.Fit(o.region) {
		return false
	}

	for i := 0; i < len(o.objects); i++ {
		// Found it ? delete it and return a copy
		if o.objects[i].Equal(object) {
			// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
			o.objects = append(o.objects[:i], o.objects[i+1:]...)
			return true
		}
	}

	if o.children != nil {
		for i := range o.children {
			if o.children[i].remove(object) {
				removedObject = true
				break
			}
		}
	}

	// Successfully removed in children
	if removedObject {
		// Try to merge nodes now that we've removed an item
		o.merge()
	}
	return removedObject
}

func (o *OctreeNode) move(object *Object, newBounds ...float64) bool {
	// Incorrect dimensions
	if (len(newBounds) != 3 && len(newBounds) != 6) || !o.remove(*object) {
		return false
	}
	if len(newBounds) == 3 {
		object.Bounds = *protometry.NewBoxOfSize(*protometry.NewVectorN(newBounds...), object.Bounds.Extents.Get(0)*2)
	} else { // Dimensions = 6
		object.Bounds = *protometry.NewBoxMinMax(newBounds...)
	}
	return o.insert(*object)
}

// Splits the OctreeNode into eight children.
func (o *OctreeNode) split() {
	subBoxes := o.region.Split()
	o.children = &[8]OctreeNode{}
	for i := range subBoxes {
		o.children[i] = OctreeNode{region: *subBoxes[i]}
	}
}

func (o *OctreeNode) getHeight() int {
	if o.children == nil {
		return 1
	}
	max := 0
	for _, c := range o.children {
		h := c.getHeight()
		if h > max {
			max = h
		}
	}
	return max + 1
}

func (o *OctreeNode) getNumberOfNodes() int {
	if o.children == nil {
		return 1
	}
	sum := len(o.children)
	for _, c := range o.children {
		n := c.getNumberOfNodes()
		sum += n
	}
	return sum
}

func (o *OctreeNode) getNumberOfObjects() int {
	if o.children == nil {
		return len(o.objects)
	}
	sum := len(o.objects)
	for _, c := range o.children {
		n := c.getNumberOfObjects()
		sum += n
	}
	return sum
}

func (o *OctreeNode) getColliding(bounds protometry.Box) []Object {
	// If current node region entirely fit inside desired Bounds,
	// No need to search somewhere else => return all objects
	if o.region.Fit(bounds) {
		return o.getAllObjects()
	}
	var objects []Object
	// If bounds doesn't intersects with region, no collision here => return empty
	if !o.region.Intersects(bounds) {
		return objects
	}
	// return objects that intersects with bounds and its children's objects
	for _, obj := range o.objects {
		if obj.Bounds.Intersects(bounds) {
			objects = append(objects, obj)
		}
	}
	// No children ? Stop here
	if o.children == nil {
		return objects
	}
	// Get the colliding children
	for _, c := range o.children {
		objects = append(objects, c.getColliding(bounds)...)
	}
	return objects
}

func (o *OctreeNode) getAllObjects() []Object {
	var objects []Object
	if o.children == nil {
		return o.objects
	}
	for _, c := range o.children {
		objects = append(objects, c.getAllObjects()...)
	}
	return objects
}

/* Merge all children into this node - the opposite of Split.
 * Note: We only have to check one level down since a merge will never happen if the children already have children,
 * since THAT won't happen unless there are already too many objects to merge.
 */
func (o *OctreeNode) merge() bool {
	totalObjects := len(o.objects)
	if o.children != nil {
		for _, child := range o.children {
			if child.children != nil {
				// If any of the *children* have children, there are definitely too many to merge,
				// or the child woudl have been merged already
				return false
			}
			totalObjects += len(child.objects)
		}
	}
	if totalObjects > CAPACITY {
		return false
	}

	// Note: We know children != null or we wouldn't be merging
	for i := range o.children {
		curChild := o.children[i]
		numObjects := len(curChild.objects)
		for j := numObjects - 1; j >= 0; j-- {
			curObj := curChild.objects[j]
			o.objects = append(o.objects, curObj)
		}
	}
	// Remove the child nodes (and the objects in them - they've been added elsewhere now)
	o.children = nil
	return true
}
