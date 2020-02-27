package octree

import (
	"errors"
)

var (
	// ErrtreeOutsideBounds ...
	ErrtreeOutsideBounds    = errors.New("Position outside bounds")
	ErrtreeFailedToFindNode = errors.New("Failed to find node")
)
