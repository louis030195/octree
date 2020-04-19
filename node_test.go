package octree

import (
	"github.com/The-Tensox/protometry"
	"testing"
)

func TestNode_merge(t *testing.T) {
	n := Node{
		objects:  nil,
		region:   protometry.Box{},
		children: &[8]Node{},
	}
	equals(t, true, n.merge())
	var nilChildren *[8]Node
	equals(t, nilChildren, n.children)
}