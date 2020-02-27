package octree

import (
	"fmt"

	"github.com/The-Tensox/protometry"
)

type OctreeNode struct {
	position *protometry.VectorN
	data     []interface{}
	region   *protometry.Box
	children [8]*OctreeNode
}

func NewEmptyOctreeNode() *OctreeNode {
	return &OctreeNode{position: protometry.NewVectorN(-1, -1, -1)}
}

func NewPointOctreeNode(position *protometry.VectorN, data []interface{}) *OctreeNode {
	return &OctreeNode{position: position, data: data}
}

func NewRegionOctreeNode(region *protometry.Box) *OctreeNode {
	o := NewEmptyOctreeNode()
	o.position = nil
	o.region = region
	for i := TLF; i <= BLB; i++ {
		o.children[i] = NewEmptyOctreeNode()
	}
	return o
}

func (o *OctreeNode) Search(position protometry.VectorN) (*OctreeNode, error) {
	pos := o.findBranch(position)
	if o.children[pos].position == nil {
		// Region node
		return o.children[pos].Search(position)
	} else if o.children[pos].position.Dimensions[0] == -1 {
		// Empty node
		return nil, nil
	}
	eq, err := o.children[pos].position.ApproxEqual(position)
	if err != nil {
		return nil, err
	}
	if eq {
		return o.children[pos], nil
	}
	return nil, nil
}
func (o *OctreeNode) Insert(position protometry.VectorN, data []interface{}) error {
	/*
		if o.size < 1 {
			return nil, nil
		}

		// If this node can host the new node
		in, err := position.In(*protometry.NewBoxOfSize(*o.position, o.size)) // Hardcoded 1 size
		// If in failed, return error
		if err != nil {
			return nil, err
		}
		// Current node is a good fit
		if in {
			if o.children[0] == nil {
				if o.position == nil {
					o.position
				}

				eq, err := o.position.ApproxEqual(position)

				// If equal failed, return error
				if err != nil {
					return nil, err
				}

				// Current node is at the wanted position
				if eq {
					o.data = append(o.data, data...) // Concat arrays
					return o, nil
				}

				// No children
				if o.children[0] == nil {
					// Create children
					subBoxes := protometry.NewBoxOfSize(*o.position, o.size).MakeSubBoxes()
					for i := 0; i < 8; i++ {
						o.children[i] = &OctreeNode{position: subBoxes[i].GetCenter(), size: o.size / 8} // size / 8 ?
					}
				}
				// Try to insert in the branch corresponding to the position
				return o.children[o.findBranch(position)].Insert(position, data)
			}
		}
	*/

	// TODO: move node rebalance
	/*
		if o == nil {
			return nil, nil
		}

		if o.children[0] == nil {
			if o.position != nil {
				eq, err := o.position.ApproxEqual(position)

				// If equal failed, return error
				if err != nil {
					return nil, err
				}

				// Current node is at the wanted position
				if eq {
					o.data = append(o.data, data...) // Concat arrays
					return o, nil
				}
				subBoxes := protometry.NewBoxOfSize(*o.position, o.region.GetSize()).MakeSubBoxes()

				// We have a position
				for i := 0; i < 8; i++ {
					o.children[i] = &OctreeNode{position: nil, region: subBoxes[i]}
				}
			} else {
				// Leaf with no position
				o.position = &position
				o.data = data
			}
		}

		if o.position == nil {
			return nil, nil
		}
		// Try to insert in the branch corresponding to the position
		return o.children[o.findBranch(position)].Insert(position, data)
	*/

	// Find the proper direction to insert
	branch := o.findBranch(position)

	// Two point on same position
	if o.children[branch].position != nil {
		eq, err := o.children[branch].position.ApproxEqual(position)
		if err != nil {
			return err
		}
		if eq {
			o.children[branch].data = append(o.children[branch].data, data...)
			return nil
		}
	}

	if o.children[branch].position == nil {
		// If region node, insert in a child
		return o.children[branch].Insert(position, data)
	} else if o.children[branch].position.Dimensions[0] == -1 {
		// If empty node, create node with new data on this leaf
		o.children[branch] = NewPointOctreeNode(&position, data)
	} else {
		// If point node, store its data, make it region node,
		// move stored data down to children
		// insert new data in children
		p := *o.children[branch].position
		d := o.children[branch].data
		// Make it region node
		o.children[branch] = o.findPosition(branch, *o.region)
		// Find new leaf for old node
		o.children[branch].Insert(p, d)
		// Find leaf for new node
		return o.children[branch].Insert(position, data)
	}
	return nil
}

func (o *OctreeNode) Remove(position protometry.VectorN) (bool, error) {
	return false, nil
}

const (
	TLF = iota // top left front
	TRF        // top right front
	BRF        // bottom right front
	BLF        // bottom left front
	TLB        // top left back
	TRB        // top right back
	BRB        // bottom right back
	BLB        // bottom left back
)

