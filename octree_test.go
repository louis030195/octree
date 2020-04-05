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

func ensureBalanced(tb testing.TB, n Node) bool {
	if len(n.objects) > CAPACITY {
		tb.Logf("Number of objects in node: %v", len(n.objects))
		return false
	}
	for _, c := range n.children {
		if r := ensureBalanced(tb, c); !r {
			return r
		}
	}
	return true
}

func boilerplateTree(t *testing.T) *Octree {
	o := NewOctree(protometry.NewBoxMinMax(1, 1, 1, 4, 4, 4))
	ok := o.Insert(*NewObjectCube(0, 2, 2, 3, 1))
	equals(t, true, ok)
	return o
}

func TestOctree_NewOctree(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), 2))
	// Should be [(-1,-1,-1), (1, 1, 1)]
	equals(t, *protometry.NewBoxMinMax(-1, -1, -1, 1, 1, 1), o.root.region)
}

func TestNode_Insert(t *testing.T) {
	o := boilerplateTree(t)
	equals(t, true, o.Insert(*NewObjectCube(5, 3, 3, 3, 2)))
	equals(t, true, o.Insert(*NewObjectCube(6, 2, 2, 2, 2)))
	equals(t, 3, len(o.root.objects))
	equals(t, true, o.Insert(*NewObjectCube(7, 2, 2, 2, 2)))

	// Go over capacity threshold, force a split
	size := 10.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	// We inserted 10 objects so we should have 10 objects ;)
	equals(t, 10, o.GetNumberOfObjects())
	// equals(t, 16, o.GetNumberOfNodes()) // FIXME
	// Let's test with more scale
	size = 100.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			equals(t, true, o.Insert(*NewObjectCube(0, i, j, i, 2)))
			equals(t, true, o.Insert(*NewObjectCube(0, i, j, i, 2)))
		}
	}
	t.Logf("Octree height: %v", o.GetHeight())
	t.Logf("Octree usage: %v", o.GetUsage())
	t.Logf("Octree nodes: %v", o.GetNumberOfNodes())
	t.Logf("Octree objects: %v", o.GetNumberOfObjects())
	equals(t, int(size*size*2), o.GetNumberOfObjects())
	equals(t, true, ensureBalanced(t, *o.root))
	equals(t, true, o.GetUsage() < 1)
}

func TestNode_GetColliding(t *testing.T) {
	o := NewOctree(protometry.NewBoxMinMax(1, 1, 1, 4, 4, 4))
	equals(t, true, o.Insert(*NewObjectCube(0, 2, 2, 3, 1)))
	equals(t, true, o.Insert(*NewObjectCube(5, 3, 3, 3, 2)))
	equals(t, true, o.Insert(*NewObjectCube(6, 2, 2, 2, 2)))
	equals(t, 3, len(o.root.objects))
	equals(t, true, o.Insert(*NewObjectCube(7, 2, 2, 2, 2)))

	colliders := o.GetColliding(*protometry.NewBoxMinMax(0, 0, 0, 0.9, 0.9, 0.9))
	equals(t, 0, len(colliders))
	colliders = o.GetColliding(*protometry.NewBoxMinMax(0, 0, 0, 1, 1, 1))
	equals(t, 2, len(colliders))
	colliders = o.GetColliding(*protometry.NewBoxMinMax(1, 1, 1, 1.1, 1.1, 1.1))
	equals(t, 2, len(colliders))
	equals(t, 6, colliders[0].Data)
	equals(t, 7, colliders[1].Data)

	/* * * */
	o = NewOctree(protometry.NewBoxMinMax(-10, -10, -10, 10, 10, 10))
	equals(t, true, o.Insert(*NewObjectCube(0, 0, 0, 0, 2)))
	equals(t, true, o.Insert(*NewObjectCube(1, 0, 2, 0, 2)))
	equals(t, true, o.Insert(*NewObjectCube(2, 0, 4, 0, 2)))
	colliders = o.GetColliding(*protometry.NewBoxMinMax(-2, -2, -2, -1.1, -1.1, -1.1))
	equals(t, 0, len(colliders))
	colliders = o.GetColliding(*protometry.NewBoxMinMax(-2, -2, -2, -1, -1, -1))
	equals(t, 1, len(colliders))
	colliders = o.GetColliding(*protometry.NewBoxMinMax(-10, -10, -10, 10, 10, 10))
	equals(t, 3, len(colliders))
	equals(t, 0, colliders[0].Data)
	// Reverse
	colliders = o.GetColliding(*protometry.NewBoxMinMax(10, 10, 10, -10, -10, -10)) // FIXME
	//equals(t, 3, len(colliders))
	//equals(t, 0, colliders[0].Data)
}

