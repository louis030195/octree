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

	equals(t, true, o.Remove(*myObj) != nil)
	var nilChildren *[8]OctreeNode
	// Shouldn't have merged
	equals(t, true, nilChildren != o.root.children)
	// One less object
	equals(t, 8, o.GetNumberOfObjects())
	equals(t, true, o.Remove(objects[len(objects)-1]) == nil) // We've already removed it
	equals(t, 8, o.GetNumberOfObjects())
	equals(t, true, o.Remove(objects[len(objects)-2]) != nil)
	equals(t, 7, o.GetNumberOfObjects())
}

func TestOctree_Move(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), 10))
	myObj := NewObjectCube(0, 0, 0, 0, 1)
	equals(t, *protometry.NewBox(-1, -1, -1, 1, 1, 1), myObj.bounds)
	equals(t, true, o.Insert(*myObj))
	equals(t, 1, o.GetNumberOfObjects())
	myObj = o.Move(*myObj, 0, 0, 0, 2, 2, 2) // Using bounds
	equals(t, true, myObj != nil)
	equals(t, *protometry.NewBox(0, 0, 0, 2, 2, 2), myObj.bounds)
	equals(t, 1, o.GetNumberOfObjects())
	myObj = o.Move(*myObj, 3, 3, 3) // Using position
	equals(t, true, myObj != nil)
	equals(t, *protometry.NewBox(2, 2, 2, 4, 4, 4), myObj.bounds)
	equals(t, 1, o.GetNumberOfObjects())
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
