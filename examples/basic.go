package main

import (
    octree "github.com/louis030195/octree/pkg"
    protometry "github.com/louis030195/protometry/pkg"
    "log"
)

func main() {
    o := octree.NewOctree(protometry.NewBoxMinMax(1, 1, 1, 4, 4, 4))
    myObj := octree.NewObjectCube(0, 2, 2, 3, 0.5)
    ok := o.Insert(*myObj)
    log.Printf("%v", ok)// true
    colliders := o.GetColliding(*protometry.NewBoxMinMax(0, 0, 0, 0.9, 2.9, 0.9))
    log.Printf("%v", colliders) // some objects
    log.Printf("%v", o.Move(myObj, 0, 0, 0, 2, 2, 2)) // Using bounds
    log.Printf("%v", o.Move(myObj, 3, 3, 3)) // Using position, assume cube of side 1
}
