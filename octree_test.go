package octree

import (
	"fmt"
	"math"
	"math/rand"
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

/*
 Returns true if all nodes have less objects than CAPACITY
 */
func ensureBalanced(tb testing.TB, n Node) bool {
	if len(n.objects) > CAPACITY {
		tb.Logf("Number of objects in node: %v", len(n.objects))
		return false
	}
	if n.children != nil {
		for _, c := range n.children {
			if r := ensureBalanced(tb, c); !r {
				return r
			}
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

func checkOctree(t *testing.T, o Octree, expectedObjects int) {
	// t.Log(o.toString(false))
	t.Logf("Octree height: %v", o.getHeight())
	t.Logf("Octree usage: %v", o.getUsage())
	t.Logf("Octree nodes: %v", o.getNumberOfNodes())
	t.Logf("Octree objects: %v", o.getNumberOfObjects())
	t.Logf("Octree is balanced %v", ensureBalanced(t, *o.root))
	equals(t, expectedObjects, o.getNumberOfObjects())
	equals(t, true, o.getUsage() < 1)
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
	equals(t, 10, o.getNumberOfObjects())
	// equals(t, 0, len(o.root.objects) < CAPACITY)
	// t.Log(o.toString(false))
	// equals(t, 16, o.getNumberOfNodes()) // FIXME
	// Let's test with more scale
	size = 100.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			equals(t, true, o.Insert(*NewObjectCube(0, i, j, i, 2)))
			equals(t, true, o.Insert(*NewObjectCube(0, i, j, i, 2)))
		}
	}
	checkOctree(t, *o, int(size*size*2))
}

func TestNode_InsertRandomPosition(t *testing.T) {
	size := math.Pow(10, 4)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		p := protometry.RandomSpherePoint(*protometry.NewVector3Zero(), size-1)
		equals(t, true, o.Insert(*NewObjectCube(0, p.Get(0), p.Get(1), p.Get(2), 1)))
	}
	checkOctree(t, *o, int(size))
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

func TestNode_GetCollidingTwo(t *testing.T) {
	size := 100.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		p := protometry.RandomSpherePoint(*protometry.NewVector3Zero(), size-1)
		equals(t, true, o.Insert(*NewObjectCube(0, p.Get(0), p.Get(1), p.Get(2), 1)))
	}
	equals(t, int(size), o.getNumberOfObjects())
	equals(t, int(size), len(o.GetColliding(*protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))))
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
	equals(t, int(size-1), o.getNumberOfObjects())
	// equals(t, int(size/8)-CAPACITY+1, o.getNumberOfNodes()) // FIXME
	for i := range objects {
		equals(t, true, o.Remove(objects[i]))
	}
	equals(t, 0, o.getNumberOfObjects())
	equals(t, 1, o.getNumberOfNodes()) // Only root left
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	objects = []Object{}
	for i := 0.; i < 9; i++ {
		myObj = NewObjectCube(0, i, i, i, 2)
		objects = append(objects, *myObj)
		equals(t, true, o.Insert(*myObj))
	}
	equals(t, 9, o.getNumberOfObjects())
	equals(t, 8, len(o.root.children))

	equals(t, true, o.Remove(*myObj))
	var nilChildren *[8]Node
	// Shouldn't have merged
	equals(t, true, nilChildren != o.root.children)
	// One less object
	equals(t, 8, o.getNumberOfObjects())
	equals(t, false, o.Remove(objects[len(objects)-1])) // We've already removed it
	equals(t, 8, o.getNumberOfObjects())
	equals(t, true, o.Remove(objects[len(objects)-2]))
	equals(t, 7, o.getNumberOfObjects())
}

func TestOctree_Move(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), 20))
	myObj := NewObjectCube(0, 0, 0, 0, 2)
	equals(t, *protometry.NewBoxMinMax(-1, -1, -1, 1, 1, 1), myObj.Bounds)
	equals(t, true, o.Insert(*myObj))
	equals(t, 1, o.getNumberOfObjects())
	// Using Bounds
	equals(t, true, o.Move(myObj, 0, 0, 0, 2, 2, 2))
	equals(t, *protometry.NewBoxMinMax(0, 0, 0, 2, 2, 2), myObj.Bounds)
	equals(t, 1, o.getNumberOfObjects())
	// Using position
	equals(t, true, o.Move(myObj, 3, 3, 3))
	equals(t, *protometry.NewBoxMinMax(2, 2, 2, 4, 4, 4), myObj.Bounds)
	equals(t, 1, o.getNumberOfObjects())
}

func TestOctree_GetAllObjects(t *testing.T) {
	for i := 0.; i < 100; i++ {
		o := octreeRandomInsertions(t, i)
		equals(t, int(i), len(o.GetAllObjects()))
	}
}

func TestOctree_Range(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVectorN(0, 0, 0), 100))

	objs := []*Object{
		NewObjectCube(0, 12, 12, 12, 1),
		NewObjectCube(0, 15, 12, 12, 1),
		NewObjectCube(0, 12, 27, 12, 1),
	}

	for i := range objs {
		equals(t, true, o.Insert(*objs[i]))
	}

	// Asserting that the first element is objs[0]
	o.Range(func(object *Object) bool {
		equals(t, true, objs[0].Equal(*object))
		return false
	})
	i := 0
	o.Range(func(object *Object) bool {
		equals(t, true, objs[i].Equal(*object))
		i++
		return true
	})

	i = 0
	// Asserting that we properly range over objs
	o.Range(func(object *Object) bool {
		equals(t, true, objs[0].Equal(*object))
		i++
		return false
	})
	// Assert that returning false properly stop the iteration
	equals(t, 1, i)
}