func TestOctree_Remove(t *testing.T) {
	o := boilerplateTree(t)
	myObj := NewObjectCube(27, 2, 2, 3, 2)
	equals(t, true, o.Insert(*myObj))
	equals(t, true, o.Remove(*myObj))

	myObj = NewObjectCube(27, 2, 2, 3, 2)
	equals(t, true, o.Insert(*myObj))
	// We didn't insert this one !
	equals(t, false, o.Remove(*NewObjectCube(12, 2, 2, 3, 2)))

	// New octree
	size := 1000.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	var objects []Object
	for i := 1.; i < size; i++ {
		myObj = NewObjectCube(0, i, i, i, 2)
		equals(t, true, o.Insert(*myObj))
		objects = append(objects, *myObj)
	}
	equals(t, int(size-1), o.GetNumberOfObjects())
	// equals(t, int(size/8)-CAPACITY+1, o.GetNumberOfNodes()) // FIXME
	for i := range objects {
		equals(t, true, o.Remove(objects[i]))
	}
	equals(t, 0, o.GetNumberOfObjects())
	equals(t, 1, o.GetNumberOfNodes()) // Only root left
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	objects = []Object{}
	for i := 0.; i < 9; i++ {
		myObj = NewObjectCube(0, i, i, i, 2)
		objects = append(objects, *myObj)
		equals(t, true, o.Insert(*myObj))
	}
	equals(t, 9, o.GetNumberOfObjects())
	equals(t, 8, len(o.root.children))

	equals(t, true, o.Remove(*myObj))
	var nilChildren *[8]Node
	// Shouldn't have merged
	equals(t, true, nilChildren != o.root.children)
	// One less object
	equals(t, 8, o.GetNumberOfObjects())
	equals(t, false, o.Remove(objects[len(objects)-1])) // We've already removed it
	equals(t, 8, o.GetNumberOfObjects())
	equals(t, true, o.Remove(objects[len(objects)-2]))
	equals(t, 7, o.GetNumberOfObjects())
}

func TestOctree_Move(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), 20))
	myObj := NewObjectCube(0, 0, 0, 0, 2)
	equals(t, *protometry.NewBoxMinMax(-1, -1, -1, 1, 1, 1), myObj.Bounds)
	equals(t, true, o.Insert(*myObj))
	equals(t, 1, o.GetNumberOfObjects())
	// Using Bounds
	equals(t, true, o.Move(myObj, 0, 0, 0, 2, 2, 2))
	equals(t, *protometry.NewBoxMinMax(0, 0, 0, 2, 2, 2), myObj.Bounds)
	equals(t, 1, o.GetNumberOfObjects())
	// Using position
	equals(t, true, o.Move(myObj, 3, 3, 3))
	equals(t, *protometry.NewBoxMinMax(2, 2, 2, 4, 4, 4), myObj.Bounds)
	equals(t, 1, o.GetNumberOfObjects())
}

func TestOctree_GetHeight(t *testing.T) {
	size := 1000.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	// equals(t, int((int(size)/CAPACITY)/8), o.GetHeight()) // FIXME
}

func TestOctree_GetNumberOfNodes(t *testing.T) {
	size := 1000.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	// equals(t, ?, o.GetNumberOfNodes()) // FIXME
}

func TestOctree_GetNumberOfObjects(t *testing.T) {
	size := 1000.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	equals(t, int(size), o.GetNumberOfObjects())
}

func TestOctree_GetUsage(t *testing.T) {
	size := 1000.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	// Any better tests ?
	equals(t, float64(o.GetNumberOfObjects())/float64(o.GetNumberOfNodes()*CAPACITY), o.GetUsage())
}


func TestOctree_ToString(t *testing.T) {
	size := 20.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	t.Log(o.ToString(false))
}

/* * * BENCHES * * */
func BenchmarkNode_Insert(b *testing.B) {
	b.StartTimer()
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 1.; i < size; i++ {
		equals(b, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	b.StopTimer()
}

func BenchmarkNode_GetColliding(b *testing.B) {
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 1.; i < size; i++ {
		equals(b, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	b.StartTimer()
	for i := 1.; i < size; i++ {
		o.GetColliding(*protometry.NewBoxOfSize(*protometry.NewVectorN(i, i, i), 1))
	}
	b.StopTimer()
}

func BenchmarkNode_Remove(b *testing.B) {
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	var objects []Object
	for i := 1.; i < size; i++ {
		ob := NewObjectCube(0, i, i, i, 1)
		equals(b, true, o.Insert(*ob))
		objects = append(objects, *ob)
	}
	b.StartTimer()
	for i := 1.; i < size-1; i++ {
		equals(b, true, o.Remove(objects[int(i)]))
	}
	b.StopTimer()
}

func BenchmarkNode_Move(b *testing.B) {
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*4)) // x4 because moving ++
	var objects []Object
	for i := 1.; i < size; i++ {
		ob := NewObjectCube(0, i, i, i, 2)
		equals(b, true, o.Insert(*ob))
		objects = append(objects, *ob)
	}
	b.StartTimer()
	for i := 1.; i < size-1; i++ {
		ob := objects[int(i)]
		equals(b, true, o.Move(&ob, ob.Bounds.Center.Scale(1.1).Dimensions...))
	}
	b.StopTimer()
}
