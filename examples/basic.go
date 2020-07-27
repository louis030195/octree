package main

import (
	octree "github.com/louis030195/octree/pkg"
    "github.com/louis030195/protometry/api/volume"
    "log"
)

func main() {
	o := octree.NewOctree(volume.NewBoxMinMax(1, 1, 1, 4, 4, 4))
	myObj := octree.NewObjectCube(0, 2, 2, 3, 0.5)
	ok := o.Insert(*myObj)
	log.Printf("%v", ok) // true
	colliders := o.GetColliding(*volume.NewBoxMinMax(0, 0, 0, 0.9, 2.9, 0.9))
	log.Printf("%v", colliders)              // some objects
	log.Printf("%v", o.Move(myObj, 2, 2, 3)) // Success
	log.Printf("%v", o.Move(myObj, 2, 2, 8)) // Failure: outside tree !
}