func TestOctree_MoveIncorrectDims(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), 20))
	myObj := NewObjectCube(0, 0, 0, 0, 2)
	equals(t, *protometry.NewBoxMinMax(-1, -1, -1, 1, 1, 1), myObj.Bounds)
	equals(t, true, o.Insert(*myObj))
	equals(t, 1, o.getNumberOfObjects())
	equals(t, false, o.Move(myObj, 0, 0, 0, 2, 2))
	equals(t, *protometry.NewBoxMinMax(-1, -1, -1, 1, 1, 1), myObj.Bounds)
	equals(t, 1, o.getNumberOfObjects())
}

func octreeRandomInsertions(t testing.TB, treeSize float64) *Octree {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), treeSize*2))
	for i := 0.; i < treeSize; i++ {
		p := protometry.RandomSpherePoint(*protometry.NewVector3Zero(), treeSize-1)
		equals(t, true, o.Insert(*NewObjectCube(0, p.Get(0), p.Get(1), p.Get(2), 1)))
	}
	return o
}

func TestOctree_GetNodes(t *testing.T) {
	ts := 1000.
	o := octreeRandomInsertions(t, ts)
	t.Log(o)
	//equals(t, ?, len(o.GetNodes())) // TODO
}

func TestOctree_GetSize(t *testing.T) {
	for i := 0.; i < 100; i++ {
		o := octreeRandomInsertions(t, i)
		equals(t, int(i), o.GetSize()/2)
	}
}


func TestOctree_GetHeight(t *testing.T) {
	size := 1000.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	// equals(t, int((int(size)/CAPACITY)/8), o.getHeight()) // FIXME
}

func TestOctree_GetNumberOfNodes(t *testing.T) {
	size := 1000.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	// equals(t, ?, o.getNumberOfNodes()) // FIXME
}

func TestOctree_GetNumberOfObjects(t *testing.T) {
	size := 1000.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	equals(t, int(size), o.getNumberOfObjects())
}

func TestOctree_GetUsage(t *testing.T) {
	size := 1000.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	// Any better tests ?
	equals(t, float64(o.getNumberOfObjects())/float64(o.getNumberOfNodes()*CAPACITY), o.getUsage())
}


func TestOctree_ToString(t *testing.T) {
	size := 20.
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 0.; i < size; i++ {
		equals(t, true, o.Insert(*NewObjectCube(0, i, i, i, 2)))
	}
	t.Log(o.toString(false))
}

/* * * BENCHES * * */
func BenchmarkNode_InsertRandomPosition(b *testing.B) {
	b.StartTimer()
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 1.; i < size; i++ {
		p := protometry.RandomSpherePoint(*protometry.NewVector3Zero(), size-1)
		equals(b, true, o.Insert(*NewObjectCube(0, p.Get(0), p.Get(1), p.Get(2),
			1)))
	}
	b.StopTimer()
}


func BenchmarkNode_GetCollidingFullRandom(b *testing.B) {
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	for i := 1.; i < size; i++ {
		p := protometry.RandomSpherePoint(*protometry.NewVector3Zero(), size-1)
		equals(b, true, o.Insert(*NewObjectCube(0, p.Get(0), p.Get(1), p.Get(2),
			1)))
	}
	b.StartTimer()
	for i := 1.; i < size; i++ {
		p := protometry.RandomSpherePoint(*protometry.NewVector3Zero(), size-1)
		o.GetColliding(*protometry.NewBoxOfSize(p, rand.ExpFloat64() / size))
	}
	b.StopTimer()
}

func BenchmarkNode_RemoveRandomPosition(b *testing.B) {
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*2))
	var objects []Object
	for i := 1.; i < size; i++ {
		p := protometry.RandomSpherePoint(*protometry.NewVector3Zero(), size-1)
		ob := NewObjectCube(0, p.Get(0), p.Get(1), p.Get(2), 1)
		equals(b, true, o.Insert(*ob))
		objects = append(objects, *ob)
	}
	b.StartTimer()
	for i := 1.; i < size-1; i++ {
		equals(b, true, o.Remove(objects[int(i)]))
	}
	b.StopTimer()
}

func BenchmarkNode_MoveRandomPosition(b *testing.B) {
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size*4)) // x4 because moving ++
	var objects []Object
	for i := 0.; i < size; i++ {
		p := protometry.RandomSpherePoint(*protometry.NewVector3Zero(), size-1)
		ob := NewObjectCube(0, p.Get(0), p.Get(1), p.Get(2), 1)
		equals(b, true, o.Insert(*ob))
		objects = append(objects, *ob)
	}
	b.StartTimer()
	for i := 0.; i < size-1; i++ {
		ob := objects[int(i)]
		p := protometry.RandomSpherePoint(*protometry.NewVector3Zero(), size-1)
		equals(b, true, o.Move(&ob, p.Dimensions...))
	}
	b.StopTimer()
}

func BenchmarkOctree_Range(b *testing.B) {
	size := float64(b.N)
	o := octreeRandomInsertions(b, size)
	b.StartTimer()
	o.Range(func(object *Object) bool {
		return true
	})
	b.StopTimer()
}

