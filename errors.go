package octree

import (
	"errors"
)

var (
	// ErrtreeOutsideBounds ...
	ErrtreeOutsideBounds       = errors.New("Position outside bounds")
	ErrtreeFailedToFindNode    = errors.New("Failed to find node")
	ErrtreeIncorrectDimensions = errors.New("Must be 3 dimensionals") // TODO: To improve
)
