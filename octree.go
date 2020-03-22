package octree

import (
	"github.com/The-Tensox/protometry"
)

// Octree ...
type Octree struct {
	root *OctreeNode
}

func NewOctree(region *protometry.Box) *Octree {
	return &Octree{
		root: &OctreeNode{region: *region},
	}
}

func (o *Octree) Insert(point Point) bool {
	return o.root.Insert(point)
}

// Range return a list of points found inside the defined region
func (o *Octree) Range(region protometry.Box) []Point {
	return o.root.Range(region)
}

/*
// Cull ...
func (o *Octree) Cull(position protometry.VectorN) ([]interface{}, error) {
	return nil, nil
}

func (o *Octree) Search(position protometry.VectorN) (*OctreeNode, error) {
	in, err := position.In(o.root.region)
	if err != nil {
		return nil, err
	}
	if !in {
		return nil, ErrtreeOutsideBounds
	}

	return o.root.Search(position)
}


func (o *Octree) Remove(position protometry.VectorN) error {
	in, err := position.In(o.root.region)
	if err != nil {
		return err
	}
	if !in {
		return ErrtreeOutsideBounds
	}

	return o.root.Remove(position)
}
*/

/*
func (o *Octree) ToString() string {
	return o.root.ToString()
}
*/
