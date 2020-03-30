package octree

import (
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

// OctreeNode ...
type OctreeNode struct {
	objects  []Object
	region   protometry.Box
	children *[8]OctreeNode
}

// Insert ...
// First case: object bounds doesn't fit in node region => return false
// Second case: number of objects < CAPACITY and children is nil => add in objects
// Third case: number of objects >= CAPACITY and children is nil => create children, try to move all objects in children
// and try to add in children otherwise add in objects
// Fourth case: children isn't nil => try to add in children otherwise add in objects
func (o *OctreeNode) insert(object Object) bool {
	// First case
	if !object.bounds.Fit(o.region) {
		return false
	}

	// Second case
	if len(o.objects) < CAPACITY && o.children == nil {
		o.objects = append(o.objects, object)
		return true
	}

	// Third case
	if len(o.objects) >= CAPACITY && o.children == nil {
		o.split()

		var objects []Object
		copy(objects, o.objects)
		o.objects = []Object{}

		// Move old objects to children
		for i := range objects {
			o.insert(objects[i])
		}
	}

	// Fourth case
	for i := range o.children {
		if o.children[i].insert(object) {
			return true
		}
	}
	o.objects = append(o.objects, object)
	return true
}

// TODO: test, probably incorrect, impl for getmultiple, maybe return object
func (o *OctreeNode) remove(object Object) *Object {
	var removedObject *Object

	// Object outside bounds
	if object.bounds.Fit(o.region) {
		return removedObject
	}

	for i := 0; i < len(o.objects); i++ {
		// Found it ? delete it and return a copy
		if o.objects[i].Equal(object) {
			removedObject = &o.objects[i]
			// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
			o.objects = append(o.objects[:i], o.objects[i+1:]...)
			return removedObject
		}
	}

	if o.children != nil {
		for i := range o.children {
			if removedObject = o.children[i].remove(object); removedObject != nil {
				break
			}
		}
	}

	// Successfully removed in children
	if removedObject != nil && o.children != nil {
		// Try to merge nodes now that we've removed an item
		// o.merge()
	}
	return removedObject
}

// Splits the OctreeNode into eight children.
func (o *OctreeNode) split() {
	subBoxes := o.region.MakeSubBoxes()
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
	// If current node region entirely fit inside desired bounds,
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
		if obj.bounds.Intersects(bounds) {
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

// /* Merge all children into this node - the opposite of Split.
//  * Note: We only have to check one level down since a merge will never happen if the children already have children,
//  * since THAT won't happen unless there are already too many objects to merge.
//  * TODO: to be tested
//  */
// func (o *OctreeNode) merge() {
// 	// Note: We know children != null or we wouldn't be merging
// 	for i := 0; i < 8; i++ {
// 		curChild := o.children[i]
// 		numObjects := len(curChild.objects)
// 		for j := numObjects - 1; j >= 0; j-- {
// 			curObj := curChild.objects[j]
// 			o.objects = append(o.objects, curObj)
// 		}
// 	}
// 	// Remove the child nodes (and the objects in them - they've been added elsewhere now)
// 	o.children = nil
// }

// /*
//  * We can shrink the octree if:
//  * - This node is >= double minLength in length
//  * - All objects in the root node are within one octant
//  * - This node doesn't have children, or does but 7/8 children are empty
//  * We can also shrink it if there are no objects left at all!
//  * TODO: to be tested
//  */
// func (o *OctreeNode) shrink() OctreeNode {
// 	return OctreeNode{}
// }

// // get return Object(s) based on their position only
// func (o *OctreeNode) get(dims ...float64) *[]Object {
// 	// Prepare an array of results
// 	var objects []Object
// 	var region *protometry.Box
// 	var ok bool
// 	var err error

// 	// We're looking for Object(s) at a precise position
// 	if len(dims) == 3 {
// 		region = protometry.NewBox(dims[0], dims[1], dims[2], dims[0], dims[1], dims[2]) // Ugly
// 	} else if len(dims) == 6 { // We're looking for Object(s) inside a region
// 		region = protometry.NewBox(dims[0], dims[1], dims[2], dims[3], dims[4], dims[5]) // Ugly
// 	}
// 	// Automatically abort if the range does not intersect this
// 	ok, err = o.region.Intersects(*region)
// 	if err != nil || !ok {
// 		return &objects
// 	}

// 	// Check objects at this level
// 	for i := range o.objects {
// 		ok, err = o.objects[i].position.In(*region)
// 		if err == nil && ok {
// 			objects = append(objects, o.objects[i])
// 		}
// 	}

// 	// Terminate here, if there are no children
// 	if o.children == nil {
// 		return &objects
// 	}

// 	// Otherwise, add the objects from the children
// 	for i := range o.children {
// 		if n := o.children[i].get(dims...); n != nil {
// 			objects = append(objects, *n...)
// 		}
// 	}

// 	return &objects
// }

// func (o *OctreeNode) move(object Object, newPosition ...float64) *Object {
// 	if len(newPosition) != 3 {
// 		return nil
// 	}
// 	// FIXME
// 	n := o.remove(object.position.Dimensions...)
// 	if n == nil {
// 		return n
// 	}
// 	newObject := NewObject(object.data, newPosition...)
// 	if res := o.insert(*newObject); res {
// 		return newObject
// 	}
// 	return nil
// }

// func (o *OctreeNode) raycast(origin, direction protometry.VectorN, maxDistance float64) *[]Object {
// 	// Prepare an array of results
// 	var objects []Object
// 	var destination *protometry.VectorN = protometry.NewVectorN(maxDistance, maxDistance, maxDistance)
// 	if maxDistance != math.MaxFloat64 {
// 		destination = direction.Mul(maxDistance)
// 	}
// 	destination = destination.Add(origin)
// 	ray := *protometry.NewBox(origin.Get(0), origin.Get(1), origin.Get(2), destination.Get(0), destination.Get(1), destination.Get(2))

// 	// Check objects at this level
// 	for i := range o.objects {
// 		in, err := o.objects[i].collider.bounds.Intersects(ray)
// 		if err == nil && in {
// 			objects = append(objects, o.objects[i])
// 		}
// 	}

// 	// Terminate here, if there are no children
// 	if o.children == nil {
// 		return &objects
// 	}

// 	// Otherwise, add the objects from the children
// 	for i := range o.children {
// 		// TODO: we can just move origin now
// 		if n := o.children[i].raycast(origin, direction, maxDistance); n != nil {
// 			objects = append(objects, *n...)
// 		}
// 	}

// 	return &objects
// }

// /*
// // ToString Get a human readable representation of the state of
// // this node and its contents.
// func (o *OctreeNode) ToString() string {
// 	return o.recursiveToString("", "  ")
// }

// func (o *OctreeNode) recursiveToString(curIndent, stepIndent string) string {
// 	singleIndent := curIndent + stepIndent

// 	// default values
// 	childStr := "nil"
// 	pointsStr := "nil"

// 	if o.children != nil {
// 		doubleIndent := singleIndent + stepIndent

// 		// accumulate child strings
// 		childStr = ""
// 		for i, child := range o.children {
// 			childStr = childStr + fmt.Sprintf("%v%d: %v,\n", doubleIndent, i, child.recursiveToString(doubleIndent, stepIndent))
// 		}

// 		childStr = fmt.Sprintf("[\n%v%v]", childStr, singleIndent)
// 	}

// 	for _, object := range o.objects {
// 		pointsStr = pointsStr + fmt.Sprintf("%v%v", singleIndent+stepIndent, object)
// 	}

// 	return fmt.Sprintf("Node{\n%vregion: %v,\n%vpoints: %v,\n%vchildren: %v,%v\n%v}", singleIndent, o.region, singleIndent, pointsStr, singleIndent, childStr, singleIndent, curIndent)
// }
// */

// /*
//  * Return the best fit this position should be placed among children
//  */
// func (o *OctreeNode) bestFit(position protometry.VectorN) int {
// 	oct := 0
// 	center := o.region.GetCenter()
// 	if position.Get(0) <= center.Get(0) {
// 		oct |= 4
// 	}
// 	if position.Get(1) <= center.Get(1) {
// 		oct |= 2
// 	}
// 	if position.Get(2) <= center.Get(2) {
// 		oct |= 1
// 	}
// 	return oct
// }