func (o *OctreeNode) findBranch(position protometry.VectorN) int {
	var pos int
	center := o.region.GetCenter()
	if position.Dimensions[0] <= center.Dimensions[0] {
		if position.Dimensions[1] <= center.Dimensions[1] {
			if position.Dimensions[2] <= center.Dimensions[2] {
				pos = TLF
			} else {
				pos = TLB
			}
		} else {
			if position.Dimensions[2] <= center.Dimensions[2] {
				pos = BLF
			} else {
				pos = BLB
			}
		}
	} else {
		if position.Dimensions[1] <= center.Dimensions[1] {
			if position.Dimensions[2] <= center.Dimensions[2] {
				pos = TRF
			} else {
				pos = TRB
			}
		} else {
			if position.Dimensions[2] <= center.Dimensions[2] {
				pos = BRF
			} else {
				pos = BRB
			}
		}
	}
	return pos
}

func (o *OctreeNode) findPosition(branch int, region protometry.Box) *OctreeNode {
	center := region.GetCenter()
	var newNode *OctreeNode
	switch branch {
	case TLF:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				region.GetMin().Dimensions[0],
				region.GetMin().Dimensions[1],
				region.GetMin().Dimensions[2]),
			*protometry.NewVectorN(
				center.Dimensions[0],
				center.Dimensions[1],
				center.Dimensions[2])))
	case TRF:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				center.Dimensions[0]+1,
				region.GetMin().Dimensions[1],
				region.GetMin().Dimensions[2]),
			*protometry.NewVectorN(
				region.GetMax().Dimensions[0],
				center.Dimensions[1],
				center.Dimensions[2])))
	case BRF:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				center.Dimensions[0]+1,
				center.Dimensions[1]+1,
				region.GetMin().Dimensions[2]),
			*protometry.NewVectorN(
				region.GetMax().Dimensions[0],
				region.GetMax().Dimensions[1],
				center.Dimensions[2])))
	case BLF:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				region.GetMin().Dimensions[0],
				center.Dimensions[1]+1,
				region.GetMin().Dimensions[2]),
			*protometry.NewVectorN(
				center.Dimensions[0],
				region.GetMax().Dimensions[1],
				center.Dimensions[2])))
	case TLB:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				region.GetMin().Dimensions[0],
				region.GetMin().Dimensions[1],
				center.Dimensions[2]+1),
			*protometry.NewVectorN(
				center.Dimensions[0],
				center.Dimensions[1],
				region.GetMax().Dimensions[2])))
	case TRB:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				center.Dimensions[0]+1,
				region.GetMin().Dimensions[1],
				center.Dimensions[2]+1),
			*protometry.NewVectorN(
				region.GetMax().Dimensions[0],
				center.Dimensions[1],
				region.GetMax().Dimensions[2])))
	case BRB:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				center.Dimensions[0]+1,
				center.Dimensions[1]+1,
				center.Dimensions[2]+1),
			*protometry.NewVectorN(
				region.GetMax().Dimensions[0],
				region.GetMax().Dimensions[1],
				region.GetMax().Dimensions[2])))
	case BLB:
		newNode = NewRegionOctreeNode(protometry.NewBox(
			*protometry.NewVectorN(
				region.GetMin().Dimensions[0],
				center.Dimensions[1]+1,
				center.Dimensions[2]+1),
			*protometry.NewVectorN(
				center.Dimensions[0],
				region.GetMax().Dimensions[1],
				region.GetMax().Dimensions[2])))
	}
	return newNode
}

// ToString Get a human readable representation of the state of
// this node and its contents.
func (o *OctreeNode) ToString() string {
	return o.recursiveToString("", "  ")
}

func (o *OctreeNode) recursiveToString(curIndent, stepIndent string) string {
	singleIndent := curIndent + stepIndent

	// default values
	childStr := "nil"
	positionStr := "nil"
	dataStr := "nil"

	if o.children[0] != nil {
		doubleIndent := singleIndent + stepIndent

		// accumulate child strings
		childStr = ""
		for i, child := range o.children {
			childStr = childStr + fmt.Sprintf("%v%d: %v,\n", doubleIndent, i, child.recursiveToString(doubleIndent, stepIndent))
		}

		childStr = fmt.Sprintf("[\n%v%v]", childStr, singleIndent)
	}

	if o.position != nil {
		positionStr = o.position.ToString()
	}

	if o.data != nil {
		// not stringifying elements since their type is unknown
		dataStr = fmt.Sprintf("[%d]", len(o.data))
	}

	return fmt.Sprintf("Node{\n%vposition: %v,\n%vdata: %v,\n%vregion: %v,\n%vchildren: %v,%v\n%v}", singleIndent, positionStr, singleIndent, dataStr, singleIndent, o.region, singleIndent, childStr, singleIndent, curIndent)
}
