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
	equals(t, nil, err)
	//equals(t, nil, n)
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

	t.Log(o.ToString())
}

func TestOctreeNode_findBranch(t *testing.T) {
	o := NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(), 1))
	t.Log(o.root.findBranch(*protometry.NewVector3One()))
	t.Log(o.root.findBranch(*protometry.NewVector3Zero()))
	t.Log(o.root.findBranch(*protometry.NewVector3One().Mul(0.1)))
	t.Log(o.root.findBranch(*protometry.NewVector3One().Mul(-1)))

}
