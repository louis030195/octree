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
	err := o.Insert(*NewPoint(1, 10, 10, 10))
	equals(t, false, err)
	err = o.Insert(*NewPoint(2, 3, 3, 3))
	equals(t, true, err)
	err = o.Insert(*NewPoint(3, 3, 3, 4))
	equals(t, true, err)
	err = o.Insert(*NewPoint(4, 3, 4, 4))
	equals(t, true, err)
	return o
}

func TestOctree_NewOctree(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), 1))
	// Should be [(-1,-1,-1), (1, 1, 1)]
	equals(t, *protometry.NewBox(-1, -1, -1, 1, 1, 1), o.root.region)
}

func TestOctreeNode_Insert(t *testing.T) {
	o := boilerplateTree(t)
	err := o.Insert(*NewPoint(5, 4, 4, 4))
	equals(t, true, err)
	err = o.Insert(*NewPoint(6, 1, 1, 1))
	equals(t, true, err)
	equals(t, 5, len(o.root.points))
	err = o.Insert(*NewPoint(7, 1, 1, 1))
	equals(t, true, err)
	if CAPACITY < 6 {
		equals(t, *NewPoint(6, 1, 1, 1), o.root.children[6].points[0])
	}

	// New octree
	size := 1000.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			err = o.Insert(*NewPoint(0, i, j, i))
			equals(t, true, err)
		}
	}
}

func BenchmarkOctreeNode_Insert(b *testing.B) {
	b.StartTimer()
	// New octree
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			o.Insert(*NewPoint(0, i, j, i))
		}
	}
	b.StopTimer()
}

func TestOctreeNode_Get(t *testing.T) {
	o := boilerplateTree(t)
	points := o.Get(1, 1, 1, 4, 4, 4)
	equals(t, 3, len(*points))
	points = o.Get(1, 1, 1, 3, 3, 3)
	equals(t, 1, len(*points))
	points = o.Get(1, 1, 1, 3, 3, 4)
	equals(t, 2, len(*points))
	points = o.Get(1, 1, 1, 3, 4, 4)
	equals(t, 3, len(*points))

	o = boilerplateTree(t)
	point := *o.Get(3, 3, 3)
	equals(t, *protometry.NewVectorN(3, 3, 3), point[0].position)
	equals(t, 2, point[0].data)
	point = *o.Get(3, 3, 4)
	equals(t, *protometry.NewVectorN(3, 3, 4), point[0].position)
	equals(t, 3, point[0].data)
	point = *o.Get(3, 4, 4)
	equals(t, *protometry.NewVectorN(3, 4, 4), point[0].position)
	equals(t, 4, point[0].data)
	var nilPointSlice *[]Point
	equals(t, nilPointSlice, o.Get(4, 4, 4))
}

func BenchmarkOctreeNode_GetMultiple(b *testing.B) {
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			o.Insert(*NewPoint(0, i, j, i))
		}
	}
	b.StartTimer()
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			o.Get(i, j, i, i, j, i)
		}
	}
	b.StopTimer()
}

func BenchmarkOctreeNode_GetOne(b *testing.B) {
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			o.Insert(*NewPoint(0, i, j, i))
		}
	}
	b.StartTimer()
	for i := 0.; i < size; i++ {
		for j := 0.; j < size; j++ {
			o.Get(i, j, i)
		}
	}
	b.StopTimer()
}

func TestOctree_Move(t *testing.T) {

}

func TestOctree_Raycast(t *testing.T) {
	o := NewOctree(protometry.NewBox(0, 0, 0, 10, 10, 10))
	p1 := protometry.NewBox(0, 0, 0, 0, 1, 0) // Line collider
	ok := o.Insert(*NewPointCollide(1, *p1, *p1.GetCenter()))
	equals(t, true, ok)
	p2 := protometry.NewBox(0, 2, 0, 0, 3, 0)
	ok = o.Insert(*NewPointCollide(2, *p2, *p2.GetCenter()))
	equals(t, true, ok)
	p3 := protometry.NewBox(0, 4, 0, 0, 5, 0)
	ok = o.Insert(*NewPointCollide(3, *p3, *p3.GetCenter()))
	equals(t, true, ok)

	// Cast a ray toward up from 0;0;0 of length 10
	points := *o.Raycast(*protometry.NewVector3Zero(), *protometry.NewVectorN(0, 1, 0), 10)
	equals(t, 3, len(points))

	// Cast a ray toward up from 0;0;0 of length 2.5
	points = *o.Raycast(*protometry.NewVector3Zero(), *protometry.NewVectorN(0, 1, 0), 2.5)
	equals(t, 2, len(points))

	// Cast a ray toward up from 0;2.1;0 of length 7
	points = *o.Raycast(*protometry.NewVectorN(0, 2.1, 0), *protometry.NewVectorN(0, 1, 0), 7)
	equals(t, 2, len(points))
	equals(t, p3, &points[len(points)-1].collider)

	// Cast a ray toward up from 0;5.1;0 of length 4
	equals(t, 0, len(*o.Raycast(*protometry.NewVectorN(0, 5.1, 0), *protometry.NewVectorN(0, 1, 0), 4)))

	// New octree
	size := 1000.
	o = NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		p1 := protometry.NewBox(0, i, 0, 0, i+1, 0)
		o.Insert(*NewPointCollide(0, *p1, *p1.GetCenter()))
		o.Raycast(*protometry.NewVectorN(0, i, 0), *protometry.NewVectorN(0, 1, 0), 1)
	}
	equals(t, 1000, len(*o.Raycast(*protometry.NewVectorN(0, 0, 0), *protometry.NewVectorN(0, 1, 0), size)))
}

func BenchmarkOctreeNode_Raycast(b *testing.B) {
	size := float64(b.N)
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), size))
	for i := 0.; i < size; i++ {
		p1 := protometry.NewBox(0, i, 0, 0, i+1, 0)
		o.Insert(*NewPointCollide(0, *p1, *p1.GetCenter()))
	}
	b.StartTimer()
	for i := 0.; i < size; i++ {
		o.Raycast(*protometry.NewVectorN(0, i, 0), *protometry.NewVectorN(0, 1, 0), 1)
	}
	b.StopTimer()
}
