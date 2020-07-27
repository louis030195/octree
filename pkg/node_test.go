package octree

import (
    "github.com/louis030195/protometry/api/volume"

    "testing"
)

func TestNode_merge(t *testing.T) {
	n := Node{
		objects:  nil,
		region:   volume.Box{},
		children: &[8]Node{},
	}
	equals(t, true, n.merge())
	var nilChildren *[8]Node
	equals(t, nilChildren, n.children)
}
