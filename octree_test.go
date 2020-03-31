package octree

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/The-Tensox/protometry"
)

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func boilerplateTree(t *testing.T) *Octree {
	o := NewOctree(protometry.NewBox(1, 1, 1, 4, 4, 4))
	ok := o.Insert(*NewObjectCube(0, 2, 2, 3, 0.5))
	equals(t, true, ok)
	return o
}

func TestOctree_NewOctree(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), 1))
	// Should be [(-1,-1,-1), (1, 1, 1)]
	equals(t, *protometry.NewBox(-1, -1, -1, 1, 1, 1), o.root.region)
}

func TestOctreeNode_Insert(t *testing.T) {
	o := boilerplateTree(t)
	ok := o.Insert(*NewObjectCube(5, 3, 3, 3, 1))
	equals(t, true, ok)
	ok = o.Insert(*NewObjectCube(6, 2, 2, 2, 1))
	equals(t, true, ok)
	equals(t, 3, len(o.root.objects))
	ok = o.Insert(*NewObjectCube(7, 2, 2, 2, 1))
	equals(t, true, ok)

	// Go over capacity threshold, force a split
	size := 10.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		ok = o.Insert(*NewObjectCube(0, i, i, i, 1))
		equals(t, true, ok)
	}

	// We inserted 10 objects so we should have 10 objects ;)
	equals(t, 10, o.GetNumberOfObjects())

	// Let's test with more scale
	size = 1000.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 1.; i < size; i++ {
		for j := 1.; j < size; j++ {
			ok = o.Insert(*NewObjectCube(0, i, j, i, 1))
			equals(t, true, ok)
			ok = o.Insert(*NewObjectCube(0, i, j, i, 1))
			equals(t, true, ok)
		}
	}
	equals(t, true, o.GetUsage() < 1)
	t.Logf("Octree height: %v", o.GetHeight())
	t.Logf("Octree usage: %v", o.GetUsage())
	t.Logf("Octree nodes: %v", o.GetNumberOfNodes())
	t.Logf("Octree objects: %v", o.GetNumberOfObjects())
}

func TestOctreeNode_GetColliding(t *testing.T) {
	o := NewOctree(protometry.NewBox(1, 1, 1, 4, 4, 4))
	ok := o.Insert(*NewObjectCube(0, 2, 2, 3, 0.5))
	equals(t, true, ok)
	ok = o.Insert(*NewObjectCube(5, 3, 3, 3, 1))
	equals(t, true, ok)
	ok = o.Insert(*NewObjectCube(6, 2, 2, 2, 1))
	equals(t, true, ok)
	equals(t, 3, len(o.root.objects))
	ok = o.Insert(*NewObjectCube(7, 2, 2, 2, 1))
	equals(t, true, ok)

	colliders := o.GetColliding(*protometry.NewBox(0, 0, 0, 0.9, 0.9, 0.9))
	equals(t, 0, len(colliders))
	colliders = o.GetColliding(*protometry.NewBox(0, 0, 0, 1, 1, 1))
	equals(t, 2, len(colliders))
	colliders = o.GetColliding(*protometry.NewBox(1, 1, 1, 1.1, 1.1, 1.1))
	equals(t, 2, len(colliders))
	equals(t, 6, colliders[0].data)
	equals(t, 7, colliders[1].data)
}

func TestOctree_Remove(t *testing.T) {
	o := boilerplateTree(t)
	myObj := NewObjectCube(27, 2, 2, 3, 0.5)
	ok := o.Insert(*myObj)
	equals(t, true, ok)
	removedObj := o.Remove(*myObj)
	equals(t, true, removedObj != nil)
	equals(t, 27, removedObj.data)
	oldBounds := myObj.bounds
	equals(t, true, removedObj.bounds.Equal(oldBounds))
	removedObj = o.Remove(*NewObjectCube(12, 2, 2, 3, 0.5))
	equals(t, true, removedObj == nil)

	// New octree
	size := 1000.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	var objects []Object
	for i := 1.; i < size; i++ {
		myObj = NewObjectCube(0, i, i, i, 1)
		ok = o.Insert(*myObj)
		equals(t, true, ok)
		objects = append(objects, *myObj)
	}
	for i := range objects {
		removedObj = o.Remove(objects[i])
		equals(t, true, removedObj != nil)
	}
	equals(t, 0, o.GetNumberOfObjects())
	equals(t, 1, o.GetNumberOfNodes()) // Only root
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	objects = []Object{}
	for i := 0.; i < 9; i++ {
		myObj = NewObjectCube(0, i, i, i, 1)
		objects = append(objects, *myObj)
		ok = o.Insert(*myObj)
		equals(t, true, ok)
	}
	equals(t, 9, o.GetNumberOfObjects())
	equals(t, 8, len(o.root.children))
	o.Remove(*myObj)
	var nilChildren *[8]OctreeNode
	// Shouldn't have merged
	equals(t, true, nilChildren != o.root.children)
	// One less object
	equals(t, 8, o.GetNumberOfObjects())
}

