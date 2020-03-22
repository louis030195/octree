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

func TestOctree_NewOctree(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), 1))
	// Should be [(-1,-1,-1), (1, 1, 1)]
	equals(t, protometry.NewBox(*protometry.NewVector3One().Mul(-1), *protometry.NewVector3One()), o.root.region)
}

func TestOctreeNode_Insert(t *testing.T) {
	o := NewOctree(protometry.NewBox(*protometry.NewVector3One(), *protometry.NewVectorN(4, 4, 4)))
	var r interface{}
	r = 0
	err := o.Insert(NewPoint(10, 10, 10, r))
	equals(t, false, err)
	err = o.Insert(NewPoint(3, 3, 3, r))
	equals(t, true, err)
	err = o.Insert(NewPoint(3, 3, 4, r))
	equals(t, true, err)
	err = o.Insert(NewPoint(3, 4, 4, r))

	err = o.Insert(NewPoint(4, 4, 4, r))
	equals(t, true, err)
	err = o.Insert(NewPoint(1, 1, 1, r))
	equals(t, true, err)
	equals(t, 5, len(o.root.points))
	err = o.Insert(NewPoint(1, 1, 1, r))
	equals(t, true, err)
	if CAPACITY < 6 {
		equals(t, NewPoint(1, 1, 1, r), o.root.children[6].points[0])
	}

	// New octree
	size := 1000.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			err = o.Insert(NewPoint(i, j, i, r))
			equals(t, true, err)
		}
	}
}

/*
func TestOctreeNode_Search(t *testing.T) {
	o := NewOctree(protometry.NewBox(*protometry.NewVector3One(), *protometry.NewVectorN(4, 4, 4)))
	err := o.Insert(*protometry.NewVectorN(10, 10, 10), []interface{}{})
	equals(t, ErrtreeOutsideBounds, err)
	err = o.Insert(*protometry.NewVectorN(3, 3, 3), []interface{}{1})
	equals(t, nil, err)
	err = o.Insert(*protometry.NewVectorN(3, 3, 4), []interface{}{2})
	equals(t, nil, err)
	n, err := o.Search(*protometry.NewVectorN(3, 3, 3))
	equals(t, nil, err)
	equals(t, protometry.NewVectorN(3, 3, 3), n.position)
	equals(t, []interface{}{1}, n.data)
	n, err = o.Search(*protometry.NewVectorN(3, 4, 4))
	equals(t, ErrtreeFailedToFindNode, err)
	err = o.Insert(*protometry.NewVectorN(3, 4, 4), []interface{}{3})
	equals(t, nil, err)
	n, err = o.Search(*protometry.NewVectorN(3, 4, 4))
	equals(t, protometry.NewVectorN(3, 4, 4), n.position)
	equals(t, []interface{}{3}, n.data)

	err = o.Insert(*protometry.NewVectorN(4, 4, 4), []interface{}{4})
	equals(t, nil, err)
	n, err = o.Search(*protometry.NewVectorN(4, 4, 4))
	equals(t, protometry.NewVectorN(4, 4, 4), n.position)
	equals(t, []interface{}{4}, n.data)
	err = o.Insert(*protometry.NewVector3One(), []interface{}{1})
	equals(t, nil, err)
	err = o.Insert(*protometry.NewVector3One(), []interface{}{2})
	equals(t, nil, err)

	// New octree
	size := 100.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			err = o.Insert(*protometry.NewVectorN(i, j, i), []interface{}{0})
			equals(t, nil, err)
		}
	}
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			_, err = o.Search(*protometry.NewVectorN(i, j, i))
			equals(t, nil, err)
		}
	}
}

func TestOctreeNode_Remove(t *testing.T) {
	o := NewOctree(protometry.NewBox(*protometry.NewVector3One(), *protometry.NewVectorN(4, 4, 4)))
	err := o.Insert(*protometry.NewVectorN(10, 10, 10), []interface{}{})
	equals(t, ErrtreeOutsideBounds, err)
	err = o.Insert(*protometry.NewVectorN(3, 3, 3), []interface{}{1})
	equals(t, nil, err)
	err = o.Remove(*protometry.NewVectorN(3, 3, 3))
	equals(t, nil, err)
	_, err = o.Search(*protometry.NewVectorN(3, 3, 3))
	equals(t, ErrtreeFailedToFindNode, err)
}

func TestOctreeNode_getOctant(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), 1))
	o.root.position = protometry.NewVector3Zero()
	t.Log(o.root.getOctant(*protometry.NewVector3One()))
	t.Log(o.root.getOctant(*protometry.NewVector3Zero()))
	t.Log(o.root.getOctant(*protometry.NewVector3One().Mul(0.1)))
	t.Log(o.root.getOctant(*protometry.NewVector3One().Mul(-1)))
}

func TestOctreeNode_getNewRegion(t *testing.T) {
	pos := protometry.NewVector3Zero()
	o := &OctreeNode{
		position: pos,
		region:   protometry.NewBox(*pos, *protometry.NewVectorN(1, 1, 1)),
	}
	b := o.getNewRegion(0)
	equals(t, protometry.NewBox(*protometry.NewVectorN(0, 0, 0), *protometry.NewVectorN(0.5, 0.5, 0.5)), b.region)
}

func BenchmarkOctreeNode_Insert(b *testing.B) {
	b.StartTimer()
	// New octree
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			o.Insert(*protometry.NewVectorN(i, j, i), []interface{}{0})
		}
	}
	b.StopTimer()
}

func BenchmarkOctreeNode_Search(b *testing.B) {
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			o.Insert(*protometry.NewVectorN(i, j, i), []interface{}{0})
		}
	}
	b.ResetTimer()
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			o.Search(*protometry.NewVectorN(i, j, i))
		}
	}
}
*/