// func BenchmarkOctreeNode_Insert(b *testing.B) {
// 	b.StartTimer()
// 	// New octree
// 	size := float64(b.N)
// 	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
// 	for i := 0.; i < size; i++ {
// 		for j := 0.; j < size; j++ {
// 			o.Insert(*NewObject(0, i, j, i))
// 		}
// 	}
// 	b.StopTimer()
// }

// func TestOctreeNode_Get(t *testing.T) {
// 	o := boilerplateTree(t)
// 	objects := o.Get(1, 1, 1, 4, 4, 4)
// 	equals(t, 3, len(*objects))
// 	objects = o.Get(1, 1, 1, 3, 3, 3)
// 	equals(t, 1, len(*objects))
// 	objects = o.Get(1, 1, 1, 3, 3, 4)
// 	equals(t, 2, len(*objects))
// 	objects = o.Get(1, 1, 1, 3, 4, 4)
// 	equals(t, 3, len(*objects))

// 	o = boilerplateTree(t)
// 	object := *o.Get(3, 3, 3)
// 	equals(t, *protometry.NewVectorN(3, 3, 3), object[0].position)
// 	equals(t, 2, object[0].data)
// 	object = *o.Get(3, 3, 4)
// 	equals(t, *protometry.NewVectorN(3, 3, 4), object[0].position)
// 	equals(t, 3, object[0].data)
// 	object = *o.Get(3, 4, 4)
// 	equals(t, *protometry.NewVectorN(3, 4, 4), object[0].position)
// 	equals(t, 4, object[0].data)
// 	var nilObjectSlice []Object
// 	equals(t, &nilObjectSlice, o.Get(4, 4, 4))
// }

// func BenchmarkOctreeNode_GetMultiple(b *testing.B) {
// 	size := float64(b.N)
// 	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
// 	for i := 0.; i < size; i++ {
// 		for j := 0.; j < size; j++ {
// 			o.Insert(*NewObject(0, i, j, i))
// 		}
// 	}
// 	b.StartTimer()
// 	for i := 0.; i < size; i++ {
// 		for j := 0.; j < size; j++ {
// 			o.Get(i, j, i, i, j, i)
// 		}
// 	}
// 	b.StopTimer()
// }

// func BenchmarkOctreeNode_GetOne(b *testing.B) {
// 	size := float64(b.N)
// 	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
// 	for i := 0.; i < size; i++ {
// 		for j := 0.; j < size; j++ {
// 			o.Insert(*NewObject(0, i, j, i))
// 		}
// 	}
// 	b.StartTimer()
// 	for i := 0.; i < size; i++ {
// 		for j := 0.; j < size; j++ {
// 			o.Get(i, j, i)
// 		}
// 	}
// 	b.StopTimer()
// }

// func TestOctree_MoveWithoutCollision(t *testing.T) {
// 	o := boilerplateTree(t)
// 	// Let's get one object previously inserted in the tree
// 	resObjects := o.Get(3, 3, 3)
// 	equals(t, true, resObjects != nil)
// 	objects := *resObjects
// 	equals(t, 1, len(objects))
// 	p := objects[0]

// 	// Let's try to move it
// 	resObject := o.Move(p, 0, 0, 0)
// 	equals(t, true, resObject != nil)
// 	p = *resObject

// 	// It should be moved to the new position
// 	equals(t, protometry.NewVectorN(0, 0, 0), p.position)

// 	// And there shouldn't be anymore objects in the old position
// 	resObjects = o.Get(3, 3, 3)
// 	equals(t, true, resObjects != nil)
// 	equals(t, 0, len(*resObjects))
// }

// func TestOctree_MoveWithCollision(t *testing.T) {
// 	// TODO: implement function find name among:[AddForce, MoveTo, Push ...]
// 	// either on Octree or ?
// 	/*
// 		Should test movement taking into account physics:

// 			_ _ _ _
// 		   |_|_|_|_|
// 		   |A|_|_|B|
// 		   |_|_|_|_|
// 		   |_|_|_|_|

// 		   A.collider = square of side 1 => |A|
// 		   B.collider = square of side 1 => |B|

// 		If I try to move B to A position, pushing B to -1;0;0, no gravity
// 		It should end up like that

// 			_ _ _ _
// 		   |_|_|_|_|
// 		   |A|B|_|_|
// 		   |_|_|_|_|
// 		   |_|_|_|_|
// 	*/
// }

// func TestOctree_Raycast(t *testing.T) {
// 	o := NewOctree(protometry.NewBox(0, 0, 0, 10, 10, 10))
// 	p1 := NewCollider(*protometry.NewBox(0, 0, 0, 0, 1, 0)) // Line collider
// 	ok := o.Insert(*NewObjectCollide(1, p1, *p1.bounds.GetCenter()))
// 	equals(t, true, ok)
// 	p2 := NewCollider(*protometry.NewBox(0, 2, 0, 0, 3, 0))
// 	ok = o.Insert(*NewObjectCollide(2, p2, *p2.bounds.GetCenter()))
// 	equals(t, true, ok)
// 	p3 := NewCollider(*protometry.NewBox(0, 4, 0, 0, 5, 0))
// 	ok = o.Insert(*NewObjectCollide(3, p3, *p3.bounds.GetCenter()))
// 	equals(t, true, ok)

// 	// Cast a ray toward up from 0;0;0 of length 10
// 	objects := *o.Raycast(*protometry.NewVector3Zero(), *protometry.NewVectorN(0, 1, 0), 10)
// 	equals(t, 3, len(objects))

// 	// Cast a ray toward up from 0;0;0 of length 2.5
// 	objects = *o.Raycast(*protometry.NewVector3Zero(), *protometry.NewVectorN(0, 1, 0), 2.5)
// 	equals(t, 2, len(objects))

// 	// Cast a ray toward up from 0;2.1;0 of length 7
// 	objects = *o.Raycast(*protometry.NewVectorN(0, 2.1, 0), *protometry.NewVectorN(0, 1, 0), 7)
// 	equals(t, 2, len(objects))
// 	equals(t, p3, objects[len(objects)-1].collider)

// 	// Cast a ray toward up from 0;5.1;0 of length 4
// 	equals(t, 0, len(*o.Raycast(*protometry.NewVectorN(0, 5.1, 0), *protometry.NewVectorN(0, 1, 0), 4)))

// 	// New octree
// 	size := 1000.
// 	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
// 	for i := 0.; i < size; i++ {
// 		p1 := NewCollider(*protometry.NewBox(0, i, 0, 0, i+1, 0))
// 		o.Insert(*NewObjectCollide(0, p1, *p1.bounds.GetCenter()))
// 		o.Raycast(*protometry.NewVectorN(0, i, 0), *protometry.NewVectorN(0, 1, 0), 1)
// 	}
// 	equals(t, 1000, len(*o.Raycast(*protometry.NewVectorN(0, 0, 0), *protometry.NewVectorN(0, 1, 0), size)))

// 	// Edge cases
// 	// New octree
// 	size = 4.
// 	o = NewOctree(protometry.NewBox(0, 0, 0, 4, 4, 4))

// 	/*
// 			_ _ _ _
// 		   |_|_|_|_|
// 		   |A|_|_|B|
// 		   |_|_|_|_|
// 		   |_|_|_|_|

// 		   A.collider = square of side 1 => |A|
// 		   B.collider = rectangle of height 1 and width 4 => |_|_|_|B|

// 	*/
// 	A := NewCollider(*protometry.NewBox(0, 2, 0, 1, 3, 0))
// 	ok = o.Insert(*NewObjectCollide(0, A, *A.bounds.GetCenter()))
// 	equals(t, true, ok)

// 	B := NewCollider(*protometry.NewBox(0, 2, 0, 4, 3, 0))
// 	ok = o.Insert(*NewObjectCollide(0, B, *protometry.NewVectorN(3.5, 2.5, 0)))
// 	equals(t, true, ok)

// 	/*
// 		Casting a ray

// 			_ _ _ _
// 		   |||_|_|_|
// 		   |A|_|_|B|
// 		   |||_|_|_|
// 		   |||_|_|_|

// 		So we should have hit both A and B colliders
// 	*/
// 	equals(t, 2, len(*o.Raycast(*protometry.NewVector3Zero(), *protometry.NewVectorN(0, 1, 0), 4)))

// 	/*
// 		Casting a ray

// 			_ _ _ _
// 		   |_|_|||_|
// 		   |A|_|||B|
// 		   |_|_|||_|
// 		   |_|_|||_|

// 		So we should have hit only B collider
// 	*/
// 	equals(t, 1, len(*o.Raycast(*protometry.NewVectorN(2, 0, 0), *protometry.NewVectorN(0, 1, 0), 4)))
// }

// func BenchmarkOctreeNode_Raycast(b *testing.B) {
// 	size := float64(b.N)
// 	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
// 	for i := 0.; i < size; i++ {
// 		p1 := NewCollider(*protometry.NewBox(0, i, 0, 0, i+1, 0))
// 		o.Insert(*NewObjectCollide(0, p1, *p1.bounds.GetCenter()))
// 	}
// 	b.StartTimer()
// 	for i := 0.; i < size; i++ {
// 		o.Raycast(*protometry.NewVectorN(0, i, 0), *protometry.NewVectorN(0, 1, 0), 1)
// 	}
// 	b.StopTimer()
// }
